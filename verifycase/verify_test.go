package verifycase

import (
	"testing"

	"rubbercure/internal/lift"
)

type purgeRecorder struct {
	lift        *lift.LiftControl
	closedAtPurge bool
}

func (p *purgeRecorder) Purge() error {
	p.closedAtPurge = p.lift.Closed()
	return nil
}

func TestRbcPressCloseOrder(t *testing.T) {
	liftControl := lift.NewLift(1)
	recorder := &purgeRecorder{lift: liftControl}
	if err := liftControl.Close(recorder); err != nil {
		t.Fatal(err)
	}
	if !liftControl.Closed() {
		t.Fatal("press did not close")
	}
	if recorder.closedAtPurge {
		t.Fatal("press closed before air purge")
	}
}
