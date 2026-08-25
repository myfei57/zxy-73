package verifycase

import (
	"errors"
	"testing"

	"rubbercure/internal/timer"
)

type failingPersister struct {
	timer      *timer.CureTimer
	seenRunning bool
	err        error
}

func (p *failingPersister) Persist() error {
	p.seenRunning = p.timer.Running()
	return p.err
}

func TestRbcCureAfterTempPersist(t *testing.T) {
	cureTimer := timer.NewCureTimer()
	persister := &failingPersister{timer: cureTimer, err: errors.New("persist failed")}
	if err := cureTimer.Start(persister); err == nil {
		t.Fatal("expected persist error")
	}
	if cureTimer.Running() {
		t.Fatal("timer started before temperature persisted")
	}
	if persister.seenRunning {
		t.Fatal("temperature persisted after the timer already started")
	}
}
