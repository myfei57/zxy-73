package timer

import "time"

type TempPersister interface {
	Persist() error
}

func (t *CureTimer) Start(persister TempPersister) error {
	t.mu.Lock()
	t.started = true
	t.running = true
	t.done = false
	t.startedAt = time.Now().UTC()
	t.mu.Unlock()
	return persister.Persist()
}
