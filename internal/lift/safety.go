package lift

import "errors"

func (l *LiftControl) ValidateRequiredSensors() error {
	if l.requiredSensors < 1 {
		return errors.New("required sensor count must be positive")
	}
	return nil
}
