package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var ErrNotFound = errors.New("store key not found")

type Store struct {
	mu      sync.RWMutex
	dir     string
	name    string
	records map[string][]byte
}

func New(dir string, name string) (*Store, error) {
	if dir == "" {
		dir = os.TempDir()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{
		dir:     dir,
		name:    name,
		records: make(map[string][]byte),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

func (s *Store) path() string {
	return filepath.Join(s.dir, s.name+".json")
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.records)
}

func (s *Store) Put(key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	clone := append([]byte(nil), value...)
	s.records[key] = clone
	return s.persistLocked()
}

func (s *Store) Get(key string) ([]byte, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.records[key]
	if !ok {
		return nil, false, nil
	}
	return append([]byte(nil), value...), true, nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[key]; !ok {
		return ErrNotFound
	}
	delete(s.records, key)
	return s.persistLocked()
}

func (s *Store) Snapshot() map[string][]byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string][]byte, len(s.records))
	for k, v := range s.records {
		out[k] = append([]byte(nil), v...)
	}
	return out
}

func (s *Store) SaveJSON(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.Put(key, data)
}

func (s *Store) LoadJSON(key string, target any) error {
	data, ok, err := s.Get(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return json.Unmarshal(data, target)
}

func (s *Store) persistLocked() error {
	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(s.path(), data)
}
