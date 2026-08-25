package timer

import "time"

type TempPersister interface {
	Persist() error
}

func (t *CureTimer) Start(persister TempPersister) error {
	if err := persister.Persist(); err != nil {
		return err
	}
	t.mu.Lock()
	t.started = true
	t.running = true
	t.done = false
	t.startedAt = time.Now().UTC()
	t.mu.Unlock()
	return nil
}
