package console

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rubbercure/internal/mold"
	"rubbercure/internal/ns"
	"rubbercure/internal/press"
	"rubbercure/internal/steam"
	"rubbercure/internal/timer"
)

func (s *Server) registerOps() {
	s.router.Get("/presses", s.handlePresses)
	s.router.Post("/presses/{id}/start", s.handlePressStart)
	s.router.Post("/presses/{id}/finish", s.handlePressFinish)
	s.router.Post("/presses/{id}/fault", s.handlePressFault)
	s.router.Post("/presses/{id}/recover", s.handlePressRecover)
	s.router.Post("/presses/{id}/state", s.handlePressState)

	s.router.Get("/molds", s.handleMolds)
	s.router.Post("/molds", s.handleMoldRegister)
	s.router.Post("/molds/{id}/swap", s.handleMoldSwap)

	s.router.Get("/steam", s.handleSteamStatus)
	s.router.Get("/steam/sequence", s.handleSteamSequence)
	s.router.Post("/steam/header/{pressID}/set", s.handleHeaderSet)
	s.router.Post("/steam/header/{pressID}/reset", s.handleHeaderReset)

	s.router.Get("/temp", s.handleTempStatus)
	s.router.Post("/temp/ramp", s.handleTempRamp)
	s.router.Post("/temp/regulate", s.handleTempRegulate)
	s.router.Post("/temp/limits", s.handleTempLimits)
	s.router.Post("/temp/persist", s.handleTempPersist)

	s.router.Get("/timer", s.handleTimerStatus)
	s.router.Post("/timer/curve", s.handleTimerCurve)

	s.router.Get("/lift", s.handleLiftStatus)
	s.router.Post("/lift/close", s.handleLiftClose)
	s.router.Post("/lift/open", s.handleLiftOpen)
	s.router.Get("/lift/check/{moldID}", s.handleLiftCheck)

	s.router.Get("/cool", s.handleCoolStatus)
	s.router.Post("/cool/speed", s.handleCoolSpeed)
	s.router.Post("/cool/drain/open", s.handleCoolDrainOpen)
	s.router.Post("/cool/drain/close", s.handleCoolDrainClose)
	s.router.Post("/cool/alarm/set", s.handleCoolAlarmSet)
	s.router.Post("/cool/alarm/clear", s.handleCoolAlarmClear)

	s.router.Get("/store", s.handleStoreStatus)
	s.router.Post("/store/clear", s.handleStoreClear)
	s.router.Get("/batch", s.handleBatchStatus)
	s.router.Post("/batch/profile", s.handleBatchProfile)
	s.router.Post("/batch/latch/reset", s.handleBatchLatchReset)
	s.router.Post("/batch/report/done", s.handleBatchReportDone)
	s.router.Get("/batch/report/{pressID}", s.handleBatchReport)
	s.router.Post("/batch/report/enqueue", s.handleReportEnqueue)
	s.router.Get("/batch/report/dequeue", s.handleReportDequeue)
	s.registerStoreOps()
}

func (s *Server) handlePresses(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"presses": s.fleet.List(),
		"summary": s.fleet.Summary(),
		"busy":    s.fleet.BusyPresses(),
	})
}

func (s *Server) handlePressStart(w http.ResponseWriter, r *http.Request) {
	p, err := s.fleet.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if err := p.StartCycle(); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	s.audit.Record(ns.ComponentPress, ns.EventCureStarted, p.ID)
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handlePressFinish(w http.ResponseWriter, r *http.Request) {
	p, err := s.fleet.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if err := p.FinishCycle(); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	s.audit.Record(ns.ComponentPress, ns.EventCureDone, p.ID)
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handlePressFault(w http.ResponseWriter, r *http.Request) {
	p, err := s.fleet.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	p.MarkFault()
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handlePressRecover(w http.ResponseWriter, r *http.Request) {
	p, err := s.fleet.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	p.Recover()
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handlePressState(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		State string `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	state := press.State(payload.State)
	if err := s.fleet.SetState(chi.URLParam(r, "id"), state); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	p, _ := s.fleet.Get(chi.URLParam(r, "id"))
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) handleMolds(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"molds":   s.registry.List(),
		"summary": s.registry.Summary(),
	})
}

func (s *Server) handleMoldRegister(w http.ResponseWriter, r *http.Request) {
	var m mold.Mold
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := mold.ValidateSpec(m.Spec); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := mold.ValidateSensorIDs(m.SensorIDs); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.registry.Register(m)
	s.audit.Record(ns.ComponentMold, ns.EventMoldRegistered, m.ID)
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleMoldSwap(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		SensorIDs []string `json:"sensor_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	id := chi.URLParam(r, "id")
	if err := s.registry.ValidateSwap(id, payload.SensorIDs); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	updated, err := s.registry.Swap(id, payload.SensorIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleSteamStatus(w http.ResponseWriter, _ *http.Request) {
	header := s.steam.Header()
	writeJSON(w, http.StatusOK, map[string]any{
		"state":      s.steam.State(),
		"purged":     s.steam.Purged(),
		"drained":    s.steam.Drained(),
		"valve":      s.steam.ValveOpen(),
		"vented":     s.steam.Vented(),
		"closed":     s.steam.Closed(),
		"pressure":   s.steam.Pressure().Current(),
		"over_limit": s.steam.Pressure().OverLimit(8.0),
		"effective":  header.Effective(),
		"max":        header.Max(),
		"min":        header.Min(),
		"requests":   header.Requests(),
	})
}

func (s *Server) handleSteamSequence(w http.ResponseWriter, _ *http.Request) {
	warmup := s.steam.WarmupSequence()
	shutdown := s.steam.ShutdownSequence(s.cool)
	names := func(steps []steam.Step) []string {
		out := make([]string, 0, len(steps))
		for _, step := range steps {
			out = append(out, step.Name)
		}
		return out
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"warmup":   names(warmup),
		"shutdown": names(shutdown),
	})
}

func (s *Server) handleHeaderSet(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Value float64 `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	header := s.steam.Header()
	effective := header.Set(chi.URLParam(r, "pressID"), payload.Value)
	writeJSON(w, http.StatusOK, map[string]float64{"effective": effective})
}

func (s *Server) handleHeaderReset(w http.ResponseWriter, r *http.Request) {
	header := s.steam.Header()
	writeJSON(w, http.StatusOK, map[string]float64{"effective": header.Reset(chi.URLParam(r, "pressID"))})
}

func (s *Server) handleTempStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"press":     s.temp.PressID(),
		"setpoint":  s.temp.Setpoint(),
		"current":   s.temp.Current(),
		"durable":   s.temp.Durable(),
		"header":    s.temp.EffectiveHeader(),
		"reached":   s.temp.SetpointReached(),
		"integral":  s.temp.Integral(),
		"low":       s.temp.LowLimit(),
		"high":      s.temp.HighLimit(),
		"in_range":  s.temp.WithinLimits(),
		"on_target": s.temp.WithinTolerance(2.0),
	})
}

func (s *Server) handleTempRamp(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Target float64 `json:"target"`
		Steps  int     `json:"steps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.temp.Ramp(payload.Target, payload.Steps); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"current": s.temp.Current()})
}

func (s *Server) handleTempPersist(w http.ResponseWriter, _ *http.Request) {
	if err := s.temp.Persist(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	reading, ok, err := s.temp.Read()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": ok, "reading": reading})
}

func (s *Server) handleTempRegulate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Target float64 `json:"target"`
		Gain   float64 `json:"gain"`
		Reset  float64 `json:"reset"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	state, err := s.temp.Regulate(payload.Target, payload.Gain, payload.Reset)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleTempLimits(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Low  float64 `json:"low"`
		High float64 `json:"high"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.temp.SetLimits(payload.Low, payload.High); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"low":  s.temp.LowLimit(),
		"high": s.temp.HighLimit(),
	})
}

func (s *Server) handleTimerStatus(w http.ResponseWriter, _ *http.Request) {
	curve := s.timer.CurrentCurve()
	lastCycle, hasCycle := s.timer.LastCycle()
	writeJSON(w, http.StatusOK, map[string]any{
		"running":     s.timer.Running(),
		"done":        s.timer.Done(),
		"latched":     s.timer.Latched(),
		"started":     s.timer.Started(),
		"cycles":      s.timer.CycleCount(),
		"last_cycle":  map[string]any{"present": hasCycle, "record": lastCycle},
		"curve":       curve,
		"duration":    s.timer.DurationSeconds(),
		"temperature": s.timer.Temperature(),
		"remaining":   s.timer.Remaining(time.Now().UTC()).String(),
	})
}

func (s *Server) handleTimerCurve(w http.ResponseWriter, r *http.Request) {
	var curve timer.Curve
	if err := json.NewDecoder(r.Body).Decode(&curve); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := timer.ValidateCurve(curve); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.batch.ProfileUpdate(curve)
	writeJSON(w, http.StatusOK, curve)
}

func (s *Server) handleLiftStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"state":    s.lift.State(),
		"position": int(s.lift.Position()),
		"closed":   s.lift.Closed(),
		"required": s.lift.RequiredSensors(),
		"valid":    s.lift.ValidateRequiredSensors() == nil,
		"open":     s.lift.IsFullyOpen(),
	})
}

func (s *Server) handleLiftClose(w http.ResponseWriter, _ *http.Request) {
	if err := s.lift.Close(s.steam); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	s.audit.Record(ns.ComponentLift, ns.EventPressClosed, "console")
	writeJSON(w, http.StatusOK, map[string]string{"state": s.lift.State()})
}

func (s *Server) handleLiftOpen(w http.ResponseWriter, _ *http.Request) {
	if err := s.lift.Open(s.steam); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	s.audit.Record(ns.ComponentLift, ns.EventMoldOpened, "console")
	writeJSON(w, http.StatusOK, map[string]string{"state": s.lift.State()})
}

func (s *Server) handleLiftCheck(w http.ResponseWriter, r *http.Request) {
	closed, err := s.lift.CheckClosed(s.registry, chi.URLParam(r, "moldID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"closed": closed})
}

func (s *Server) handleCoolStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"state":   s.cool.State(),
		"running": s.cool.Running(),
		"stopped": s.cool.Stopped(),
		"speed":   s.cool.Speed(),
		"flow":    s.cool.Flow(),
		"pump":    s.cool.PumpState(),
		"drain":   s.cool.Drain(),
		"alarm":   s.cool.Alarm(),
	})
}

func (s *Server) handleCoolSpeed(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Speed int `json:"speed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.cool.SetSpeed(payload.Speed)
	writeJSON(w, http.StatusOK, s.cool.PumpState())
}

func (s *Server) handleCoolDrainOpen(w http.ResponseWriter, _ *http.Request) {
	s.cool.OpenDrain()
	writeJSON(w, http.StatusOK, map[string]string{"drain": string(s.cool.Drain())})
}

func (s *Server) handleCoolDrainClose(w http.ResponseWriter, _ *http.Request) {
	s.cool.CloseDrain()
	writeJSON(w, http.StatusOK, map[string]string{"drain": string(s.cool.Drain())})
}

func (s *Server) handleCoolAlarmSet(w http.ResponseWriter, _ *http.Request) {
	s.cool.SetAlarm()
	writeJSON(w, http.StatusOK, map[string]bool{"alarm": s.cool.Alarm()})
}

func (s *Server) handleCoolAlarmClear(w http.ResponseWriter, _ *http.Request) {
	s.cool.ClearAlarm()
	writeJSON(w, http.StatusOK, map[string]bool{"alarm": s.cool.Alarm()})
}

func (s *Server) handleStoreStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"keys":     s.store.ListKeys(),
		"count":    s.store.Count(),
		"snapshot": len(s.store.Snapshot()),
	})
}

func (s *Server) handleStoreClear(w http.ResponseWriter, _ *http.Request) {
	if err := s.store.Clear(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": s.store.Count()})
}

func (s *Server) handleBatchStatus(w http.ResponseWriter, _ *http.Request) {
	next, err := s.batch.NextAvailablePress()
	writeJSON(w, http.StatusOK, map[string]any{
		"assignments":    s.batch.Assignments(),
		"next_press":     next,
		"next_error":     errText(err),
		"available":      s.batch.AvailablePressCount(),
		"busy":           s.batch.BusyPressIDs(),
		"next_mold":      s.batch.NextMoldSeq(),
		"pending":        s.batch.PendingMoldCount(),
		"latch_released": s.batch.LatchReleased(),
		"fleet":          s.batch.Fleet().Summary(),
		"registry":       s.batch.Registry().Summary(),
		"timer_curve":    s.batch.Timer().CurrentCurve(),
	})
}

func (s *Server) handleBatchProfile(w http.ResponseWriter, r *http.Request) {
	var curve timer.Curve
	if err := json.NewDecoder(r.Body).Decode(&curve); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := timer.ValidateCurve(curve); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.batch.ProfileUpdate(curve)
	writeJSON(w, http.StatusOK, s.timer.CurrentCurve())
}

func (s *Server) handleBatchLatchReset(w http.ResponseWriter, _ *http.Request) {
	s.batch.ResetTimerLatch()
	writeJSON(w, http.StatusOK, map[string]bool{"released": s.batch.LatchReleased()})
}

func (s *Server) handleBatchReportDone(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PressID string `json:"press_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.batch.ReportDone(payload.PressID); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusOK, s.batch.Assignments())
}

func (s *Server) handleBatchReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.batch.BuildReport(chi.URLParam(r, "pressID"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleReportEnqueue(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PressID string `json:"press_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	s.batch.EnqueueReport(payload.PressID)
	writeJSON(w, http.StatusOK, map[string]int{"queued": s.batch.QueueLength()})
}

func (s *Server) handleReportDequeue(w http.ResponseWriter, _ *http.Request) {
	pressID, ok := s.batch.DequeueReport()
	writeJSON(w, http.StatusOK, map[string]any{"press_id": pressID, "present": ok})
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
