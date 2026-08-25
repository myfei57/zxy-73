package steam

func (s *SteamSystem) Vent() error {
	s.ventOpen = true
	s.pressure.Release()
	return nil
}

func (s *SteamSystem) Vented() bool {
	return s.ventOpen
}
