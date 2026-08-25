package mold

import (
	"errors"
	"sort"
	"sync"
)

var ErrMoldNotFound = errors.New("mold not found")

type Mold struct {
	ID         string
	Spec       string
	MapVersion int
	SensorIDs  []string
}

type Registry struct {
	mu    sync.RWMutex
	molds map[string]Mold
}

func NewRegistry() *Registry {
	return &Registry{molds: make(map[string]Mold)}
}

func (r *Registry) Register(m Mold) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if m.MapVersion == 0 {
		m.MapVersion = 1
	}
	r.molds[m.ID] = m
}

func (r *Registry) Get(id string) (Mold, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.molds[id]
	if !ok {
		return Mold{}, ErrMoldNotFound
	}
	return m, nil
}

func (r *Registry) IDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.molds))
	for id := range r.molds {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
