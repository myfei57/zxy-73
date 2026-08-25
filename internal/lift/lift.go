package lift

type LiftControl struct {
	closed          bool
	requiredSensors int
}

func NewLift(requiredSensors int) *LiftControl {
	return &LiftControl{requiredSensors: requiredSensors}
}

func (l *LiftControl) Closed() bool {
	return l.closed
}

func (l *LiftControl) RequiredSensors() int {
	return l.requiredSensors
}
