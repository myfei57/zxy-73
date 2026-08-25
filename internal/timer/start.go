package timer

import "time"

type TempPersister interface {
	Persist() error
}

// Start 持久化模温读数后再翻转计时状态。
// 顺序必须是“先落盘、后启动”：只有模温成功落盘，计时启动这一动作
// 才意味着恢复时能取到与计时配套的模温记录。若先置 started/running
// 再落盘，一旦在 Persist 之前断电，恢复后系统只见计时标志、不见模温，
// 只能沿用旧温度继续计时，导致硫化不足。
func (t *CureTimer) Start(persister TempPersister) error {
	if persister != nil {
		if err := persister.Persist(); err != nil {
			return err
		}
	}
	t.mu.Lock()
	t.started = true
	t.running = true
	t.done = false
	t.startedAt = time.Now().UTC()
	t.mu.Unlock()
	return nil
}
