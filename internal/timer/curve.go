package timer

func (t *CureTimer) SetCurve(curve Curve) {
	t.mu.Lock()
	t.curve = curve
	t.mu.Unlock()
}

func (t *CureTimer) CurrentCurve() Curve {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.curve
}
