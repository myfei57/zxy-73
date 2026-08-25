package batch

import "time"

type CycleReport struct {
	PressID   string    `json:"press_id"`
	MoldID    string    `json:"mold_id"`
	Completed time.Time `json:"completed"`
	Sequence  int       `json:"sequence"`
}

func (s *Scheduler) BuildReport(pressID string) (CycleReport, error) {
	p, err := s.fleet.Get(pressID)
	if err != nil {
		return CycleReport{}, err
	}
	return CycleReport{
		PressID:   p.ID,
		MoldID:    p.MoldID,
		Completed: time.Now().UTC(),
		Sequence:  s.NextMoldSeq(),
	}, nil
}
