package timer

func (t *CureTimer) Started() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.started
}
