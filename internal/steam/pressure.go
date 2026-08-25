package steam

import "sync"

type Pressure struct {
	mu  sync.Mutex
	bar float64
}

func NewPressure() *Pressure {
	return &Pressure{}
}

func (p *Pressure) Set(bar float64) {
	p.mu.Lock()
	p.bar = bar
	p.mu.Unlock()
}

func (p *Pressure) Release() {
	p.mu.Lock()
	p.bar = 0
	p.mu.Unlock()
}

func (p *Pressure) Current() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.bar
}

func (p *Pressure) OverLimit(limit float64) bool {
	return p.Current() > limit
}
