package console

import (
	"encoding/json"
	"net/http"

	"rubbercure/internal/ns"
)

type cureRequest struct {
	PressID     string  `json:"press_id"`
	MoldID      string  `json:"mold_id"`
	Temperature float64 `json:"temperature"`
}

type cureResponse struct {
	Status      string  `json:"status"`
	Header      float64 `json:"header_setpoint"`
	TimerStart  bool    `json:"timer_started"`
	LiftState   string  `json:"lift_state"`
	MoldVersion int     `json:"mold_version"`
}

func (s *Server) handleCure(w http.ResponseWriter, r *http.Request) {
	var req cureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	current, err := s.registry.CurrentMap(req.MoldID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.temp.SetPoint(req.Temperature)
	if err := s.lift.Close(s.steam); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.temp.Warmup(s.steam); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.timer.Start(s.temp); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.batch.AssignMold(req.PressID, req.MoldID); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	s.audit.Record(ns.ComponentConsole, ns.EventCureStarted, req.PressID)
	writeJSON(w, http.StatusOK, cureResponse{
		Status:      "curing",
		Header:      s.temp.EffectiveHeader(),
		TimerStart:  s.timer.Started(),
		LiftState:   s.lift.State(),
		MoldVersion: current.MapVersion,
	})
}

func (s *Server) handleShutdown(w http.ResponseWriter, _ *http.Request) {
	_ = s.steam.Shutdown(s.cool)
	s.audit.Record(ns.ComponentConsole, ns.EventShutdown, "console")
	writeJSON(w, http.StatusOK, map[string]string{
		"steam": s.steam.State(),
		"cool":  s.cool.State(),
	})
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	component := ns.Component(r.URL.Query().Get("component"))
	event := ns.EventType(r.URL.Query().Get("event"))
	entries := s.audit.Filter(component, event)
	last, hasLast := s.audit.Last()
	payload := map[string]any{
		"total":   s.audit.Total(),
		"summary": s.audit.Summary(),
		"entries": entries,
		"last":    map[string]any{"present": hasLast, "entry": last},
	}
	writeJSON(w, http.StatusOK, payload)
}
