package steam

func (s *SteamSystem) Purge() error {
	s.purgeOpen = true
	return nil
}

func (s *SteamSystem) Purged() bool {
	return s.purgeOpen
}
