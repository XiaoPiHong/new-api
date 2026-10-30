package model

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createReserveTestUser(t *testing.T, quota int) User {
	t.Helper()
	user := User{
		Username:    "reserve-user-" + common.GetRandomString(6),
		Password:    "unused-password-hash",
		Role:        common.RoleCommonUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		AuthVersion: 1,
		Quota:       quota,
		AffCode:     "reserve-aff-" + common.GetRandomString(8),
	}
	require.NoError(t, DB.Create(&user).Error)
	return user
}

func createReserveTestToken(t *testing.T, remainQuota int) Token {
	t.Helper()
	token := Token{
		UserId:      1,
		Key:         "reserve-token-" + common.GetRandomString(8),
		Name:        "reserve-test",
		Status:      common.TokenStatusEnabled,
		ExpiredTime: -1,
		RemainQuota: remainQuota,
	}
	require.NoError(t, token.Insert())
	return token
}

func getUserQuotaFromDB(t *testing.T, id int) int {
	t.Helper()
	var user User
	require.NoError(t, DB.Select("quota").First(&user, id).Error)
	return user.Quota
}

func getTokenFromDB(t *testing.T, id int) Token {
	t.Helper()
	var token Token
	require.NoError(t, DB.First(&token, id).Error)
	return token
}

func resetBatchUpdateTestState(t *testing.T) {
	t.Helper()
	oldBatchEnabled := common.BatchUpdateEnabled
	common.BatchUpdateEnabled = false
	for i := range BatchUpdateTypeCount {
		batchUpdateLocks[i].Lock()
		batchUpdateStores[i] = make(map[int]int)
		batchUpdateLocks[i].Unlock()
	}
	t.Cleanup(func() {
		common.BatchUpdateEnabled = oldBatchEnabled
		for i := range BatchUpdateTypeCount {
			batchUpdateLocks[i].Lock()
			batchUpdateStores[i] = make(map[int]int)
			batchUpdateLocks[i].Unlock()
		}
	})
}

func TestTryReserveQuotaWithoutRedis(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)

	user := createReserveTestUser(t, 100)
	reserved, err := TryReserveUserQuota(user.Id, 60)
	require.NoError(t, err)
	assert.True(t, reserved)
	assert.Equal(t, 40, getUserQuotaFromDB(t, user.Id))

	reserved, err = TryReserveUserQuota(user.Id, 41)
	require.NoError(t, err)
	assert.False(t, reserved)
	assert.Equal(t, 40, getUserQuotaFromDB(t, user.Id))

	token := createReserveTestToken(t, 80)
	reserved, err = TryReserveTokenQuota(token.Id, token.Key, 25, false)
	require.NoError(t, err)
	assert.True(t, reserved)
	reloaded := getTokenFromDB(t, token.Id)
	assert.Equal(t, 55, reloaded.RemainQuota)
	assert.Equal(t, 25, reloaded.UsedQuota)

	reserved, err = TryReserveTokenQuota(token.Id, token.Key, 56, false)
	require.NoError(t, err)
	assert.False(t, reserved)
	assert.Equal(t, 55, getTokenFromDB(t, token.Id).RemainQuota)
}

func TestRedisBatchReserveNeverFallsBackToStaleDatabaseBalance(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	useUserCacheMiniRedis(t)
	common.BatchUpdateEnabled = true

	user := createReserveTestUser(t, 10)
	reserved, err := TryReserveUserQuota(user.Id, 8)
	require.NoError(t, err)
	assert.True(t, reserved)
	assert.Equal(t, 10, getUserQuotaFromDB(t, user.Id), "batch delta is not flushed yet")

	reserved, err = TryReserveUserQuota(user.Id, 3)
	require.NoError(t, err)
	assert.False(t, reserved, "stale DB balance must not authorize a second spend")
	cachedUser, err := GetUserCache(user.Id)
	require.NoError(t, err)
	assert.Equal(t, 2, cachedUser.Quota)

	token := createReserveTestToken(t, 9)
	reserved, err = TryReserveTokenQuota(token.Id, token.Key, 7, false)
	require.NoError(t, err)
	assert.True(t, reserved)
	reserved, err = TryReserveTokenQuota(token.Id, token.Key, 3, false)
	require.NoError(t, err)
	assert.False(t, reserved)
	assert.Equal(t, 9, getTokenFromDB(t, token.Id).RemainQuota)

	batchUpdate()
	assert.Equal(t, 2, getUserQuotaFromDB(t, user.Id))
	reloadedToken := getTokenFromDB(t, token.Id)
	assert.Equal(t, 2, reloadedToken.RemainQuota)
	assert.Equal(t, 7, reloadedToken.UsedQuota)
}

func TestBatchUpdateAccumulatesTwoMaximumRequestCharges(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	common.BatchUpdateEnabled = true

	user := createReserveTestUser(t, common.MaxQuota*2+100)
	require.NoError(t, DecreaseUserQuota(user.Id, common.MaxQuota, false))
	require.NoError(t, DecreaseUserQuota(user.Id, common.MaxQuota, false))

	batchUpdate()
	assert.Equal(t, 100, getUserQuotaFromDB(t, user.Id))
}

func TestBatchUpdateAccumulatorSaturatesOverflow(t *testing.T) {
	resetBatchUpdateTestState(t)

	addNewRecord(BatchUpdateTypeUserQuota, 1, math.MaxInt)
	addNewRecord(BatchUpdateTypeUserQuota, 1, 1)
	batchUpdateLocks[BatchUpdateTypeUserQuota].Lock()
	assert.Equal(t, math.MaxInt, batchUpdateStores[BatchUpdateTypeUserQuota][1])
	batchUpdateLocks[BatchUpdateTypeUserQuota].Unlock()

	batchUpdateLocks[BatchUpdateTypeUserQuota].Lock()
	batchUpdateStores[BatchUpdateTypeUserQuota] = make(map[int]int)
	batchUpdateLocks[BatchUpdateTypeUserQuota].Unlock()
	addNewRecord(BatchUpdateTypeUserQuota, 1, math.MinInt)
	addNewRecord(BatchUpdateTypeUserQuota, 1, -1)
	batchUpdateLocks[BatchUpdateTypeUserQuota].Lock()
	assert.Equal(t, math.MinInt, batchUpdateStores[BatchUpdateTypeUserQuota][1])
	batchUpdateLocks[BatchUpdateTypeUserQuota].Unlock()
}

func TestReserveFallsBackToDatabaseWhenRedisIsUnavailable(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	server := useUserCacheMiniRedis(t)

	user := createReserveTestUser(t, 20)
	require.NoError(t, populateUserCache(user))
	server.Close()

	// Redis 故障时降级为数据库条件更新：服务保持可用且不会超扣。
	reserved, err := TryReserveUserQuota(user.Id, 5)
	require.NoError(t, err)
	assert.True(t, reserved)
	assert.Equal(t, 15, getUserQuotaFromDB(t, user.Id))

	reserved, err = TryReserveUserQuota(user.Id, 16)
	require.NoError(t, err)
	assert.False(t, reserved)
	assert.Equal(t, 15, getUserQuotaFromDB(t, user.Id))
}

func TestSynchronousReserveCompensatesCacheWhenPersistenceFails(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	useUserCacheMiniRedis(t)

	user := createReserveTestUser(t, 10)
	require.NoError(t, populateUserCache(user))
	require.NoError(t, DB.Delete(&user).Error)

	reserved, err := TryReserveUserQuota(user.Id, 6)
	assert.False(t, reserved)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	cached, cacheErr := cacheGetUserBase(user.Id)
	require.NoError(t, cacheErr)
	assert.Equal(t, 10, cached.Quota)

	token := createReserveTestToken(t, 12)
	_, err = GetTokenByKey(token.Key, true)
	require.NoError(t, err)
	require.NoError(t, DB.Delete(&token).Error)
	reserved, err = TryReserveTokenQuota(token.Id, token.Key, 7, false)
	assert.False(t, reserved)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	cachedToken, cacheErr := cacheGetTokenByKey(token.Key)
	require.NoError(t, cacheErr)
	assert.Equal(t, 12, cachedToken.RemainQuota)
	assert.Zero(t, cachedToken.UsedQuota)
}

func TestTokenCacheInitPreservesLiveQuotaAndFenceBlocksStaleSnapshot(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	server := useUserCacheMiniRedis(t)

	token := createReserveTestToken(t, 100)
	loaded, err := GetTokenByKey(token.Key, true)
	require.NoError(t, err)
	stale := *loaded

	result, err := cacheApplyTokenQuotaDelta(token.Id, token.Key, -70)
	require.NoError(t, err)
	require.Equal(t, cacheQuotaOK, result)

	// 已存在的哈希只刷新 TTL：数据库快照不得覆盖已被原子预扣的余额。
	code, err := cacheInitToken(stale)
	require.NoError(t, err)
	assert.Equal(t, 2, code)
	cached, err := cacheGetTokenByKey(token.Key)
	require.NoError(t, err)
	assert.Equal(t, 30, cached.RemainQuota)

	// 变更期间：fence 删除缓存并拦截并发读者手中的过期快照。
	require.NoError(t, invalidateTokenCacheForMutation(token.Key))
	code, err = cacheInitToken(stale)
	require.NoError(t, err)
	assert.Zero(t, code, "the pre-mutation snapshot must not be published while fenced")
	_, err = cacheGetTokenByKey(token.Key)
	assert.Error(t, err)

	// fence 过期后可重新从数据库水合。
	server.FastForward(time.Duration(tokenCacheFenceSeconds+1) * time.Second)
	fresh, err := GetTokenByKey(token.Key, false)
	require.NoError(t, err)
	assert.Equal(t, 100, fresh.RemainQuota)
	cached, err = cacheGetTokenByKey(token.Key)
	require.NoError(t, err)
	assert.Equal(t, 100, cached.RemainQuota)
}

func TestBillingRefundCachePreservesConcurrentReservations(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	useUserCacheMiniRedis(t)
	user := createReserveTestUser(t, 920)
	token := createReserveTestToken(t, 920)
	token.UsedQuota = 80
	require.NoError(t, populateUserCache(user))
	_, err := cacheInitToken(token)
	require.NoError(t, err)
	changes := []billingQuotaCacheChange{
		{key: getUserCacheKey(user.Id), id: user.Id, credit: 80, operation: "refund-user"},
		{key: getTokenCacheKey(token.Key), id: token.Id, token: true, credit: 80, operation: "refund-token"},
	}
	for _, change := range changes {
		require.NoError(t, change.begin())
	}
	// Another request has reserved Redis quota but has not yet persisted it.
	result, err := cacheTryReserveUserQuota(user.Id, 100)
	require.NoError(t, err)
	require.Equal(t, cacheQuotaOK, result)
	result, err = cacheTryReserveTokenQuota(token.Id, token.Key, 100)
	require.NoError(t, err)
	require.Equal(t, cacheQuotaOK, result)
	for _, change := range changes {
		require.NoError(t, change.finish(true))
		require.NoError(t, change.finish(true), "lost responses must be safe to retry")
	}
	cachedUser, err := cacheGetUserBase(user.Id)
	require.NoError(t, err)
	assert.Equal(t, 900, cachedUser.Quota, "refund must preserve the concurrent 100 quota reservation")
	cachedToken, err := cacheGetTokenByKey(token.Key)
	require.NoError(t, err)
	assert.Equal(t, 900, cachedToken.RemainQuota)
	assert.Equal(t, 100, cachedToken.UsedQuota)
}

func TestBillingRefundCacheFencesColdAndDelayedSnapshots(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	useUserCacheMiniRedis(t)
	user := createReserveTestUser(t, 920)
	token := createReserveTestToken(t, 920)
	token.UsedQuota = 80
	userKey, tokenKey := getUserCacheKey(user.Id), getTokenCacheKey(token.Key)
	beforeUser, err := readBillingQuotaEpoch(userKey)
	require.NoError(t, err)
	beforeToken, err := readBillingQuotaEpoch(tokenKey)
	require.NoError(t, err)
	changes := []billingQuotaCacheChange{
		{key: userKey, id: user.Id, credit: 80, operation: "cold-user"},
		{key: tokenKey, id: token.Id, token: true, credit: 80, operation: "cold-token"},
	}
	for _, change := range changes {
		require.NoError(t, change.begin())
	}
	duringUser, err := readBillingQuotaEpoch(userKey)
	require.NoError(t, err)
	duringToken, err := readBillingQuotaEpoch(tokenKey)
	require.NoError(t, err)
	// Even a fresh post-commit DB snapshot cannot be cached before finish.
	user.Quota = 1000
	token.RemainQuota, token.UsedQuota = 1000, 0
	require.NoError(t, populateUserCache(user, duringUser))
	code, err := cacheInitToken(token, duringToken)
	require.NoError(t, err)
	assert.Zero(t, code)
	assert.EqualValues(t, 0, common.RDB.Exists(context.Background(), userKey, tokenKey).Val())
	for _, change := range changes {
		require.NoError(t, change.finish(true))
	}
	// A reader delayed from either before or during the transaction is rejected.
	for _, epoch := range []string{beforeUser, duringUser} {
		require.NoError(t, populateUserCache(user, epoch))
	}
	for _, epoch := range []string{beforeToken, duringToken} {
		code, err := cacheInitToken(token, epoch)
		require.NoError(t, err)
		assert.Zero(t, code)
	}
	assert.EqualValues(t, 0, common.RDB.Exists(context.Background(), userKey, tokenKey).Val())
	afterUser, err := readBillingQuotaEpoch(userKey)
	require.NoError(t, err)
	afterToken, err := readBillingQuotaEpoch(tokenKey)
	require.NoError(t, err)
	require.NoError(t, populateUserCache(user, afterUser))
	code, err = cacheInitToken(token, afterToken)
	require.NoError(t, err)
	assert.Equal(t, 1, code)
	for _, change := range changes {
		require.NoError(t, change.finish(true))
	}
	cachedUser, err := cacheGetUserBase(user.Id)
	require.NoError(t, err)
	assert.Equal(t, 1000, cachedUser.Quota, "new DB snapshot must not get another 80 credit")
	cachedToken, err := cacheGetTokenByKey(token.Key)
	require.NoError(t, err)
	assert.Equal(t, 1000, cachedToken.RemainQuota)
	assert.Zero(t, cachedToken.UsedQuota)
	// The new quota guards cannot override either existing auth revocation fence.
	require.NoError(t, SetUserAuthVersionFence(user.Id, user.AuthVersion+1))
	assert.ErrorIs(t, populateUserCache(user, afterUser), ErrUserAuthCachePending)
	require.NoError(t, invalidateTokenCacheForMutation(token.Key))
	code, err = cacheInitToken(token, afterToken)
	require.NoError(t, err)
	assert.Zero(t, code)
}

func TestBillingSettlementCacheRollbackAndPendingDatabaseFallback(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	server := useUserCacheMiniRedis(t)
	user := createReserveTestUser(t, 100)
	token := createReserveTestToken(t, 100)
	require.NoError(t, populateUserCache(user))
	_, err := cacheInitToken(token)
	require.NoError(t, err)
	changes := []billingQuotaCacheChange{
		{key: getUserCacheKey(user.Id), id: user.Id, credit: -80, operation: "charge-user"},
		{key: getTokenCacheKey(token.Key), id: token.Id, token: true, credit: -80, operation: "charge-token"},
	}
	for _, change := range changes {
		require.NoError(t, change.begin())
	}
	result, err := cacheTryReserveUserQuota(user.Id, 21)
	require.NoError(t, err)
	assert.Equal(t, cacheQuotaInsufficient, result)
	result, err = cacheTryReserveTokenQuota(token.Id, token.Key, 21)
	require.NoError(t, err)
	assert.Equal(t, cacheQuotaInsufficient, result)
	result, err = cacheTryReserveUserQuota(user.Id, 10)
	require.NoError(t, err)
	require.Equal(t, cacheQuotaOK, result)
	result, err = cacheTryReserveTokenQuota(token.Id, token.Key, 10)
	require.NoError(t, err)
	require.Equal(t, cacheQuotaOK, result)
	for _, change := range changes {
		require.NoError(t, change.finish(false))
		require.NoError(t, change.finish(false))
	}
	cachedUser, err := cacheGetUserBase(user.Id)
	require.NoError(t, err)
	assert.Equal(t, 90, cachedUser.Quota)
	cachedToken, err := cacheGetTokenByKey(token.Key)
	require.NoError(t, err)
	assert.Equal(t, 90, cachedToken.RemainQuota)
	assert.Equal(t, 10, cachedToken.UsedQuota)

	// Simulate an interrupted process after begin: once hashes expire, pending
	// fences prevent rehydration but DB conditional reservations remain bounded.
	for i := range changes {
		changes[i].operation += "-interrupted"
		require.NoError(t, changes[i].begin())
	}
	server.FastForward(time.Duration(max(userCacheTTLSeconds(), tokenCacheTTLSeconds())+1) * time.Second)
	reserved, err := TryReserveUserQuota(user.Id, 70)
	require.NoError(t, err)
	assert.True(t, reserved)
	reserved, err = TryReserveUserQuota(user.Id, 31)
	require.NoError(t, err)
	assert.False(t, reserved)
	reserved, err = TryReserveTokenQuota(token.Id, token.Key, 70, false)
	require.NoError(t, err)
	assert.True(t, reserved)
	reserved, err = TryReserveTokenQuota(token.Id, token.Key, 31, false)
	require.NoError(t, err)
	assert.False(t, reserved)
	assert.EqualValues(t, 0, common.RDB.Exists(context.Background(), changes[0].key, changes[1].key).Val())
	assert.EqualValues(t, 2, common.RDB.Exists(context.Background(), billingQuotaPendingKey(changes[0].key), billingQuotaPendingKey(changes[1].key)).Val())
}

func TestBillingQuotaCacheDoesNotCompensateUnknownCommit(t *testing.T) {
	truncateTables(t)
	resetBatchUpdateTestState(t)
	useUserCacheMiniRedis(t)
	user := createReserveTestUser(t, 100)
	require.NoError(t, populateUserCache(user))
	change := billingQuotaCacheChange{key: getUserCacheKey(user.Id), id: user.Id, credit: -80, operation: "ambiguous-commit"}
	require.NoError(t, change.begin())
	resolveBillingQuotaCaches([]billingQuotaCacheChange{change}, true, errors.New("commit response lost"))
	cached, err := cacheGetUserBase(user.Id)
	require.NoError(t, err)
	assert.Equal(t, 20, cached.Quota, "an ambiguously committed charge must not become spendable again")
	assert.EqualValues(t, 1, common.RDB.Exists(context.Background(), billingQuotaPendingKey(change.key)).Val())
	assert.True(t, common.RDB.HExists(context.Background(), change.key, "BillingQuotaOp:"+change.operation).Val())
}
