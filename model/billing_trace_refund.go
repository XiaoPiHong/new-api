package model

import (
	"context"
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

var (
	ErrBillingRefundConflict       = errors.New("billing refund accounting conflict")
	ErrBillingRefundReconciliation = errors.New("Historical billing refund status is unknown. Reconcile this trace before refunding.")
	ErrBillingRefundPending        = errors.New("Task billing is not complete. Retry the refund later.")
	ErrBillingRefundBatchUpdate    = errors.New("manual billing refunds require batch update disabled")
)

// BillingRefundReceiptProtocol identifies consumption written by a gateway
// whose trace refunds always commit a permanent primary-database receipt.
const BillingRefundReceiptProtocol = "receipt_v1"

// BillingTraceRefund is an accounting receipt, not a usage log. It lives in the
// primary database so deleting/failing the separate log database cannot make a
// successful legacy trace refund payable again.
type BillingTraceRefund struct {
	Id                  int64  `gorm:"primaryKey"`
	TraceId             string `gorm:"type:varchar(128);uniqueIndex"`
	TaskId              int64  `gorm:"index"`
	UserId              int
	TokenId             int
	ChannelId           int
	BillingSource       string `gorm:"type:varchar(32)"`
	SubscriptionId      int
	ReservationId       string `gorm:"type:varchar(64)"`
	ItemKey             string `gorm:"type:varchar(128)"`
	ConsumeQuota        int
	PreviousRefundQuota int
	Quota               int
	Source              string `gorm:"type:varchar(64)"`
	Operator            string `gorm:"type:varchar(128)"`
	Reason              string `gorm:"type:text"`
	CreatedAt           int64  `gorm:"autoCreateTime"`
}

func GetBillingTraceRefund(traceId string) (*BillingTraceRefund, error) {
	var receipt BillingTraceRefund
	err := DB.Where("trace_id = ?", strings.TrimSpace(traceId)).First(&receipt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &receipt, err
}

// RefundLoggedBillingTrace atomically credits a validated, non-task trace and
// records its permanent idempotency key. The service validates every log's owner
// and funding identity before passing the request here.
func RefundLoggedBillingTrace(ctx context.Context, input BillingTraceRefund) (*BillingTraceRefund, bool, error) {
	input.TraceId = strings.TrimSpace(input.TraceId)
	if input.TraceId == "" || len(input.TraceId) > 128 || input.UserId <= 0 ||
		input.Quota <= 0 || input.Quota > common.MaxQuota || input.ConsumeQuota < input.Quota ||
		input.ConsumeQuota > common.MaxQuota || input.PreviousRefundQuota < 0 ||
		input.ConsumeQuota-input.PreviousRefundQuota != input.Quota {
		return nil, false, ErrBillingRefundConflict
	}
	if common.BatchUpdateEnabled {
		return nil, false, ErrBillingRefundBatchUpdate
	}
	changed := false
	result := input
	var cacheChanges []billingQuotaCacheChange
	callbackCompleted := false
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (callbackErr error) {
		defer func() { callbackCompleted = callbackErr == nil }()
		// A stable owner row serializes two refunds even before the unique
		// receipt exists. SQLite's conflicting writer is rejected atomically.
		var owner User
		if err := lockForUpdate(tx).Select("id").First(&owner, input.UserId).Error; err != nil {
			return err
		}
		err := tx.Where("trace_id = ?", input.TraceId).First(&result).Error
		if err == nil {
			if result.UserId != input.UserId || result.TokenId != input.TokenId {
				return ErrBillingRefundConflict
			}
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		var tasks int64
		if err := tx.Model(&Task{}).Where("billing_trace_id = ?", input.TraceId).Count(&tasks).Error; err != nil {
			return err
		}
		if tasks != 0 {
			return ErrBillingRefundConflict
		}
		cacheChanges, err = beginBillingQuotaCaches(tx, input.UserId, input.TokenId, input.BillingSource, input.Quota)
		if err != nil {
			return err
		}
		if err := applyBillingQuotaDelta(tx, input.UserId, input.TokenId, input.ChannelId,
			input.BillingSource, input.SubscriptionId, input.ReservationId, input.ItemKey,
			-input.Quota, 0); err != nil {
			return err
		}
		result = input
		result.Id = 0
		if err := tx.Create(&result).Error; err != nil {
			return err
		}
		changed = true
		return nil
	})
	resolveBillingQuotaCaches(cacheChanges, callbackCompleted, err)
	if err != nil {
		return nil, false, err
	}
	return &result, changed, nil
}

// AdjustTaskBillingQuota serializes trace-backed settlement and refunds against
// the authoritative task row. A durable refund marker also prevents a terminal
// poller that loaded the task earlier from charging it again after a refund.
func AdjustTaskBillingQuota(ctx context.Context, taskID int64, actualQuota int, refund bool, requireTerminal bool) (*Task, int, error) {
	if taskID <= 0 || actualQuota < 0 || actualQuota > common.MaxQuota {
		return nil, 0, ErrBillingRefundConflict
	}
	if common.BatchUpdateEnabled {
		return nil, 0, ErrBillingRefundBatchUpdate
	}
	var task Task
	previous := 0
	delta := 0
	var cacheChanges []billingQuotaCacheChange
	callbackCompleted := false
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) (callbackErr error) {
		defer func() { callbackCompleted = callbackErr == nil }()
		if err := lockForUpdate(tx).First(&task, taskID).Error; err != nil {
			return err
		}
		if requireTerminal && task.Status != TaskStatusSuccess && task.Status != TaskStatusFailure {
			return ErrBillingRefundPending
		}
		if task.UserId <= 0 || task.PrivateData.TokenId < 0 || task.Quota < 0 || task.Quota > common.MaxQuota {
			return ErrBillingRefundConflict
		}
		previous = task.Quota
		traceID := strings.TrimSpace(task.BillingTraceId)
		if traceID == "" {
			traceID = strings.TrimSpace(task.PrivateData.BillingTraceId)
		}
		if traceID == "" || len(traceID) > 128 {
			return ErrBillingRefundConflict
		}
		if task.BillingTraceId != "" && task.PrivateData.BillingTraceId != "" && task.BillingTraceId != task.PrivateData.BillingTraceId {
			return ErrBillingRefundConflict
		}
		var receipt BillingTraceRefund
		receiptErr := tx.Where("trace_id = ?", traceID).First(&receipt).Error
		if receiptErr == nil {
			if receipt.TaskId != task.ID || receipt.UserId != task.UserId || receipt.TokenId != task.PrivateData.TokenId {
				return ErrBillingRefundConflict
			}
			// The durable receipt survives a deleted task, or an unrelated save
			// of an older task snapshot that accidentally loses the JSON flag.
			previous = 0
			task.Quota = 0
			task.PrivateData.BillingRefunded = true
			return tx.Model(&task).Updates(map[string]any{"quota": 0, "private_data": task.PrivateData}).Error
		}
		if !errors.Is(receiptErr, gorm.ErrRecordNotFound) {
			return receiptErr
		}
		if task.PrivateData.BillingRefunded {
			if task.Quota != 0 {
				return ErrBillingRefundConflict
			}
			return nil
		}
		if refund {
			actualQuota = 0
		}
		delta = actualQuota - task.Quota
		if delta != 0 {
			var cacheErr error
			cacheChanges, cacheErr = beginBillingQuotaCaches(tx, task.UserId, task.PrivateData.TokenId, task.PrivateData.BillingSource, -delta)
			if cacheErr != nil {
				return cacheErr
			}
			if err := applyBillingQuotaDelta(tx, task.UserId, task.PrivateData.TokenId, task.ChannelId,
				task.PrivateData.BillingSource, task.PrivateData.SubscriptionId,
				task.PrivateData.WorkflowQuotaReservationId, task.PrivateData.WorkflowQuotaItemKey,
				delta, actualQuota); err != nil {
				return err
			}
		}
		task.Quota = actualQuota
		task.PrivateData.BillingRefunded = refund
		if err := tx.Model(&task).Updates(map[string]any{
			"quota": task.Quota, "private_data": task.PrivateData,
		}).Error; err != nil {
			return err
		}
		if refund {
			source := "AUTOMATIC"
			if requireTerminal {
				source = "MANUAL"
			}
			return tx.Create(&BillingTraceRefund{
				TraceId: traceID, TaskId: task.ID, UserId: task.UserId, TokenId: task.PrivateData.TokenId,
				ChannelId: task.ChannelId, BillingSource: task.PrivateData.BillingSource,
				SubscriptionId: task.PrivateData.SubscriptionId, ReservationId: task.PrivateData.WorkflowQuotaReservationId,
				ItemKey: task.PrivateData.WorkflowQuotaItemKey, ConsumeQuota: previous, Quota: previous, Source: source,
			}).Error
		}
		return nil
	})
	resolveBillingQuotaCaches(cacheChanges, callbackCompleted, err)
	if err != nil {
		return nil, 0, err
	}
	return &task, previous, nil
}

// applyBillingQuotaDelta changes funding, token accounting and usage together.
// delta > 0 charges, delta < 0 refunds. Workflow funding changes only the node;
// its wallet and token balances remain the reservation release's responsibility.
func applyBillingQuotaDelta(tx *gorm.DB, userID, tokenID, channelID int, source string, subscriptionID int, reservationID, itemKey string, delta, actualQuota int) error {
	if delta == 0 {
		return nil
	}
	if delta < -common.MaxQuota || delta > common.MaxQuota || userID <= 0 || tokenID < 0 {
		return ErrBillingRefundConflict
	}
	switch source {
	case "workflow_reservation":
		if reservationID == "" || itemKey == "" || tokenID <= 0 {
			return ErrBillingRefundConflict
		}
		reservation, err := lockWorkflowQuotaReservation(tx, reservationID, userID, tokenID)
		if err != nil {
			return err
		}
		var item WorkflowQuotaReservationItem
		if err := lockForUpdate(tx).Where("reservation_id = ? AND item_key = ?", reservation.Id, itemKey).First(&item).Error; err != nil {
			return err
		}
		if reservation.Status != WorkflowQuotaReservationStatusReserved && reservation.Status != WorkflowQuotaReservationStatusActive {
			return ErrWorkflowQuotaReservationClosed
		}
		if item.Status != WorkflowQuotaItemStatusSettled || item.ConsumedQuota < 0 ||
			item.ConsumedQuota > common.MaxQuota || item.ConsumedQuota+delta != actualQuota || actualQuota > item.QuotedQuota {
			return ErrBillingRefundConflict
		}
		newConsumed := reservation.ConsumedQuota + delta
		if newConsumed < 0 || newConsumed > reservation.ReservedQuota {
			return ErrBillingRefundConflict
		}
		status := WorkflowQuotaItemStatusSettled
		if actualQuota == 0 {
			status = WorkflowQuotaItemStatusRefunded
		}
		if err := tx.Model(&item).Updates(map[string]any{"consumed_quota": actualQuota, "status": status}).Error; err != nil {
			return err
		}
		if err := tx.Model(reservation).Update("consumed_quota", newConsumed).Error; err != nil {
			return err
		}
	case "subscription":
		if subscriptionID <= 0 || reservationID != "" || itemKey != "" {
			return ErrBillingRefundConflict
		}
		var sub UserSubscription
		if err := lockForUpdate(tx).Where("id = ? AND user_id = ?", subscriptionID, userID).First(&sub).Error; err != nil {
			return err
		}
		if sub.AmountUsed < 0 || sub.AmountUsed > int64(common.MaxWalletQuota) ||
			sub.AmountUsed+int64(delta) < 0 || sub.AmountUsed+int64(delta) > int64(common.MaxWalletQuota) ||
			(sub.AmountTotal > 0 && sub.AmountUsed+int64(delta) > sub.AmountTotal) {
			return ErrBillingRefundConflict
		}
		if err := tx.Model(&sub).Update("amount_used", sub.AmountUsed+int64(delta)).Error; err != nil {
			return err
		}
	case "wallet", "":
		if subscriptionID != 0 || reservationID != "" || itemKey != "" {
			return ErrBillingRefundConflict
		}
		if delta < 0 {
			if err := creditTopUpQuota(tx, userID, -delta, nil); err != nil {
				return err
			}
		} else {
			result := tx.Model(&User{}).Where("id = ? AND quota >= ?", userID, delta).
				Update("quota", gorm.Expr("quota - ?", delta))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return ErrBillingRefundConflict
			}
		}
	default:
		return ErrBillingRefundConflict
	}
	if source != "workflow_reservation" && tokenID > 0 {
		var token Token
		err := lockForUpdate(tx.Unscoped()).Where("id = ?", tokenID).First(&token).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if delta > 0 {
				return err
			}
		} else {
			if token.UserId != userID || token.UsedQuota < 0 || token.UsedQuota > common.MaxWalletQuota ||
				token.UsedQuota+delta < 0 || token.UsedQuota+delta > common.MaxWalletQuota ||
				token.RemainQuota < -common.MaxWalletQuota || token.RemainQuota > common.MaxWalletQuota ||
				token.RemainQuota-delta > common.MaxWalletQuota || token.RemainQuota-delta < -common.MaxWalletQuota ||
				(delta > 0 && !token.UnlimitedQuota && token.RemainQuota < delta) {
				return ErrBillingRefundConflict
			}
			if err := tx.Unscoped().Model(&token).Updates(map[string]any{
				"remain_quota": token.RemainQuota - delta, "used_quota": token.UsedQuota + delta,
			}).Error; err != nil {
				return err
			}
		}
	}
	// Request count is historical and must not be incremented/decremented here.
	userUpdate := tx.Model(&User{}).Where("id = ? AND used_quota >= ? AND used_quota <= ?", userID, max(-delta, 0), common.MaxWalletQuota-max(delta, 0)).
		Update("used_quota", gorm.Expr("used_quota + ?", delta))
	if userUpdate.Error != nil {
		return userUpdate.Error
	}
	if userUpdate.RowsAffected != 1 {
		return ErrBillingRefundConflict
	}
	if channelID > 0 {
		var channel Channel
		err := lockForUpdate(tx).Select("id", "used_quota").First(&channel, channelID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil {
			if channel.UsedQuota < 0 || channel.UsedQuota > int64(common.MaxWalletQuota) ||
				channel.UsedQuota+int64(delta) < 0 || channel.UsedQuota+int64(delta) > int64(common.MaxWalletQuota) {
				return ErrBillingRefundConflict
			}
			return tx.Model(&channel).Update("used_quota", channel.UsedQuota+int64(delta)).Error
		}
	}
	return nil
}
