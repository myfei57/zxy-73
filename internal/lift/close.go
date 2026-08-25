package lift

type SteamPurgePort interface {
	Purge() error
}

func (l *LiftControl) Close(steam SteamPurgePort) error {
	l.closed = true
	return steam.Purge()
}
