package steam

func (s *SteamSystem) State() string {
	switch {
	case s.ventOpen:
		return "venting"
	case s.valveOpen:
		return "heating"
	case s.purgeOpen:
		return "purging"
	default:
		return "idle"
	}
}
