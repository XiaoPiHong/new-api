package model

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// These tables live only in the explicitly configured TEST_* databases; the
// normal application DSN is never used. A prefix also separates this fixture
// from other model tests sharing the same disposable database.
func workflowQuotaTestDB(t *testing.T, dialect string) *gorm.DB {
	t.Helper()
	var driver gorm.Dialector
	switch dialect {
	case "sqlite":
		driver = sqlite.Open(":memory:")
	case "mysql":
		dsn := os.Getenv("TEST_MYSQL_DSN")
		if dsn == "" {
			t.Skip("TEST_MYSQL_DSN is not configured")
		}
		driver = mysqlMigrationDialector{mysql.Dialector{Config: &mysql.Config{DSN: dsn}}}
	case "postgres":
		dsn := os.Getenv("TEST_POSTGRES_DSN")
		if dsn == "" {
			t.Skip("TEST_POSTGRES_DSN is not configured")
		}
		driver = postgresMigrationDialector{postgres.Dialector{Config: &postgres.Config{DSN: dsn, PreferSimpleProtocol: true}}}
	}
	db, err := gorm.Open(driver, &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "workflow_test_"}})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	if dialect == "sqlite" {
		// SQLite has no FOR UPDATE; one connection gives its transaction model
		// deterministic serialization without pretending to test row locks.
		sqlDB.SetMaxOpenConns(1)
	} else {
		sqlDB.SetMaxOpenConns(8)
	}
	previousDB, previousLogDB := DB, LOG_DB
	previousType, previousLogType := common.MainDatabaseType(), common.LogDatabaseType()
	previousRedis, previousBatch := common.RedisEnabled, common.BatchUpdateEnabled
	DB, LOG_DB = db, db
	common.SetDatabaseTypes(common.DatabaseType(dialect), common.DatabaseType(dialect))
	common.RedisEnabled, common.BatchUpdateEnabled = false, false
	initCol()
	t.Cleanup(func() {
		assert.NoError(t, db.Migrator().DropTable(&Log{}, &TopUp{}, &Task{}, &WorkflowQuotaReservationItem{}, &WorkflowQuotaReservation{}, &Token{}, &User{}))
		DB, LOG_DB = previousDB, previousLogDB
		common.SetDatabaseTypes(previousType, previousLogType)
		common.RedisEnabled, common.BatchUpdateEnabled = previousRedis, previousBatch
		initCol()
		assert.NoError(t, sqlDB.Close())
	})
	require.NoError(t, db.AutoMigrate(&User{}, &Token{}, &WorkflowQuotaReservation{}, &WorkflowQuotaReservationItem{}, &Task{}, &TopUp{}, &Log{}))
	return db
}

func workflowQuotaTestOwner(t *testing.T, db *gorm.DB, id int, wallet, tokenQuota int, unlimited bool) {
	t.Helper()
	require.NoError(t, db.Create(&User{
		Id: id, Username: fmt.Sprintf("workflow_%d", id), AffCode: fmt.Sprintf("workflow_%d", id),
		Status: common.UserStatusEnabled, Quota: wallet,
	}).Error)
	require.NoError(t, db.Create(&Token{
		Id: id, UserId: id, Key: fmt.Sprintf("workflow-test-token-%d", id),
		RemainQuota: tokenQuota, UnlimitedQuota: unlimited,
	}).Error)
}

func workflowQuotaTestBalance(t *testing.T, db *gorm.DB, id, wallet, remain, used int) {
	t.Helper()
	var user User
	var token Token
	require.NoError(t, db.First(&user, id).Error)
	require.NoError(t, db.First(&token, id).Error)
	assert.Equal(t, wallet, user.Quota, "wallet")
	assert.Equal(t, remain, token.RemainQuota, "token remaining quota")
	assert.Equal(t, used, token.UsedQuota, "token used quota")
}

func workflowQuotaTestReserve(t *testing.T, id int, key string, expiresAt int64) *WorkflowQuotaReservation {
	t.Helper()
	reservation, created, err := ReserveWorkflowQuota(ReserveWorkflowQuotaRequest{
		ReservationKey: key, UserId: id, TokenId: id, ExpiresAt: expiresAt,
		Items: []WorkflowQuotaReservationItemInput{
			{ItemKey: "chat", Adapter: "newapi.chat", ModelName: "test-chat", QuotedQuota: 100},
			{ItemKey: "video", Adapter: "newapi.video", ModelName: "test-video", QuotedQuota: 200},
		},
	})
	require.NoError(t, err)
	require.True(t, created)
	return reservation
}

func TestWorkflowQuotaDatabaseMatrix(t *testing.T) {
	for _, dialect := range []string{"sqlite", "mysql", "postgres"} {
		t.Run(dialect, func(t *testing.T) {
			db := workflowQuotaTestDB(t, dialect)
			t.Run("settlement_release_and_final_refund", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 1, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 1, "settle", common.GetTimestamp()+3600)
				retried, created, err := ReserveWorkflowQuota(ReserveWorkflowQuotaRequest{
					ReservationKey: "settle", UserId: 1, TokenId: 1,
					Items: []WorkflowQuotaReservationItemInput{
						{ItemKey: "chat", Adapter: "newapi.chat", ModelName: "test-chat", QuotedQuota: 100},
						{ItemKey: "video", Adapter: "newapi.video", ModelName: "test-video", QuotedQuota: 200},
					},
				})
				require.NoError(t, err)
				assert.False(t, created)
				assert.Equal(t, reservation.PublicId, retried.PublicId)
				workflowQuotaTestBalance(t, db, 1, 700, 700, 300)
				require.NoError(t, ClaimWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 80, "request-1"))
				require.NoError(t, ClaimWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 80, "request-1"))
				require.ErrorIs(t, ClaimWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 80, "another-request"), ErrWorkflowQuotaItemConflict)
				require.ErrorIs(t, ClaimWorkflowQuotaItem(reservation.PublicId, "video", 2, 1, 10, "wrong-owner"), ErrWorkflowQuotaReservationNotFound)
				require.NoError(t, SettleWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 80))
				require.NoError(t, SettleWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 80))
				require.ErrorIs(t, SettleWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 81), ErrWorkflowQuotaItemConflict)
				require.NoError(t, AdjustSettledWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 60))
				require.NoError(t, AdjustSettledWorkflowQuotaItem(reservation.PublicId, "chat", 1, 1, 60))
				workflowQuotaTestBalance(t, db, 1, 700, 700, 300)
				released, changed, err := ReleaseWorkflowQuotaReservation(reservation.PublicId, 1, 1, false)
				require.NoError(t, err)
				require.True(t, changed)
				assert.Equal(t, WorkflowQuotaReservationStatusSettled, released.Status)
				assert.Equal(t, 240, released.ReleasedQuota)
				workflowQuotaTestBalance(t, db, 1, 940, 940, 60)
				_, changed, err = ReleaseWorkflowQuotaReservation(reservation.PublicId, 1, 1, false)
				require.NoError(t, err)
				assert.False(t, changed)
				_, changed, err = ReleaseWorkflowQuotaReservation(reservation.PublicId, 1, 1, true)
				require.NoError(t, err)
				assert.True(t, changed)
				_, changed, err = ReleaseWorkflowQuotaReservation(reservation.PublicId, 1, 1, true)
				require.NoError(t, err)
				assert.False(t, changed)
				workflowQuotaTestBalance(t, db, 1, 1000, 1000, 0)
			})
			t.Run("concurrent_release", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 2, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 2, "parallel-release", common.GetTimestamp()+3600)
				var wg sync.WaitGroup
				var changedCount atomic.Int32
				errs := make(chan error, 6)
				for i := 0; i < 6; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						_, changed, err := ReleaseWorkflowQuotaReservation(reservation.PublicId, 2, 2, true)
						if changed {
							changedCount.Add(1)
						}
						errs <- err
					}()
				}
				wg.Wait()
				close(errs)
				for err := range errs {
					require.NoError(t, err)
				}
				assert.EqualValues(t, 1, changedCount.Load())
				workflowQuotaTestBalance(t, db, 2, 1000, 1000, 0)
			})
			t.Run("expiry_and_idempotent_retry", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 3, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 3, "expired", common.GetTimestamp()-1)
				require.ErrorIs(t, ClaimWorkflowQuotaItem(reservation.PublicId, "chat", 3, 3, 10, "late"), ErrWorkflowQuotaReservationExpired)
				released, err := ReleaseExpiredWorkflowQuotaReservations(100)
				require.NoError(t, err)
				assert.Equal(t, 1, released)
				released, err = ReleaseExpiredWorkflowQuotaReservations(100)
				require.NoError(t, err)
				assert.Zero(t, released)
				workflowQuotaTestBalance(t, db, 3, 1000, 1000, 0)
			})
			t.Run("insufficient_token_rolls_back_wallet", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 4, 1000, 299, false)
				_, _, err := ReserveWorkflowQuota(ReserveWorkflowQuotaRequest{
					ReservationKey: "insufficient", UserId: 4, TokenId: 4,
					Items: []WorkflowQuotaReservationItemInput{{ItemKey: "one", QuotedQuota: 300}},
				})
				require.ErrorIs(t, err, ErrWorkflowQuotaTokenInsufficient)
				workflowQuotaTestBalance(t, db, 4, 1000, 299, 0)
				var count int64
				require.NoError(t, db.Model(&WorkflowQuotaReservation{}).Where("reservation_key = ?", "insufficient").Count(&count).Error)
				assert.Zero(t, count)
			})
			t.Run("unlimited_token_and_migration_preservation", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 5, 1000, 7, true)
				reservation := workflowQuotaTestReserve(t, 5, "unlimited", common.GetTimestamp()+3600)
				workflowQuotaTestBalance(t, db, 5, 700, 7, 0)
				// The two workflow schemas are unchanged from the custom build.
				// Re-run migration twice with representative persisted lifecycle data.
				for i := 0; i < 2; i++ {
					require.NoError(t, db.AutoMigrate(&WorkflowQuotaReservation{}, &WorkflowQuotaReservationItem{}))
				}
				var stored WorkflowQuotaReservation
				require.NoError(t, db.Preload("Items").First(&stored, reservation.Id).Error)
				assert.Equal(t, reservation.PublicId, stored.PublicId)
				assert.Equal(t, reservation.ReservedQuota, stored.ReservedQuota)
				require.Len(t, stored.Items, 2)
				_, _, err := ReleaseWorkflowQuotaReservation(reservation.PublicId, 5, 5, true)
				require.NoError(t, err)
				workflowQuotaTestBalance(t, db, 5, 1000, 7, 0)
			})
			t.Run("zero_settlement_closes_the_item", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 6, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 6, "free-result", common.GetTimestamp()+3600)
				require.NoError(t, ClaimWorkflowQuotaItem(reservation.PublicId, "chat", 6, 6, 80, "zero"))
				require.NoError(t, SettleWorkflowQuotaItem(reservation.PublicId, "chat", 6, 6, 0))
				require.NoError(t, SettleWorkflowQuotaItem(reservation.PublicId, "chat", 6, 6, 0))
				require.ErrorIs(t, ClaimWorkflowQuotaItem(reservation.PublicId, "chat", 6, 6, 1, "reuse"), ErrWorkflowQuotaItemConflict)
				_, _, err := ReleaseWorkflowQuotaReservation(reservation.PublicId, 6, 6, false)
				require.NoError(t, err)
				workflowQuotaTestBalance(t, db, 6, 1000, 1000, 0)
			})
			t.Run("deleted_token_does_not_block_wallet_refund", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 7, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 7, "deleted-token", common.GetTimestamp()+3600)
				require.NoError(t, db.Delete(&Token{}, 7).Error)
				_, changed, err := ReleaseWorkflowQuotaReservationInternal(reservation.PublicId, true)
				require.NoError(t, err)
				assert.True(t, changed)
				var user User
				var token Token
				require.NoError(t, db.First(&user, 7).Error)
				require.NoError(t, db.Unscoped().First(&token, 7).Error)
				assert.Equal(t, 1000, user.Quota)
				assert.Equal(t, 1000, token.RemainQuota)
				assert.Zero(t, token.UsedQuota)
				assert.True(t, token.DeletedAt.Valid, "refund must not re-enable the deleted token")
			})
			t.Run("missing_token_refunds_wallet_without_recreation", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 8, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 8, "missing-token", common.GetTimestamp()+3600)
				require.NoError(t, db.Unscoped().Delete(&Token{}, 8).Error)
				_, _, err := ReleaseWorkflowQuotaReservationInternal(reservation.PublicId, true)
				require.NoError(t, err)
				var user User
				var count int64
				require.NoError(t, db.First(&user, 8).Error)
				require.NoError(t, db.Unscoped().Model(&Token{}).Where("id = ?", 8).Count(&count).Error)
				assert.Equal(t, 1000, user.Quota)
				assert.Zero(t, count)
			})
			t.Run("reused_token_id_never_credits_another_owner", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 9, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 9, "foreign-token", common.GetTimestamp()+3600)
				require.NoError(t, db.Unscoped().Delete(&Token{}, 9).Error)
				require.NoError(t, db.Create(&Token{Id: 9, UserId: 999, Key: "foreign-owner-token", RemainQuota: 77}).Error)
				_, _, err := ReleaseWorkflowQuotaReservationInternal(reservation.PublicId, true)
				require.NoError(t, err)
				var user User
				var token Token
				require.NoError(t, db.First(&user, 9).Error)
				require.NoError(t, db.First(&token, 9).Error)
				assert.Equal(t, 1000, user.Quota)
				assert.Equal(t, 999, token.UserId)
				assert.Equal(t, 77, token.RemainQuota)
				assert.Zero(t, token.UsedQuota)
			})
			t.Run("expiry_can_release_a_soft_deleted_token", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 10, 1000, 1000, false)
				workflowQuotaTestReserve(t, 10, "expired-deleted-token", common.GetTimestamp()-1)
				require.NoError(t, db.Delete(&Token{}, 10).Error)
				count, err := ReleaseExpiredWorkflowQuotaReservations(100)
				require.NoError(t, err)
				assert.Equal(t, 1, count)
				var user User
				var token Token
				require.NoError(t, db.First(&user, 10).Error)
				require.NoError(t, db.Unscoped().First(&token, 10).Error)
				assert.Equal(t, 1000, user.Quota)
				assert.Equal(t, 1000, token.RemainQuota)
				assert.Zero(t, token.UsedQuota)
				assert.True(t, token.DeletedAt.Valid)
			})
			t.Run("legacy_task_trace_and_workflow_json_roundtrip", func(t *testing.T) {
				for i, private := range []TaskPrivateData{
					{BillingTraceId: "trace-only"},
					{
						BillingTraceId: "workflow-trace", BillingSource: "workflow_reservation",
						WorkflowQuotaReservationId: "wqr_fixture", WorkflowQuotaItemKey: "video", TokenId: 6,
						UpstreamTaskID: "vendor-task", ResultURL: "https://example.invalid/result.mp4", NodeName: "fixture-node",
						Execution: &TaskExecutionSnapshot{RequestID: "request-fixture", RequestPath: "/v1/videos"},
					},
				} {
					task := Task{
						TaskID: fmt.Sprintf("legacy-workflow-task-%d", i), Platform: constant.TaskPlatform("58"),
						UserId: 6, Quota: 10, BillingTraceId: private.BillingTraceId, PrivateData: private,
					}
					require.NoError(t, db.Create(&task).Error)
					stored, exists, err := GetByBillingTraceId(private.BillingTraceId)
					require.NoError(t, err)
					require.True(t, exists)
					assert.Equal(t, task.Platform, stored.Platform)
					assert.Equal(t, private, stored.PrivateData)
				}
			})
			t.Run("concurrent_payment_callback_keeps_workflow_wallet_balanced", func(t *testing.T) {
				workflowQuotaTestOwner(t, db, 11, 1000, 1000, false)
				reservation := workflowQuotaTestReserve(t, 11, "payment-workflow", common.GetTimestamp()+3600)
				require.NoError(t, db.Create(&TopUp{
					UserId: 11, Amount: 2, Money: 9.99, TradeNo: "workflow-payment-order",
					PaymentMethod: "alipay", PaymentProvider: PaymentProviderEpay, Status: common.TopUpStatusPending,
				}).Error)
				var wg sync.WaitGroup
				var creditedCount atomic.Int32
				errs := make(chan error, 4)
				for i := 0; i < 4; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						credited, err := CompleteEpayTopUp("workflow-payment-order", "alipay", "9.99", "127.0.0.1")
						if credited {
							creditedCount.Add(1)
						}
						errs <- err
					}()
				}
				wg.Wait()
				close(errs)
				for err := range errs {
					require.NoError(t, err)
				}
				assert.EqualValues(t, 1, creditedCount.Load())
				credit, err := common.WalletQuotaFromDecimalStrict(decimal.NewFromInt(2).Mul(decimal.NewFromFloat(common.QuotaPerUnit)))
				require.NoError(t, err)
				workflowQuotaTestBalance(t, db, 11, 700+credit, 700, 300)
				_, _, err = ReleaseWorkflowQuotaReservationInternal(reservation.PublicId, true)
				require.NoError(t, err)
				workflowQuotaTestBalance(t, db, 11, 1000+credit, 1000, 0)
			})
			t.Run("redis_quota_matches_committed_reservation", func(t *testing.T) {
				useUserCacheMiniRedis(t)
				workflowQuotaTestOwner(t, db, 12, 1000, 1000, false)
				var user User
				var token Token
				require.NoError(t, db.First(&user, 12).Error)
				require.NoError(t, db.First(&token, 12).Error)
				require.NoError(t, populateUserCache(user))
				require.NoError(t, cacheSetTokenForTest(token))
				reservation := workflowQuotaTestReserve(t, 12, "redis-reservation", common.GetTimestamp()+3600)
				cachedUser, err := cacheGetUserBase(12)
				require.NoError(t, err)
				cachedToken, err := cacheGetTokenByKey(token.Key)
				require.NoError(t, err)
				assert.Equal(t, 700, cachedUser.Quota)
				assert.Equal(t, 700, cachedToken.RemainQuota)
				assert.Equal(t, 300, cachedToken.UsedQuota)
				_, _, err = ReleaseWorkflowQuotaReservationInternal(reservation.PublicId, true)
				require.NoError(t, err)
				cachedUser, err = cacheGetUserBase(12)
				require.NoError(t, err)
				cachedToken, err = cacheGetTokenByKey(token.Key)
				require.NoError(t, err)
				assert.Equal(t, 1000, cachedUser.Quota)
				assert.Equal(t, 1000, cachedToken.RemainQuota)
				assert.Zero(t, cachedToken.UsedQuota)
				workflowQuotaTestBalance(t, db, 12, 1000, 1000, 0)
			})
		})
	}
}
