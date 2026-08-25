package console

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) registerStoreOps() {
	s.router.Post("/store/record", s.handleStorePutRecord)
	s.router.Get("/store/record/{collection}/{id}", s.handleStoreGetRecord)
	s.router.Post("/store/put", s.handleStorePutJSON)
	s.router.Get("/store/string/{key}", s.handleStoreGetString)
	s.router.Get("/store/int/{key}", s.handleStoreGetInt)
	s.router.Delete("/store/{key}", s.handleStoreDelete)
	s.router.Delete("/store/prefix/{prefix}", s.handleStoreDeletePrefix)
	s.router.Get("/store/collection/{collection}", s.handleStoreCollection)
}

func (s *Server) handleStorePutRecord(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Collection string `json:"collection"`
		ID         string `json:"id"`
		Value      any    `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.store.PutRecord(payload.Collection, payload.ID, payload.Value); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"key": payload.Collection + ":" + payload.ID})
}

func (s *Server) handleStoreGetRecord(w http.ResponseWriter, r *http.Request) {
	var value map[string]any
	if err := s.store.GetRecord(chi.URLParam(r, "collection"), chi.URLParam(r, "id"), &value); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func (s *Server) handleStorePutJSON(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.store.PutJSON(payload.Key, payload.Value); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"key": payload.Key})
}

func (s *Server) handleStoreGetString(w http.ResponseWriter, r *http.Request) {
	value, ok, err := s.store.GetString(chi.URLParam(r, "key"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"present": ok, "value": value})
}

func (s *Server) handleStoreGetInt(w http.ResponseWriter, r *http.Request) {
	value, ok, err := s.store.GetInt(chi.URLParam(r, "key"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"present": ok, "value": value})
}

func (s *Server) handleStoreDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Delete(chi.URLParam(r, "key")); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (s *Server) handleStoreDeletePrefix(w http.ResponseWriter, r *http.Request) {
	removed := s.store.DeletePrefix(chi.URLParam(r, "prefix"))
	writeJSON(w, http.StatusOK, map[string]int{"removed": removed})
}

func (s *Server) handleStoreCollection(w http.ResponseWriter, r *http.Request) {
	collection := chi.URLParam(r, "collection")
	records, err := s.store.MarshalRecords(collection)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ids":     s.store.ListRecords(collection),
		"records": json.RawMessage(records),
	})
}
