package console

import (
	"encoding/json"
	"net/http"

	"rubbercure/internal/ns"
)

type statusPayload struct {
	Presses     []pressStatus `json:"presses"`
	Molds       []string      `json:"molds"`
	Steam       string        `json:"steam"`
	Cool        string        `json:"cool"`
	Lift        string        `json:"lift"`
	Timer       timerStatus   `json:"timer"`
	Assignments int           `json:"assignments"`
	Header      float64       `json:"header_setpoint"`
}

type pressStatus struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Mold  string `json:"mold"`
}

type timerStatus struct {
	Running bool   `json:"running"`
	Done    bool   `json:"done"`
	Latched bool   `json:"latched"`
	Curve   string `json:"curve"`
}

func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	presses := make([]pressStatus, 0)
	for _, p := range s.fleet.List() {
		presses = append(presses, pressStatus{ID: p.ID, State: string(p.State), Mold: p.MoldID})
	}
	payload := statusPayload{
		Presses:     presses,
		Molds:       s.registry.IDs(),
		Steam:       s.steam.State(),
		Cool:        s.cool.State(),
		Lift:        s.lift.State(),
		Timer:       timerStatus{Running: s.timer.Running(), Done: s.timer.Done(), Latched: s.timer.Latched(), Curve: s.timer.CurrentCurve().Name},
		Assignments: len(s.batch.Assignments()),
		Header:      s.temp.EffectiveHeader(),
	}
	s.audit.Record(ns.ComponentConsole, ns.EventType("status.read"), "console")
	writeJSON(w, http.StatusOK, payload)
}

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}
