package timer

func (t *CureTimer) ResetLatch() {
	t.mu.Lock()
	t.latched = false
	t.mu.Unlock()
}
