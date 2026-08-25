package batch

func (s *Scheduler) PendingMoldCount() int {
	return s.NextMoldSeq()
}
