package audit

func (r *Recorder) Summary() map[string]int {
	counts := r.CountByEvent()
	out := make(map[string]int, len(counts))
	for event, count := range counts {
		out[event.String()] = count
	}
	return out
}
