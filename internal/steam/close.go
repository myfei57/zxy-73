package steam

type Stopper interface {
	Stop()
}

func (s *SteamSystem) Close() error {
	s.valveOpen = false
	s.pressure.Release()
	return nil
}

func (s *SteamSystem) Closed() bool {
	return !s.valveOpen
}

func (s *SteamSystem) Shutdown(stopper Stopper) error {
	stopper.Stop()
	return s.Close()
}
