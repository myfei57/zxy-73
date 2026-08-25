package console

import (
	"net/http"
)

func (s *Server) routes() {
	s.router.Get("/healthz", s.handleHealth)
	s.router.Get("/status", s.handleStatus)
	s.router.Post("/cure", s.handleCure)
	s.router.Post("/shutdown", s.handleShutdown)
	s.router.Get("/audit", s.handleAudit)
	s.registerOps()
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
