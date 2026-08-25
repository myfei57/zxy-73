package audit

import "rubbercure/internal/ns"

func (r *Recorder) Filter(component ns.Component, event ns.EventType) []Entry {
	entries := r.Entries()
	out := make([]Entry, 0)
	for _, entry := range entries {
		if component != "" && entry.Component != component {
			continue
		}
		if event != "" && entry.Event != event {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func (r *Recorder) CountByEvent() map[ns.EventType]int {
	entries := r.Entries()
	out := make(map[ns.EventType]int)
	for _, entry := range entries {
		out[entry.Event]++
	}
	return out
}

func (r *Recorder) Total() int {
	return len(r.Entries())
}
