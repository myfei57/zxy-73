package timer

import "time"

func (t *CureTimer) Complete() {
	t.mu.Lock()
	t.running = false
	t.done = true
	// Release the timing latch so the next cycle re-enters timing;
	// without this reset the latch stays engaged and the next mold's
	// cure is skipped, leaving the product undercured.
	t.latched = false
	t.recordCycleLocked(time.Now().UTC())
	t.mu.Unlock()
}
