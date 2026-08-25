package batch

func (s *Scheduler) NextMoldSeq() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.nextMold
}
