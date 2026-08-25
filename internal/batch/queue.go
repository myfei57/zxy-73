package batch

func (s *Scheduler) EnqueueReport(pressID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reportQueue = append(s.reportQueue, pressID)
}

func (s *Scheduler) DequeueReport() (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.reportQueue) == 0 {
		return "", false
	}
	pressID := s.reportQueue[0]
	s.reportQueue = s.reportQueue[1:]
	return pressID, true
}

func (s *Scheduler) QueueLength() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.reportQueue)
}
