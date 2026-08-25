package lift

import "rubbercure/internal/mold"

func (l *LiftControl) CheckClosed(registry *mold.Registry, moldID string) (bool, error) {
	current, err := registry.CurrentMap(moldID)
	if err != nil {
		return false, err
	}
	return len(current.SensorIDs) >= l.requiredSensors, nil
}
