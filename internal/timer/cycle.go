package timer

import "time"

type CycleRecord struct {
	StartedAt   time.Time
	CompletedAt time.Time
	Duration    time.Duration
	CurveName   string
}

func (t *CureTimer) CycleCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.cycles)
}

func (t *CureTimer) LastCycle() (CycleRecord, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.cycles) == 0 {
		return CycleRecord{}, false
	}
	return t.cycles[len(t.cycles)-1], true
}

func (t *CureTimer) recordCycleLocked(completedAt time.Time) {
	if t.startedAt.IsZero() {
		return
	}
	record := CycleRecord{
		StartedAt:   t.startedAt,
		CompletedAt: completedAt,
		Duration:    completedAt.Sub(t.startedAt),
		CurveName:   t.curve.Name,
	}
	t.cycles = append(t.cycles, record)
	if len(t.cycles) > 256 {
		t.cycles = append([]CycleRecord(nil), t.cycles[len(t.cycles)-256:]...)
	}
}
