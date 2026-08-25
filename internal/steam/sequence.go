package steam

type Step struct {
	Name string
	Run  func() error
}

func (s *SteamSystem) WarmupSequence() []Step {
	return []Step{
		{Name: "trap", Run: s.Trap},
		{Name: "valve", Run: s.OpenValve},
	}
}

func (s *SteamSystem) ShutdownSequence(stopper Stopper) []Step {
	return []Step{
		{Name: "cool-stop", Run: func() error { stopper.Stop(); return nil }},
		{Name: "steam-close", Run: s.Close},
	}
}
