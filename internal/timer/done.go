package timer

import "time"

func (t *CureTimer) Complete() {
	t.mu.Lock()
	t.running = false
	t.done = true
	t.recordCycleLocked(time.Now().UTC())
	t.mu.Unlock()
}
