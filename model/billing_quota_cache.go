package model

import (
	"context"
	"errors"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// These keys are independent of authentication versions and revocation fences.
// The epoch is read BEFORE the database snapshot and checked atomically when a
// cold cache is filled. Pending operations have no expiry: an interrupted funds
// transaction must not silently allow an old snapshot to become spendable.
func billingQuotaEpochKey(cacheKey string) string {
	return "billing:quota:epoch:" + cacheKey
}

func billingQuotaPendingKey(cacheKey string) string {
	return "billing:quota:pending:" + cacheKey
}

func readBillingQuotaEpoch(cacheKey string) (string, error) {
	if !common.RedisEnabled {
		return "0", nil
	}
	epoch, err := common.RDB.Get(context.Background(), billingQuotaEpochKey(cacheKey)).Result()
	if errors.Is(err, redis.Nil) {
		return "0", nil
	}
	return epoch, err
}

type billingQuotaCacheChange struct {
	key       string
	id        int
	token     bool
	unlimited bool
	credit    int
	operation string
}

// A marker belongs to the existing hash, never to a newly hydrated snapshot.
// Refunds are added only after commit. Charges reserve the old cache before the
// DB write so a concurrent cache pre-consume cannot spend the same funds.
const beginBillingQuotaCacheScript = `
if redis.call('HEXISTS', KEYS[3], ARGV[1]) == 1 then return 1 end
local matches = redis.call('HGET', KEYS[1], 'Id') == ARGV[2]
local credit = tonumber(ARGV[3])
local balanceField = 'Quota'
if ARGV[4] == '1' then
  balanceField = 'RemainQuota'
  matches = matches and redis.call('HEXISTS', KEYS[1], 'UsedQuota') == 1
else
  matches = matches and redis.call('HGET', KEYS[1], 'CacheSchema') == ARGV[5]
end
local balance = tonumber(redis.call('HGET', KEYS[1], balanceField))
matches = matches and balance ~= nil
if matches and credit < 0 then
  if balance + credit < -tonumber(ARGV[6]) then return 0 end
  if ARGV[7] ~= '1' and balance + credit < 0 then return 0 end
  if ARGV[4] == '1' then
    local used = tonumber(redis.call('HGET', KEYS[1], 'UsedQuota'))
    if not used or used < 0 or used - credit > tonumber(ARGV[6]) then return 0 end
  end
end
redis.call('INCR', KEYS[2])
redis.call('HSET', KEYS[3], ARGV[1], ARGV[3])
if matches then
  redis.call('HSET', KEYS[1], 'BillingQuotaOp:' .. ARGV[1], ARGV[3])
  if credit < 0 then
    redis.call('HINCRBY', KEYS[1], balanceField, ARGV[3])
    if ARGV[4] == '1' then redis.call('HINCRBY', KEYS[1], 'UsedQuota', -credit) end
  end
end
return 1`

const finishBillingQuotaCacheScript = `
local creditText = redis.call('HGET', KEYS[3], ARGV[1])
if not creditText then return 1 end
local marker = 'BillingQuotaOp:' .. ARGV[1]
if redis.call('HGET', KEYS[1], marker) == creditText then
  if redis.call('HGET', KEYS[1], 'Id') ~= ARGV[2] then return 0 end
  local credit = tonumber(creditText)
  local adjustment = 0
  if ARGV[3] == '1' and credit > 0 then adjustment = credit end
  if ARGV[3] == '0' and credit < 0 then adjustment = -credit end
  local field = 'Quota'
  if ARGV[4] == '1' then field = 'RemainQuota' end
  if adjustment ~= 0 then
    local balance = tonumber(redis.call('HGET', KEYS[1], field))
    if not balance or balance + adjustment > tonumber(ARGV[5]) then return 0 end
    if ARGV[4] == '1' then
      local used = tonumber(redis.call('HGET', KEYS[1], 'UsedQuota'))
      if not used or used < adjustment then return 0 end
      redis.call('HINCRBY', KEYS[1], 'UsedQuota', -adjustment)
    end
    redis.call('HINCRBY', KEYS[1], field, adjustment)
  end
  redis.call('HDEL', KEYS[1], marker)
end
redis.call('HDEL', KEYS[3], ARGV[1])
redis.call('INCR', KEYS[2])
return 1`

func (change billingQuotaCacheChange) keys() []string {
	return []string{change.key, billingQuotaEpochKey(change.key), billingQuotaPendingKey(change.key)}
}

func (change billingQuotaCacheChange) begin() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := common.RDB.Eval(ctx, beginBillingQuotaCacheScript, change.keys(),
		change.operation, change.id, change.credit, boolToInt(change.token),
		userCacheSchemaVersion, common.MaxWalletQuota, boolToInt(change.unlimited)).Int()
	if err != nil {
		return err
	}
	if result != 1 {
		return ErrBillingRefundConflict
	}
	return nil
}

func (change billingQuotaCacheChange) finish(committed bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := common.RDB.Eval(ctx, finishBillingQuotaCacheScript, change.keys(),
		change.operation, change.id, boolToInt(committed), boolToInt(change.token), common.MaxWalletQuota).Int()
	if err != nil {
		return err
	}
	if result != 1 {
		return ErrBillingRefundConflict
	}
	return nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// Start inside the DB transaction, before any funds write. Even when Redis
// reports an error, return attempted changes so the caller can compensate a
// command whose response was lost. No transaction may commit after a begin error.
func beginBillingQuotaCaches(tx *gorm.DB, userID, tokenID int, source string, credit int) ([]billingQuotaCacheChange, error) {
	if !common.RedisEnabled || credit == 0 || source == "workflow_reservation" {
		return nil, nil
	}
	operation, err := common.GenerateRandomCharsKey(24)
	if err != nil {
		return nil, err
	}
	changes := make([]billingQuotaCacheChange, 0, 2)
	if source == "wallet" || source == "" {
		changes = append(changes, billingQuotaCacheChange{key: getUserCacheKey(userID), id: userID, credit: credit, operation: operation})
	}
	if tokenID > 0 {
		var token Token
		err := tx.Unscoped().Select("id", "key", "user_id", "unlimited_quota").First(&token, tokenID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if err == nil {
			if token.UserId != userID {
				return nil, ErrBillingRefundConflict
			}
			changes = append(changes, billingQuotaCacheChange{key: getTokenCacheKey(token.Key), id: tokenID,
				token: true, unlimited: token.UnlimitedQuota, credit: credit, operation: operation})
		}
	}
	for i, change := range changes {
		if err := change.begin(); err != nil {
			return changes[:i+1], err
		}
	}
	return changes, nil
}

func finishBillingQuotaCaches(changes []billingQuotaCacheChange, committed bool) {
	for _, change := range changes {
		if err := change.finish(committed); err != nil {
			// A retry is idempotent even when the first response was lost. Keep
			// the pending fence on failure; never undo committed database funds.
			if retryErr := change.finish(committed); retryErr != nil {
				common.SysError(common.LogText("billing quota cache operation %s remains pending: %s", change.operation, retryErr.Error()))
			}
		}
	}
}

// A successful transaction callback followed by a Commit error is ambiguous:
// the database may have committed before the response was lost. Compensating a
// cached charge in that case would create spendable funds. Keep its marker and
// fence until the database outcome is reconciled; cold reads fall back to DB.
func resolveBillingQuotaCaches(changes []billingQuotaCacheChange, callbackCompleted bool, transactionErr error) {
	if transactionErr == nil || !callbackCompleted {
		finishBillingQuotaCaches(changes, transactionErr == nil)
		return
	}
	for _, change := range changes {
		common.SysError(common.LogText("billing quota cache operation %s has an unknown database commit outcome", change.operation))
	}
}
