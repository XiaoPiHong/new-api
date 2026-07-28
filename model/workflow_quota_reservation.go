package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// 主单状态描述整条工作流从预占、消费到释放的生命周期。
	WorkflowQuotaReservationStatusReserved = "RESERVED"
	WorkflowQuotaReservationStatusActive   = "ACTIVE"
	WorkflowQuotaReservationStatusSettled  = "SETTLED"
	WorkflowQuotaReservationStatusReleased = "RELEASED"
	WorkflowQuotaReservationStatusExpired  = "EXPIRED"

	// 明细状态约束单个节点只能认领自己的报价额度，并支持失败撤销与终态退款。
	WorkflowQuotaItemStatusReserved = "RESERVED"
	WorkflowQuotaItemStatusActive   = "ACTIVE"
	WorkflowQuotaItemStatusSettled  = "SETTLED"
	WorkflowQuotaItemStatusReleased = "RELEASED"
	WorkflowQuotaItemStatusRefunded = "REFUNDED"
)

var (
	// 模型层只返回可判定的哨兵错误，由 controller/service 映射为稳定错误码。
	ErrWorkflowQuotaReservationNotFound = errors.New("workflow quota reservation not found")
	ErrWorkflowQuotaReservationConflict = errors.New("workflow quota reservation conflict")
	ErrWorkflowQuotaReservationExpired  = errors.New("workflow quota reservation expired")
	ErrWorkflowQuotaReservationClosed   = errors.New("workflow quota reservation closed")
	ErrWorkflowQuotaItemNotFound        = errors.New("workflow quota reservation item not found")
	ErrWorkflowQuotaItemConflict        = errors.New("workflow quota reservation item conflict")
	ErrWorkflowQuotaItemExceeded        = errors.New("workflow quota reservation item quota exceeded")
	ErrWorkflowQuotaUserInsufficient    = errors.New("workflow quota user quota insufficient")
	ErrWorkflowQuotaTokenInsufficient   = errors.New("workflow quota token quota insufficient")
	ErrWorkflowQuotaBatchUpdateEnabled  = errors.New("workflow quota reservation requires batch update disabled")
)

// WorkflowQuotaReservation 是整个工作流的额度预占主单。
// 预占时直接从钱包和有限额 Token 的可用额度中扣除，任务结束后再统一释放未消费部分。
type WorkflowQuotaReservation struct {
	Id              int64  `json:"id" gorm:"primaryKey"`
	PublicId        string `json:"reservation_id" gorm:"type:varchar(64);uniqueIndex"`
	ReservationKey  string `json:"reservation_key" gorm:"type:varchar(128);uniqueIndex:idx_workflow_quota_key_owner"`
	UserId          int    `json:"user_id" gorm:"uniqueIndex:idx_workflow_quota_key_owner;index"`
	TokenId         int    `json:"token_id" gorm:"uniqueIndex:idx_workflow_quota_key_owner;index"`
	Status          string `json:"status" gorm:"type:varchar(24);index"`
	ReservedQuota   int    `json:"reserved_quota"`
	ConsumedQuota   int    `json:"consumed_quota"`
	ReleasedQuota   int    `json:"released_quota"`
	ExpiresAt       int64  `json:"expires_at" gorm:"index"`
	CreatedAt       int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int64  `json:"updated_at" gorm:"autoUpdateTime"`
	Items           []WorkflowQuotaReservationItem `json:"items" gorm:"foreignKey:ReservationId"`
}

// WorkflowQuotaReservationItem 将预占额度绑定到具体工作流节点。
// 节点不能借用其他节点的报价额度，避免前置节点挤占后续高成本节点的预算。
type WorkflowQuotaReservationItem struct {
	Id               int64  `json:"id" gorm:"primaryKey"`
	ReservationId    int64  `json:"-" gorm:"uniqueIndex:idx_workflow_quota_item"`
	ItemKey           string `json:"item_key" gorm:"type:varchar(128);uniqueIndex:idx_workflow_quota_item"`
	Adapter           string `json:"adapter" gorm:"type:varchar(32)"`
	ModelName         string `json:"model" gorm:"type:varchar(191)"`
	QuotedQuota       int    `json:"quoted_quota"`
	PreConsumedQuota  int    `json:"pre_consumed_quota"`
	ConsumedQuota     int    `json:"consumed_quota"`
	Status            string `json:"status" gorm:"type:varchar(24);index"`
	RequestId         string `json:"request_id,omitempty" gorm:"type:varchar(128)"`
	CreatedAt         int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

// WorkflowQuotaReservationItemInput 是 controller 完成报价后的可信节点额度输入。
type WorkflowQuotaReservationItemInput struct {
	ItemKey     string
	Adapter     string
	ModelName   string
	QuotedQuota int
}

// ReserveWorkflowQuotaRequest 只接收 NewAPI 自己计算出的报价，调用方不能直接指定总额度。
type ReserveWorkflowQuotaRequest struct {
	ReservationKey string
	UserId         int
	TokenId        int
	ExpiresAt      int64
	Items          []WorkflowQuotaReservationItemInput
}

// ReserveWorkflowQuota 原子完成用户钱包、Token 额度扣减和预占单写入。
// BATCH_UPDATE_ENABLED 会让数据库额度短暂落后于缓存，无法满足跨进程原子预占，因此明确拒绝。
func ReserveWorkflowQuota(input ReserveWorkflowQuotaRequest) (*WorkflowQuotaReservation, bool, error) {
	if common.BatchUpdateEnabled {
		return nil, false, ErrWorkflowQuotaBatchUpdateEnabled
	}

	key := strings.TrimSpace(input.ReservationKey)
	if key == "" || input.UserId <= 0 || input.TokenId <= 0 || len(input.Items) == 0 {
		return nil, false, ErrWorkflowQuotaReservationConflict
	}

	total := 0
	seenItems := make(map[string]struct{}, len(input.Items))
	// node.id 是后续核销主键，创建时必须唯一；总额度只能由可信报价明细汇总。
	for _, item := range input.Items {
		itemKey := strings.TrimSpace(item.ItemKey)
		if itemKey == "" || item.QuotedQuota < 0 {
			return nil, false, ErrWorkflowQuotaReservationConflict
		}
		if _, exists := seenItems[itemKey]; exists {
			return nil, false, ErrWorkflowQuotaReservationConflict
		}
		seenItems[itemKey] = struct{}{}
		total += item.QuotedQuota
	}

	var result WorkflowQuotaReservation
	created := false
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 同一用户和 Token 对同一任务的重试复用原预占单，但请求内容必须完全一致。
		var existing WorkflowQuotaReservation
		existingErr := tx.Preload("Items").Where(
			"reservation_key = ? AND user_id = ? AND token_id = ?",
			key,
			input.UserId,
			input.TokenId,
		).First(&existing).Error
		switch {
		case existingErr == nil:
			switch existing.Status {
			case WorkflowQuotaReservationStatusReserved, WorkflowQuotaReservationStatusActive:
				if !workflowQuotaReservationMatches(&existing, input.Items, total) {
					return ErrWorkflowQuotaReservationConflict
				}
				result = existing
				return nil
			default:
				return ErrWorkflowQuotaReservationConflict
			}
		case !errors.Is(existingErr, gorm.ErrRecordNotFound):
			return existingErr
		}

		// 条件更新把“检查余额”和“扣余额”合并为一条 SQL，避免并发超额预占。
		if total > 0 {
			userUpdate := tx.Model(&User{}).
				Where("id = ? AND quota >= ?", input.UserId, total).
				Update("quota", gorm.Expr("quota - ?", total))
			if userUpdate.Error != nil {
				return userUpdate.Error
			}
			if userUpdate.RowsAffected != 1 {
				return ErrWorkflowQuotaUserInsufficient
			}
		}

		// 有限额 Token 与钱包在同一事务内扣减；无限额 Token 只扣用户钱包。
		var token Token
		if err := tx.Where("id = ? AND user_id = ?", input.TokenId, input.UserId).First(&token).Error; err != nil {
			return err
		}
		if total > 0 && !token.UnlimitedQuota {
			tokenUpdate := tx.Model(&Token{}).
				Where("id = ? AND user_id = ? AND remain_quota >= ?", input.TokenId, input.UserId, total).
				Updates(map[string]interface{}{
					"remain_quota": gorm.Expr("remain_quota - ?", total),
					"used_quota":   gorm.Expr("used_quota + ?", total),
				})
			if tokenUpdate.Error != nil {
				return tokenUpdate.Error
			}
			if tokenUpdate.RowsAffected != 1 {
				return ErrWorkflowQuotaTokenInsufficient
			}
		}

		// 只有钱包、Token 都扣减成功后才写主单和节点明细，任一步失败都会整体回滚。
		result = WorkflowQuotaReservation{
			PublicId:       "wqr_" + common.GetRandomString(24),
			ReservationKey: key,
			UserId:         input.UserId,
			TokenId:        input.TokenId,
			Status:         WorkflowQuotaReservationStatusReserved,
			ReservedQuota:  total,
			ExpiresAt:      input.ExpiresAt,
		}
		if err := tx.Create(&result).Error; err != nil {
			return err
		}

		items := make([]WorkflowQuotaReservationItem, 0, len(input.Items))
		for _, item := range input.Items {
			items = append(items, WorkflowQuotaReservationItem{
				ReservationId: result.Id,
				ItemKey:        strings.TrimSpace(item.ItemKey),
				Adapter:        strings.TrimSpace(item.Adapter),
				ModelName:      strings.TrimSpace(item.ModelName),
				QuotedQuota:    item.QuotedQuota,
				Status:         WorkflowQuotaItemStatusReserved,
			})
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		result.Items = items
		created = true
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	// Redis 更新必须发生在事务提交后，数据库始终作为额度事实来源。
	if created && total > 0 {
		adjustWorkflowQuotaCaches(input.UserId, input.TokenId, -total)
	}
	return &result, created, nil
}

// workflowQuotaReservationMatches 保证同一个幂等键只能代表同一份节点报价。
// 调用方重试时若模型、adapter 或报价发生变化，必须拒绝复用旧预占单。
func workflowQuotaReservationMatches(
	existing *WorkflowQuotaReservation,
	items []WorkflowQuotaReservationItemInput,
	total int,
) bool {
	if existing == nil ||
		existing.ReservedQuota != total ||
		len(existing.Items) != len(items) {
		return false
	}
	existingByKey := make(map[string]WorkflowQuotaReservationItem, len(existing.Items))
	for _, item := range existing.Items {
		existingByKey[item.ItemKey] = item
	}
	for _, item := range items {
		stored, ok := existingByKey[strings.TrimSpace(item.ItemKey)]
		if !ok || stored.Adapter != strings.TrimSpace(item.Adapter) ||
			stored.ModelName != strings.TrimSpace(item.ModelName) || stored.QuotedQuota != item.QuotedQuota {
			return false
		}
	}
	return true
}

// ClaimWorkflowQuotaItem 在请求发送上游前锁定对应节点额度。
// 同一 requestId 可以幂等重试，不同请求不能重复占用已经 ACTIVE 的节点。
func ClaimWorkflowQuotaItem(publicId string, itemKey string, userId int, tokenId int, quota int, requestId string) error {
	if quota < 0 {
		return ErrWorkflowQuotaItemExceeded
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		reservation, err := lockWorkflowQuotaReservation(tx, publicId, userId, tokenId)
		if err != nil {
			return err
		}
		// 过期后禁止启动新上游请求，未消费额度交给定时回收任务释放。
		if reservation.ExpiresAt > 0 && reservation.ExpiresAt <= common.GetTimestamp() {
			return ErrWorkflowQuotaReservationExpired
		}
		switch reservation.Status {
		case WorkflowQuotaReservationStatusReserved, WorkflowQuotaReservationStatusActive:
		default:
			return ErrWorkflowQuotaReservationClosed
		}

		var item WorkflowQuotaReservationItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
			"reservation_id = ? AND item_key = ?",
			reservation.Id,
			strings.TrimSpace(itemKey),
		).First(&item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWorkflowQuotaItemNotFound
			}
			return err
		}

		switch item.Status {
		case WorkflowQuotaItemStatusReserved:
			if quota > item.QuotedQuota {
				return fmt.Errorf("%w: quoted=%d required=%d", ErrWorkflowQuotaItemExceeded, item.QuotedQuota, quota)
			}
			if err := tx.Model(&item).Updates(map[string]interface{}{
				"status":             WorkflowQuotaItemStatusActive,
				"pre_consumed_quota": quota,
				"request_id":          strings.TrimSpace(requestId),
			}).Error; err != nil {
				return err
			}
			return tx.Model(reservation).Update("status", WorkflowQuotaReservationStatusActive).Error
		case WorkflowQuotaItemStatusActive:
			if item.RequestId == strings.TrimSpace(requestId) {
				if quota > item.QuotedQuota {
					return fmt.Errorf("%w: quoted=%d required=%d", ErrWorkflowQuotaItemExceeded, item.QuotedQuota, quota)
				}
				if item.PreConsumedQuota == quota {
					return nil
				}
				return tx.Model(&item).Update("pre_consumed_quota", quota).Error
			}
			return ErrWorkflowQuotaItemConflict
		default:
			return ErrWorkflowQuotaItemConflict
		}
	})
}

// AdjustSettledWorkflowQuotaItem 供异步任务轮询阶段修正实际费用或执行退款。
// 它只调整预占单内部的已消费金额，钱包和 Token 的最终返还仍由释放接口统一完成。
func AdjustSettledWorkflowQuotaItem(publicId string, itemKey string, userId int, tokenId int, actualQuota int) error {
	if actualQuota < 0 {
		return ErrWorkflowQuotaItemExceeded
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		reservation, err := lockWorkflowQuotaReservation(tx, publicId, userId, tokenId)
		if err != nil {
			return err
		}
		switch reservation.Status {
		case WorkflowQuotaReservationStatusReserved, WorkflowQuotaReservationStatusActive:
		default:
			return ErrWorkflowQuotaReservationClosed
		}

		var item WorkflowQuotaReservationItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
			"reservation_id = ? AND item_key = ?",
			reservation.Id,
			strings.TrimSpace(itemKey),
		).First(&item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWorkflowQuotaItemNotFound
			}
			return err
		}
		switch item.Status {
		case WorkflowQuotaItemStatusSettled, WorkflowQuotaItemStatusRefunded:
		default:
			return ErrWorkflowQuotaItemConflict
		}
		if actualQuota > item.QuotedQuota {
			return fmt.Errorf("%w: quoted=%d actual=%d", ErrWorkflowQuotaItemExceeded, item.QuotedQuota, actualQuota)
		}
		// 使用“主单原消费 - 明细原消费 + 明细新消费”保证异步差额调整可重复执行。
		newConsumed := reservation.ConsumedQuota - item.ConsumedQuota + actualQuota
		if newConsumed < 0 || newConsumed > reservation.ReservedQuota {
			return ErrWorkflowQuotaItemExceeded
		}
		itemStatus := WorkflowQuotaItemStatusSettled
		if actualQuota == 0 {
			itemStatus = WorkflowQuotaItemStatusRefunded
		}
		if err := tx.Model(&item).Updates(map[string]interface{}{
			"status":         itemStatus,
			"consumed_quota": actualQuota,
		}).Error; err != nil {
			return err
		}
		return tx.Model(reservation).Update("consumed_quota", newConsumed).Error
	})
}

// SettleWorkflowQuotaItem 把节点实际消费写入预占单，不再重复修改钱包或 Token 额度。
func SettleWorkflowQuotaItem(publicId string, itemKey string, userId int, tokenId int, actualQuota int) error {
	if actualQuota < 0 {
		return ErrWorkflowQuotaItemExceeded
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		reservation, err := lockWorkflowQuotaReservation(tx, publicId, userId, tokenId)
		if err != nil {
			return err
		}
		switch reservation.Status {
		case WorkflowQuotaReservationStatusReserved, WorkflowQuotaReservationStatusActive:
		default:
			return ErrWorkflowQuotaReservationClosed
		}

		var item WorkflowQuotaReservationItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
			"reservation_id = ? AND item_key = ?",
			reservation.Id,
			strings.TrimSpace(itemKey),
		).First(&item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWorkflowQuotaItemNotFound
			}
			return err
		}
		switch item.Status {
		case WorkflowQuotaItemStatusSettled:
			if item.ConsumedQuota == actualQuota {
				return nil
			}
			return ErrWorkflowQuotaItemConflict
		case WorkflowQuotaItemStatusActive:
		default:
			return ErrWorkflowQuotaItemConflict
		}
		if actualQuota > item.QuotedQuota {
			return fmt.Errorf("%w: quoted=%d actual=%d", ErrWorkflowQuotaItemExceeded, item.QuotedQuota, actualQuota)
		}
		// 首次结算只把当前节点实际费用累加到主单，不能超过整单预占总额。
		newConsumed := reservation.ConsumedQuota + actualQuota
		if newConsumed > reservation.ReservedQuota {
			return ErrWorkflowQuotaItemExceeded
		}
		if err := tx.Model(&item).Updates(map[string]interface{}{
			"status":         WorkflowQuotaItemStatusSettled,
			"consumed_quota": actualQuota,
		}).Error; err != nil {
			return err
		}
		return tx.Model(reservation).Update("consumed_quota", newConsumed).Error
	})
}

// RefundWorkflowQuotaItem 只撤销本次失败请求的节点占用，钱包退款仍由预占单释放接口统一处理。
func RefundWorkflowQuotaItem(publicId string, itemKey string, userId int, tokenId int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		reservation, err := lockWorkflowQuotaReservation(tx, publicId, userId, tokenId)
		if err != nil {
			if errors.Is(err, ErrWorkflowQuotaReservationClosed) {
				return nil
			}
			return err
		}
		var item WorkflowQuotaReservationItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
			"reservation_id = ? AND item_key = ?",
			reservation.Id,
			strings.TrimSpace(itemKey),
		).First(&item).Error; err != nil {
			return err
		}
		switch item.Status {
		case WorkflowQuotaItemStatusActive:
			return tx.Model(&item).Updates(map[string]interface{}{
				"status":             WorkflowQuotaItemStatusReserved,
				"pre_consumed_quota": 0,
				"request_id":          "",
			}).Error
		case WorkflowQuotaItemStatusReserved, WorkflowQuotaItemStatusReleased, WorkflowQuotaItemStatusRefunded:
			return nil
		default:
			return ErrWorkflowQuotaItemConflict
		}
	})
}

// ReleaseWorkflowQuotaReservation 结束预占单。
// refundConsumed=false 只返还未消费额度；失败任务传 true，按现有业务策略返还整单额度。
func ReleaseWorkflowQuotaReservation(
	publicId string,
	userId int,
	tokenId int,
	refundConsumed bool,
) (*WorkflowQuotaReservation, bool, error) {
	var result WorkflowQuotaReservation
	refundedQuota := 0
	changed := false
	err := DB.Transaction(func(tx *gorm.DB) error {
		reservation, err := lockWorkflowQuotaReservation(tx, publicId, userId, tokenId)
		if err != nil {
			return err
		}
		switch reservation.Status {
		case WorkflowQuotaReservationStatusReleased:
			result = *reservation
			return tx.Preload("Items").First(&result, reservation.Id).Error
		case WorkflowQuotaReservationStatusSettled, WorkflowQuotaReservationStatusExpired:
			// 成功释放后，XphAI 仍可能在最终结果落库时失败；这时允许把已结算单
			// 幂等升级为整单退款。过期单也允许任务超时处理补退已消费部分。
			if !refundConsumed {
				result = *reservation
				return tx.Preload("Items").First(&result, reservation.Id).Error
			}
		case WorkflowQuotaReservationStatusReserved, WorkflowQuotaReservationStatusActive:
		default:
			return ErrWorkflowQuotaReservationClosed
		}

		// released_quota 记录此前已经返还的额度；重试只补差额，避免重复退款。
		targetReleasedQuota := reservation.ReservedQuota - reservation.ConsumedQuota
		if refundConsumed {
			targetReleasedQuota = reservation.ReservedQuota
		}
		refundedQuota = targetReleasedQuota - reservation.ReleasedQuota
		if refundedQuota < 0 {
			return ErrWorkflowQuotaReservationConflict
		}
		if refundedQuota > 0 {
			if err := tx.Model(&User{}).Where("id = ?", userId).
				Update("quota", gorm.Expr("quota + ?", refundedQuota)).Error; err != nil {
				return err
			}
			var token Token
			if err := tx.Where("id = ? AND user_id = ?", tokenId, userId).First(&token).Error; err != nil {
				return err
			}
			if !token.UnlimitedQuota {
				if err := tx.Model(&Token{}).Where("id = ?", tokenId).Updates(map[string]interface{}{
					"remain_quota": gorm.Expr("remain_quota + ?", refundedQuota),
					"used_quota":   gorm.Expr("used_quota - ?", refundedQuota),
				}).Error; err != nil {
					return err
				}
			}
		}

		// 成功且有真实消费记为 SETTLED；无消费或失败整单退款记为 RELEASED。
		status := WorkflowQuotaReservationStatusSettled
		itemUnusedStatus := WorkflowQuotaItemStatusReleased
		if refundConsumed || reservation.ConsumedQuota == 0 {
			status = WorkflowQuotaReservationStatusReleased
			itemUnusedStatus = WorkflowQuotaItemStatusRefunded
		}
		if err := tx.Model(&WorkflowQuotaReservationItem{}).
			Where("reservation_id = ? AND status <> ?", reservation.Id, WorkflowQuotaItemStatusSettled).
			Update("status", itemUnusedStatus).Error; err != nil {
			return err
		}
		if refundConsumed {
			if err := tx.Model(&WorkflowQuotaReservationItem{}).
				Where("reservation_id = ?", reservation.Id).
				Update("status", WorkflowQuotaItemStatusRefunded).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(reservation).Updates(map[string]interface{}{
			"status":         status,
			"released_quota": targetReleasedQuota,
		}).Error; err != nil {
			return err
		}
		changed = true
		return tx.Preload("Items").First(&result, reservation.Id).Error
	})
	if err != nil {
		return nil, false, err
	}
	if changed && refundedQuota > 0 {
		adjustWorkflowQuotaCaches(userId, tokenId, refundedQuota)
	}
	return &result, changed, nil
}

// ReleaseWorkflowQuotaReservationInternal 供持有服务间共享密钥的 SSO 释放预占单。
// 终态退款不能依赖用户 Token 仍然有效，因此先按不可猜测的 public_id 读取所有者，
// 再复用同一个带行锁和差额幂等保护的释放实现。
func ReleaseWorkflowQuotaReservationInternal(
	publicId string,
	refundConsumed bool,
) (*WorkflowQuotaReservation, bool, error) {
	var reservation WorkflowQuotaReservation
	err := DB.Select("user_id", "token_id").Where(
		"public_id = ?",
		strings.TrimSpace(publicId),
	).First(&reservation).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, ErrWorkflowQuotaReservationNotFound
	}
	if err != nil {
		return nil, false, err
	}
	return ReleaseWorkflowQuotaReservation(
		publicId,
		reservation.UserId,
		reservation.TokenId,
		refundConsumed,
	)
}

// ReleaseExpiredWorkflowQuotaReservations 只释放过期单的未消费额度，避免误退已经成功提交上游的节点费用。
func ReleaseExpiredWorkflowQuotaReservations(limit int) (int, error) {
	if limit <= 0 {
		limit = 100
	}
	var reservations []WorkflowQuotaReservation
	if err := DB.Where(
		"status IN ? AND expires_at > 0 AND expires_at <= ?",
		[]string{WorkflowQuotaReservationStatusReserved, WorkflowQuotaReservationStatusActive},
		common.GetTimestamp(),
	).Order("id ASC").Limit(limit).Find(&reservations).Error; err != nil {
		return 0, err
	}
	released := 0
	for _, reservation := range reservations {
		result, changed, err := releaseExpiredWorkflowQuotaReservation(reservation.PublicId, reservation.UserId, reservation.TokenId)
		if err != nil {
			common.SysLog(fmt.Sprintf("release expired workflow quota reservation %s failed: %s", reservation.PublicId, err.Error()))
			continue
		}
		if changed {
			released++
			adjustWorkflowQuotaCaches(result.UserId, result.TokenId, result.ReleasedQuota)
		}
	}
	return released, nil
}

// releaseExpiredWorkflowQuotaReservation 用行锁回收单张过期预占单，重复调用不会重复退款。
func releaseExpiredWorkflowQuotaReservation(publicId string, userId int, tokenId int) (*WorkflowQuotaReservation, bool, error) {
	var result WorkflowQuotaReservation
	refundedQuota := 0
	changed := false
	err := DB.Transaction(func(tx *gorm.DB) error {
		reservation, err := lockWorkflowQuotaReservation(tx, publicId, userId, tokenId)
		if err != nil {
			return err
		}
		switch reservation.Status {
		case WorkflowQuotaReservationStatusSettled, WorkflowQuotaReservationStatusReleased, WorkflowQuotaReservationStatusExpired:
			result = *reservation
			return nil
		case WorkflowQuotaReservationStatusReserved, WorkflowQuotaReservationStatusActive:
		default:
			return ErrWorkflowQuotaReservationClosed
		}
		// 自动过期只退未消费差额，已结算节点是否退款由 SSO 的任务终态决定。
		refundedQuota = reservation.ReservedQuota - reservation.ConsumedQuota
		if refundedQuota < 0 {
			return ErrWorkflowQuotaReservationConflict
		}
		if refundedQuota > 0 {
			if err := tx.Model(&User{}).Where("id = ?", userId).
				Update("quota", gorm.Expr("quota + ?", refundedQuota)).Error; err != nil {
				return err
			}
			var token Token
			if err := tx.Where("id = ? AND user_id = ?", tokenId, userId).First(&token).Error; err != nil {
				return err
			}
			if !token.UnlimitedQuota {
				if err := tx.Model(&Token{}).Where("id = ?", tokenId).Updates(map[string]interface{}{
					"remain_quota": gorm.Expr("remain_quota + ?", refundedQuota),
					"used_quota":   gorm.Expr("used_quota - ?", refundedQuota),
				}).Error; err != nil {
					return err
				}
			}
		}
		if err := tx.Model(&WorkflowQuotaReservationItem{}).
			Where("reservation_id = ? AND status <> ?", reservation.Id, WorkflowQuotaItemStatusSettled).
			Update("status", WorkflowQuotaItemStatusReleased).Error; err != nil {
			return err
		}
		if err := tx.Model(reservation).Updates(map[string]interface{}{
			"status":         WorkflowQuotaReservationStatusExpired,
			"released_quota": refundedQuota,
		}).Error; err != nil {
			return err
		}
		changed = true
		return tx.First(&result, reservation.Id).Error
	})
	return &result, changed, err
}

// lockWorkflowQuotaReservation 同时校验预占单所有者并加行锁，串行化认领、结算和释放。
func lockWorkflowQuotaReservation(tx *gorm.DB, publicId string, userId int, tokenId int) (*WorkflowQuotaReservation, error) {
	var reservation WorkflowQuotaReservation
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
		"public_id = ? AND user_id = ? AND token_id = ?",
		strings.TrimSpace(publicId),
		userId,
		tokenId,
	).First(&reservation).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrWorkflowQuotaReservationNotFound
	}
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// adjustWorkflowQuotaCaches 在数据库事务提交后同步 Redis 中的可用额度。
// 数据库仍是事实来源；缓存失败只记录日志，后续缓存过期会自动回源。
func adjustWorkflowQuotaCaches(userId int, tokenId int, delta int) {
	if delta == 0 || !common.RedisEnabled {
		return
	}
	if err := cacheIncrUserQuota(userId, int64(delta)); err != nil {
		common.SysLog("adjust workflow reservation user quota cache failed: " + err.Error())
	}
	var token Token
	if err := DB.Select("id", "key", "unlimited_quota").First(&token, tokenId).Error; err != nil {
		common.SysLog("load workflow reservation token cache key failed: " + err.Error())
		return
	}
	if token.UnlimitedQuota {
		return
	}
	if err := cacheIncrTokenQuota(token.Key, int64(delta)); err != nil {
		common.SysLog("adjust workflow reservation token quota cache failed: " + err.Error())
	}
}
