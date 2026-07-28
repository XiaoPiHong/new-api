package service

import (
	"time"

	"github.com/QuantumNous/new-api/model"
)

// ---------------------------------------------------------------------------
// FundingSource — 资金来源接口（钱包 or 订阅）
// ---------------------------------------------------------------------------

// FundingSource 抽象了预扣费的资金来源。
type FundingSource interface {
	// Source 返回资金来源标识：钱包、订阅或工作流额度预占。
	Source() string
	// PreConsume 从该资金来源预扣 amount 额度
	PreConsume(amount int) error
	// Settle 根据差额调整资金来源（正数补扣，负数退还）
	Settle(delta int) error
	// Refund 退还所有预扣费
	Refund() error
}

// ---------------------------------------------------------------------------
// WorkflowReservationFunding — 工作流预占资金来源
// ---------------------------------------------------------------------------

// WorkflowReservationFunding 不直接增减钱包额度。
// 整个工作流的额度已经在任务开始前一次性扣除，这里只负责节点认领和实际消费登记。
type WorkflowReservationFunding struct {
	reservationId string
	itemKey        string
	userId         int
	tokenId        int
	requestId      string
	preConsumed    int
	claimed        bool
}

func (w *WorkflowReservationFunding) Source() string { return BillingSourceWorkflowReservation }

// PreConsume 在调用上游前认领节点预算，相同 requestId 重试时保持幂等。
func (w *WorkflowReservationFunding) PreConsume(amount int) error {
	if err := model.ClaimWorkflowQuotaItem(
		w.reservationId,
		w.itemKey,
		w.userId,
		w.tokenId,
		amount,
		w.requestId,
	); err != nil {
		return err
	}
	w.preConsumed = amount
	w.claimed = true
	return nil
}

// Settle 登记节点最终消费；钱包余额已预扣，因此这里只更新预占主单和明细。
func (w *WorkflowReservationFunding) Settle(delta int) error {
	actualQuota := w.preConsumed + delta
	if err := model.SettleWorkflowQuotaItem(
		w.reservationId,
		w.itemKey,
		w.userId,
		w.tokenId,
		actualQuota,
	); err != nil {
		return err
	}
	w.claimed = false
	return nil
}

// Refund 撤销尚未结算的节点认领，等待任务终态决定是否整单退款。
func (w *WorkflowReservationFunding) Refund() error {
	if !w.claimed {
		return nil
	}
	if err := model.RefundWorkflowQuotaItem(w.reservationId, w.itemKey, w.userId, w.tokenId); err != nil {
		return err
	}
	w.claimed = false
	return nil
}

// Reserve 将图片等响应后才能确认的预扣目标限制在该节点报价额度内。
func (w *WorkflowReservationFunding) Reserve(delta int) error {
	if delta <= 0 {
		return nil
	}
	targetQuota := w.preConsumed + delta
	if err := model.ClaimWorkflowQuotaItem(
		w.reservationId,
		w.itemKey,
		w.userId,
		w.tokenId,
		targetQuota,
		w.requestId,
	); err != nil {
		return err
	}
	w.preConsumed = targetQuota
	return nil
}

// ---------------------------------------------------------------------------
// WalletFunding — 钱包资金来源实现
// ---------------------------------------------------------------------------

type WalletFunding struct {
	userId   int
	consumed int // 实际预扣的用户额度
}

func (w *WalletFunding) Source() string { return BillingSourceWallet }

func (w *WalletFunding) PreConsume(amount int) error {
	if amount <= 0 {
		return nil
	}
	if err := model.DecreaseUserQuota(w.userId, amount, false); err != nil {
		return err
	}
	w.consumed = amount
	return nil
}

func (w *WalletFunding) Settle(delta int) error {
	if delta == 0 {
		return nil
	}
	if delta > 0 {
		return model.DecreaseUserQuota(w.userId, delta, false)
	}
	return model.IncreaseUserQuota(w.userId, -delta, false)
}

func (w *WalletFunding) Refund() error {
	if w.consumed <= 0 {
		return nil
	}
	// IncreaseUserQuota 是 quota += N 的非幂等操作，不能重试，否则会多退额度。
	// 订阅的 RefundSubscriptionPreConsume 有 requestId 幂等保护所以可以重试。
	return model.IncreaseUserQuota(w.userId, w.consumed, false)
}

// ---------------------------------------------------------------------------
// SubscriptionFunding — 订阅资金来源实现
// ---------------------------------------------------------------------------

type SubscriptionFunding struct {
	requestId      string
	userId         int
	modelName      string
	amount         int64 // 预扣的订阅额度（subConsume）
	subscriptionId int
	preConsumed    int64
	// 以下字段在 PreConsume 成功后填充，供 RelayInfo 同步使用
	AmountTotal     int64
	AmountUsedAfter int64
	PlanId          int
	PlanTitle       string
}

func (s *SubscriptionFunding) Source() string { return BillingSourceSubscription }

func (s *SubscriptionFunding) PreConsume(_ int) error {
	// amount 参数被忽略，使用内部 s.amount（已在构造时根据 preConsumedQuota 计算）
	res, err := model.PreConsumeUserSubscription(s.requestId, s.userId, s.modelName, 0, s.amount)
	if err != nil {
		return err
	}
	s.subscriptionId = res.UserSubscriptionId
	s.preConsumed = res.PreConsumed
	s.AmountTotal = res.AmountTotal
	s.AmountUsedAfter = res.AmountUsedAfter
	// 获取订阅计划信息
	if planInfo, err := model.GetSubscriptionPlanInfoByUserSubscriptionId(res.UserSubscriptionId); err == nil && planInfo != nil {
		s.PlanId = planInfo.PlanId
		s.PlanTitle = planInfo.PlanTitle
	}
	return nil
}

func (s *SubscriptionFunding) Settle(delta int) error {
	if delta == 0 {
		return nil
	}
	return model.PostConsumeUserSubscriptionDelta(s.subscriptionId, int64(delta))
}

func (s *SubscriptionFunding) Refund() error {
	if s.preConsumed <= 0 {
		return nil
	}
	return refundWithRetry(func() error {
		return model.RefundSubscriptionPreConsume(s.requestId)
	})
}

// refundWithRetry 尝试多次执行退款操作以提高成功率，只能用于基于事务的退款函数！！！！！！
// try to refund with retries, only for refund functions based on transactions!!!
func refundWithRetry(fn func() error) error {
	if fn == nil {
		return nil
	}
	const maxAttempts = 3
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if i < maxAttempts-1 {
			time.Sleep(time.Duration(200*(i+1)) * time.Millisecond)
		}
	}
	return lastErr
}
