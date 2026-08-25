package mold

type Summary struct {
	Total    int            `json:"total"`
	BySpec   map[string]int `json:"by_spec"`
	Versions map[string]int `json:"versions"`
}

func (r *Registry) Summary() Summary {
	ids := r.IDs()
	summary := Summary{Total: len(ids), BySpec: make(map[string]int), Versions: make(map[string]int)}
	for _, id := range ids {
		m, err := r.Get(id)
		if err != nil {
			continue
		}
		summary.BySpec[m.Spec]++
		summary.Versions[id] = m.MapVersion
	}
	return summary
}

func (r *Registry) List() []Mold {
	ids := r.IDs()
	out := make([]Mold, 0, len(ids))
	for _, id := range ids {
		if m, err := r.Get(id); err == nil {
			out = append(out, m)
		}
	}
	return out
}
