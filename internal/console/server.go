package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"rubbercure/internal/audit"
	"rubbercure/internal/batch"
	"rubbercure/internal/cool"
	"rubbercure/internal/lift"
	"rubbercure/internal/mold"
	"rubbercure/internal/press"
	"rubbercure/internal/steam"
	"rubbercure/internal/store"
	"rubbercure/internal/temp"
	"rubbercure/internal/timer"
)

type Server struct {
	router   chi.Router
	fleet    *press.Fleet
	registry *mold.Registry
	steam    *steam.SteamSystem
	cool     *cool.CoolSystem
	temp     *temp.Controller
	timer    *timer.CureTimer
	lift     *lift.LiftControl
	batch    *batch.Scheduler
	audit    *audit.Recorder
	store    *store.Store
}

func NewServer(
	fleet *press.Fleet,
	registry *mold.Registry,
	steamSystem *steam.SteamSystem,
	coolSystem *cool.CoolSystem,
	tempController *temp.Controller,
	cureTimer *timer.CureTimer,
	liftControl *lift.LiftControl,
	scheduler *batch.Scheduler,
	recorder *audit.Recorder,
	st *store.Store,
) *Server {
	s := &Server{
		fleet:    fleet,
		registry: registry,
		steam:    steamSystem,
		cool:     coolSystem,
		temp:     tempController,
		timer:    cureTimer,
		lift:     liftControl,
		batch:    scheduler,
		audit:    recorder,
		store:    st,
	}
	s.router = chi.NewRouter()
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}
