package timer

import "time"

func (t *CureTimer) Remaining(startedAt time.Time) time.Duration {
	curve := t.CurrentCurve()
	if curve.DurationSeconds <= 0 {
		return 0
	}
	elapsed := time.Since(startedAt)
	total := time.Duration(curve.DurationSeconds) * time.Second
	if elapsed >= total {
		return 0
	}
	return total - elapsed
}

func (t *CureTimer) DurationSeconds() int {
	return t.CurrentCurve().DurationSeconds
}

func (t *CureTimer) Temperature() float64 {
	return t.CurrentCurve().Temperature
}
