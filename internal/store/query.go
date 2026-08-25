package store

import (
	"encoding/json"
	"strconv"
)

func (s *Store) GetString(key string) (string, bool, error) {
	data, ok, err := s.Get(key)
	if err != nil || !ok {
		return "", ok, err
	}
	return string(data), true, nil
}

func (s *Store) GetInt(key string) (int, bool, error) {
	data, ok, err := s.Get(key)
	if err != nil || !ok {
		return 0, ok, err
	}
	value, err := strconv.Atoi(string(data))
	return value, err == nil, err
}

func (s *Store) PutJSON(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.Put(key, data)
}
