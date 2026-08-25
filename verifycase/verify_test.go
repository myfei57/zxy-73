package verifycase

import (
	"testing"

	"rubbercure/internal/steam"
)

type stopRecorder struct {
	steam        *steam.SteamSystem
	closedAtStop bool
}

func (s *stopRecorder) Stop() {
	s.closedAtStop = s.steam.Closed()
}

func TestRbcStopOrder(t *testing.T) {
	steamSystem := steam.NewSystem(nil)
	if err := steamSystem.OpenValve(); err != nil {
		t.Fatal(err)
	}
	recorder := &stopRecorder{steam: steamSystem}
	if err := steamSystem.Shutdown(recorder); err != nil {
		t.Fatal(err)
	}
	if !steamSystem.Closed() {
		t.Fatal("steam valve is still open after shutdown")
	}
	if recorder.closedAtStop {
		t.Fatal("steam valve closed before circulation pump stopped")
	}
}
