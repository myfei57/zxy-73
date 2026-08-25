package timer

import (
	"sync"
	"time"
)

type Curve struct {
	Name            string
	Temperature     float64
	DurationSeconds int
}

type CureTimer struct {
	mu        sync.Mutex
	started   bool
	running   bool
	done      bool
	latched   bool
	curve     Curve
	startedAt time.Time
	cycles    []CycleRecord
}

func NewCureTimer() *CureTimer {
	return &CureTimer{latched: true}
}

func (t *CureTimer) Running() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}

func (t *CureTimer) Done() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.done
}

func (t *CureTimer) Latched() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.latched
}
