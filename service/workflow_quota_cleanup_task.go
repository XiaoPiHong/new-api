package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

var workflowQuotaCleanupTaskOnce sync.Once

// StartWorkflowQuotaReservationCleanupTask 定时释放进程异常后遗留的过期预占额度。
// 过期回收只返还未消费额度，已经提交上游的节点费用仍保留。
func StartWorkflowQuotaReservationCleanupTask() {
	workflowQuotaCleanupTaskOnce.Do(func() {
		if !common.IsMasterNode {
			return
		}
		go func() {
			runWorkflowQuotaReservationCleanup()
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				runWorkflowQuotaReservationCleanup()
			}
		}()
	})
}

// runWorkflowQuotaReservationCleanup 每批最多处理 100 张，避免定时任务长时间占用数据库。
func runWorkflowQuotaReservationCleanup() {
	released, err := model.ReleaseExpiredWorkflowQuotaReservations(100)
	if err != nil {
		common.SysLog("workflow quota reservation cleanup failed: " + err.Error())
		return
	}
	if released > 0 {
		common.SysLog(fmt.Sprintf("released %d expired workflow quota reservations", released))
	}
}
