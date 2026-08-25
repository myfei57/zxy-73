package press

import "errors"

func (p *Press) StartCycle() error {
	if p.State != StateIdle {
		return errors.New("press must be idle to start a cycle")
	}
	p.State = StateCuring
	p.CycleSeq++
	return nil
}

func (p *Press) FinishCycle() error {
	if p.State != StateCuring {
		return errors.New("press is not curing")
	}
	p.State = StateIdle
	return nil
}

func (p *Press) MarkFault() {
	p.State = StateFault
}

func (p *Press) Recover() {
	p.State = StateIdle
}
