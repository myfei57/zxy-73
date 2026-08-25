package mold

import "errors"

func ValidateSpec(spec string) error {
	if spec == "" {
		return errors.New("mold spec is empty")
	}
	return nil
}

func ValidateSensorIDs(ids []string) error {
	seen := make(map[string]bool)
	for _, id := range ids {
		if id == "" {
			return errors.New("sensor id is empty")
		}
		if seen[id] {
			return errors.New("duplicate sensor id")
		}
		seen[id] = true
	}
	return nil
}

func (r *Registry) ValidateSwap(id string, sensorIDs []string) error {
	if _, err := r.Get(id); err != nil {
		return err
	}
	return ValidateSensorIDs(sensorIDs)
}
