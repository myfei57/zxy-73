package lift

type SteamPurgePort interface {
	Purge() error
}

func (l *LiftControl) Close(steam SteamPurgePort) error {
	if err := steam.Purge(); err != nil {
		return err
	}
	l.closed = true
	return nil
}
