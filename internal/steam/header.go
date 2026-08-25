package steam

import "sync"

type Header struct {
	mu        sync.Mutex
	requests  map[string]float64
	effective float64
}

func NewHeader() *Header {
	return &Header{requests: make(map[string]float64)}
}

func (h *Header) Set(pressID string, value float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.requests[pressID] = value
	h.effective = h.computeEffectiveLocked()
	return h.effective
}

func (h *Header) Effective() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.effective
}

func (h *Header) Requests() map[string]float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make(map[string]float64, len(h.requests))
	for k, v := range h.requests {
		out[k] = v
	}
	return out
}

func (h *Header) computeEffectiveLocked() float64 {
	max := 0.0
	for _, value := range h.requests {
		if value > max {
			max = value
		}
	}
	return max
}
