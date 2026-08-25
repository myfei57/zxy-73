package verifycase

import (
	"testing"

	"rubbercure/internal/audit"
	"rubbercure/internal/batch"
	"rubbercure/internal/mold"
	"rubbercure/internal/press"
	"rubbercure/internal/timer"
)

func TestRbcTimerLatchReset(t *testing.T) {
	fleet := press.NewFleet()
	fleet.Add(&press.Press{ID: "press-a", State: press.StateCuring})
	registry := mold.NewRegistry()
	cureTimer := timer.NewCureTimer()
	recorder := audit.NewRecorder(nil)
	scheduler := batch.NewScheduler(fleet, registry, cureTimer, recorder)
	if err := scheduler.ReportDone("press-a"); err != nil {
		t.Fatal(err)
	}
	if cureTimer.Latched() {
		t.Fatal("cure timer latch was not released after the cycle")
	}
}
