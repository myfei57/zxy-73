package lift

func (l *LiftControl) State() string {
	if l.closed {
		return "closed"
	}
	return "open"
}
