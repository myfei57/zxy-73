package verifycase

import (
	"testing"

	"rubbercure/internal/temp"
)

type warmupOrderRecorder struct {
	order []string
}

func (w *warmupOrderRecorder) Trap() error {
	w.order = append(w.order, "trap")
	return nil
}

func (w *warmupOrderRecorder) OpenValve() error {
	w.order = append(w.order, "valve")
	return nil
}

func TestRbcSteamVentOrder(t *testing.T) {
	controller := temp.NewController("press-a", nil, nil)
	recorder := &warmupOrderRecorder{}
	if err := controller.Warmup(recorder); err != nil {
		t.Fatal(err)
	}
	if len(recorder.order) != 2 || recorder.order[0] != "trap" || recorder.order[1] != "valve" {
		t.Fatalf("wrong warmup order: %v", recorder.order)
	}
}
