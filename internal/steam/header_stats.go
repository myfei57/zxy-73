package steam

func (h *Header) Max() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	max := 0.0
	for _, value := range h.requests {
		if value > max {
			max = value
		}
	}
	return max
}

func (h *Header) Min() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.requests) == 0 {
		return 0
	}
	min := 0.0
	first := true
	for _, value := range h.requests {
		if first || value < min {
			min = value
			first = false
		}
	}
	return min
}

func (h *Header) Reset(pressID string) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.requests, pressID)
	h.effective = h.computeEffectiveLocked()
	return h.effective
}
