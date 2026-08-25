package steam

func (s *SteamSystem) OpenValve() error {
	s.valveOpen = true
	s.pressure.Set(s.header.Effective())
	return nil
}

func (s *SteamSystem) ValveOpen() bool {
	return s.valveOpen
}
