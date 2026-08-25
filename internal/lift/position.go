package lift

type Position int

const (
	PositionOpen Position = iota
	PositionClosed
)

func (l *LiftControl) Position() Position {
	if l.closed {
		return PositionClosed
	}
	return PositionOpen
}

func (l *LiftControl) IsFullyOpen() bool {
	return !l.closed
}
