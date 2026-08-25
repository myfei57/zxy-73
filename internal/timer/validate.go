package timer

import "errors"

func ValidateCurve(curve Curve) error {
	if curve.Name == "" {
		return errors.New("curve name is empty")
	}
	if curve.Temperature <= 0 {
		return errors.New("curve temperature must be positive")
	}
	if curve.DurationSeconds <= 0 {
		return errors.New("curve duration must be positive")
	}
	return nil
}
