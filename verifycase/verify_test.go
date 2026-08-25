package verifycase

import (
	"testing"

	"rubbercure/internal/lift"
	"rubbercure/internal/mold"
)

func TestRbcMoldMapFresh(t *testing.T) {
	registry := mold.NewRegistry()
	registry.Register(mold.Mold{ID: "mold-a", Spec: "tire", SensorIDs: []string{"s1"}, MapVersion: 1})
	liftControl := lift.NewLift(2)
	closedBefore, err := liftControl.CheckClosed(registry, "mold-a")
	if err != nil {
		t.Fatal(err)
	}
	if closedBefore {
		t.Fatal("mold should be open before the sensor map refresh")
	}
	if _, err := registry.Swap("mold-a", []string{"s1", "s2"}); err != nil {
		t.Fatal(err)
	}
	closedAfter, err := liftControl.CheckClosed(registry, "mold-a")
	if err != nil {
		t.Fatal(err)
	}
	if !closedAfter {
		t.Fatal("closed state did not follow the refreshed mold map")
	}
}
