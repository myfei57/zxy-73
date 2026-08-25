package verifycase

import (
	"testing"

	"rubbercure/internal/audit"
	"rubbercure/internal/batch"
	"rubbercure/internal/mold"
	"rubbercure/internal/press"
	"rubbercure/internal/timer"
)

func TestRbcCureCurveBaselineFresh(t *testing.T) {
	cureTimer := timer.NewCureTimer()
	cureTimer.SetCurve(timer.Curve{Name: "old", Temperature: 100, DurationSeconds: 60})
	fleet := press.NewFleet()
	registry := mold.NewRegistry()
	recorder := audit.NewRecorder(nil)
	scheduler := batch.NewScheduler(fleet, registry, cureTimer, recorder)
	scheduler.ProfileUpdate(timer.Curve{Name: "new", Temperature: 120, DurationSeconds: 90})
	if cureTimer.CurrentCurve().Name != "new" {
		t.Fatalf("cure curve was not refreshed: got %q", cureTimer.CurrentCurve().Name)
	}
}
