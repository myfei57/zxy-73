package steam

func (s *SteamSystem) Trap() error {
	s.trapOpen = true
	return nil
}

func (s *SteamSystem) Drained() bool {
	return s.trapOpen
}
