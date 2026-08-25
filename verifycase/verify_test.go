package verifycase

import (
	"testing"

	"rubbercure/internal/lift"
)

type noopPurge struct{}

func (noopPurge) Purge() error {
	return nil
}

type ventRecorder struct {
	lift          *lift.LiftControl
	closedAtVent bool
}

func (v *ventRecorder) Vent() error {
	v.closedAtVent = v.lift.Closed()
	return nil
}

func TestRbcVentBeforeOpen(t *testing.T) {
	liftControl := lift.NewLift(1)
	if err := liftControl.Close(noopPurge{}); err != nil {
		t.Fatal(err)
	}
	recorder := &ventRecorder{lift: liftControl}
	if err := liftControl.Open(recorder); err != nil {
		t.Fatal(err)
	}
	if liftControl.Closed() {
		t.Fatal("mold is still closed after open")
	}
	if !recorder.closedAtVent {
		t.Fatal("mold opened before pressure vented")
	}
}
