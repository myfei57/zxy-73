package lift

type SteamVentPort interface {
	Vent() error
}

func (l *LiftControl) Open(steam SteamVentPort) error {
	l.closed = false
	return steam.Vent()
}
