package mold

func (r *Registry) Swap(id string, sensorIDs []string) (Mold, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.molds[id]
	if !ok {
		return Mold{}, ErrMoldNotFound
	}
	return m, nil
}

func (r *Registry) CurrentMap(id string) (Mold, error) {
	return r.Get(id)
}
