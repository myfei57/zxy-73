package press

type State string

const (
	StateIdle    State = "idle"
	StateClosing State = "closing"
	StateCuring  State = "curing"
	StateOpening State = "opening"
	StateFault   State = "fault"
)

type Press struct {
	ID       string
	Name     string
	State    State
	MoldID   string
	CycleSeq uint64
}

func New(id string, name string) *Press {
	return &Press{ID: id, Name: name, State: StateIdle}
}

func (p *Press) Busy() bool {
	return p.State == StateClosing || p.State == StateCuring || p.State == StateOpening
}
