package store

import "sort"

func (s *Store) ListKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.records))
	for key := range s.records {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.records)
}

func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = make(map[string][]byte)
	return s.persistLocked()
}

func (s *Store) DeletePrefix(prefix string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	removed := 0
	for key := range s.records {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(s.records, key)
			removed++
		}
	}
	if removed > 0 {
		_ = s.persistLocked()
	}
	return removed
}
