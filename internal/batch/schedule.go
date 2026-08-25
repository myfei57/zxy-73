package batch

func (s *Scheduler) ResetTimerLatch() {
	s.timer.ResetLatch()
}

func (s *Scheduler) LatchReleased() bool {
	return !s.timer.Latched()
}
