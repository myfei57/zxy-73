package audit

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"rubbercure/internal/ns"
	"rubbercure/internal/store"
)

type Entry struct {
	ID        string       `json:"id"`
	At        time.Time    `json:"at"`
	Component ns.Component `json:"component"`
	Event     ns.EventType `json:"event"`
	Message   string       `json:"message"`
}

type Recorder struct {
	mu      sync.Mutex
	store   *store.Store
	entries []Entry
}

func NewRecorder(st *store.Store) *Recorder {
	return &Recorder{store: st}
}

func (r *Recorder) Record(component ns.Component, event ns.EventType, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry := Entry{
		ID:        uuid.NewString(),
		At:        time.Now().UTC(),
		Component: component,
		Event:     event,
		Message:   message,
	}
	r.entries = append(r.entries, entry)
	if len(r.entries) > 512 {
		r.entries = append([]Entry(nil), r.entries[len(r.entries)-512:]...)
	}
	if r.store != nil {
		_ = r.store.SaveJSON("audit:latest", r.entries)
	}
}

func (r *Recorder) Entries() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Entry, len(r.entries))
	copy(out, r.entries)
	return out
}

func (r *Recorder) Last() (Entry, bool) {
	entries := r.Entries()
	if len(entries) == 0 {
		return Entry{}, false
	}
	return entries[len(entries)-1], true
}
