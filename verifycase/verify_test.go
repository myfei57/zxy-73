package verifycase

import (
	"sync"
	"testing"

	"rubbercure/internal/audit"
	"rubbercure/internal/batch"
	"rubbercure/internal/mold"
	"rubbercure/internal/press"
	"rubbercure/internal/timer"
)

func TestRbcConcurrentDoneReport(t *testing.T) {
	fleet := press.NewFleet()
	fleet.Add(&press.Press{ID: "press-a", State: press.StateCuring})
	fleet.Add(&press.Press{ID: "press-b", State: press.StateCuring})
	registry := mold.NewRegistry()
	cureTimer := timer.NewCureTimer()
	recorder := audit.NewRecorder(nil)
	scheduler := batch.NewScheduler(fleet, registry, cureTimer, recorder)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = scheduler.ReportDone("press-a")
		}()
		go func() {
			defer wg.Done()
			_ = scheduler.ReportDone("press-b")
		}()
		wg.Wait()
	}
	seen := make(map[string]bool)
	for _, assignment := range scheduler.Assignments() {
		if seen[assignment.MoldID] {
			t.Fatalf("duplicate mold assignment: %s", assignment.MoldID)
		}
		seen[assignment.MoldID] = true
	}
}
