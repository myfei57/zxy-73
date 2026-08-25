package verifycase

import (
	"sync"
	"testing"

	"rubbercure/internal/steam"
	"rubbercure/internal/temp"
)

func TestRbcConcurrentSteamSetpoint(t *testing.T) {
	header := steam.NewHeader()
	pressA := temp.NewController("press-a", header, nil)
	pressB := temp.NewController("press-b", header, nil)
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			pressA.SetPoint(70)
		}()
		go func() {
			defer wg.Done()
			pressB.SetPoint(90)
		}()
		wg.Wait()
	}
	if header.Effective() != 90 {
		t.Fatalf("header setpoint was overwritten: got %v want 90", header.Effective())
	}
	if len(header.Requests()) != 2 {
		t.Fatalf("expected both press setpoints to be registered, got %d", len(header.Requests()))
	}
}
