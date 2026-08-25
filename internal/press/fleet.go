package press

import (
	"errors"
	"sort"
	"sync"
)

var ErrPressNotFound = errors.New("press not found")

type Fleet struct {
	mu      sync.RWMutex
	presses map[string]*Press
}

func NewFleet() *Fleet {
	return &Fleet{presses: make(map[string]*Press)}
}

func (f *Fleet) Add(p *Press) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.presses[p.ID] = p
}

func (f *Fleet) Get(id string) (*Press, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	p, ok := f.presses[id]
	if !ok {
		return nil, ErrPressNotFound
	}
	return p, nil
}

func (f *Fleet) List() []*Press {
	f.mu.RLock()
	defer f.mu.RUnlock()
	ids := make([]string, 0, len(f.presses))
	for id := range f.presses {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]*Press, 0, len(ids))
	for _, id := range ids {
		out = append(out, f.presses[id])
	}
	return out
}

func (f *Fleet) AssignMold(id string, moldID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.presses[id]
	if !ok {
		return ErrPressNotFound
	}
	if p.Busy() {
		return errors.New("press is busy")
	}
	p.MoldID = moldID
	p.CycleSeq++
	return nil
}

func (f *Fleet) SetState(id string, state State) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.presses[id]
	if !ok {
		return ErrPressNotFound
	}
	p.State = state
	return nil
}

func (f *Fleet) ReleaseMold(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.presses[id]
	if !ok {
		return ErrPressNotFound
	}
	p.MoldID = ""
	p.State = StateIdle
	return nil
}
