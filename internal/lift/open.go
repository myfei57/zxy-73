package lift

type SteamVentPort interface {
	Vent() error
}

func (l *LiftControl) Open(steam SteamVentPort) error {
	if err := steam.Vent(); err != nil {
		return err
	}
	l.closed = false
	return nil
}
