package press

type Summary struct {
	Total   int           `json:"total"`
	Busy    int           `json:"busy"`
	Idle    int           `json:"idle"`
	Fault   int           `json:"fault"`
	ByState map[State]int `json:"by_state"`
}

func (f *Fleet) Summary() Summary {
	presses := f.List()
	summary := Summary{Total: len(presses), ByState: make(map[State]int)}
	for _, p := range presses {
		summary.ByState[p.State]++
		switch p.State {
		case StateIdle:
			summary.Idle++
		case StateFault:
			summary.Fault++
		default:
			summary.Busy++
		}
	}
	return summary
}

func (f *Fleet) BusyPresses() []*Press {
	out := make([]*Press, 0)
	for _, p := range f.List() {
		if p.Busy() {
			out = append(out, p)
		}
	}
	return out
}
