package store

import "encoding/json"

func (s *Store) PutRecord(collection string, id string, value any) error {
	return s.PutJSON(collection+":"+id, value)
}

func (s *Store) GetRecord(collection string, id string, target any) error {
	return s.LoadJSON(collection+":"+id, target)
}

func (s *Store) ListRecords(collection string) []string {
	prefix := collection + ":"
	keys := s.ListKeys()
	out := make([]string, 0)
	for _, key := range keys {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			out = append(out, key[len(prefix):])
		}
	}
	return out
}

func (s *Store) MarshalRecords(collection string) ([]byte, error) {
	ids := s.ListRecords(collection)
	records := make(map[string]json.RawMessage, len(ids))
	for _, id := range ids {
		data, ok, err := s.Get(collection + ":" + id)
		if err != nil {
			return nil, err
		}
		if ok {
			records[id] = append(json.RawMessage(nil), data...)
		}
	}
	return json.Marshal(records)
}
