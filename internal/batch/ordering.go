package batch

import "errors"

func (s *Scheduler) NextAvailablePress() (string, error) {
	for _, p := range s.fleet.List() {
		if !p.Busy() {
			return p.ID, nil
		}
	}
	return "", errors.New("no idle press available")
}

func (s *Scheduler) AvailablePressCount() int {
	count := 0
	for _, p := range s.fleet.List() {
		if !p.Busy() {
			count++
		}
	}
	return count
}

func (s *Scheduler) BusyPressIDs() []string {
	out := make([]string, 0)
	for _, p := range s.fleet.List() {
		if p.Busy() {
			out = append(out, p.ID)
		}
	}
	return out
}
