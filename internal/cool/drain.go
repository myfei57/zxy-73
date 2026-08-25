package cool

type DrainState string

const (
	DrainClosed DrainState = "closed"
	DrainOpen   DrainState = "open"
	Draining    DrainState = "draining"
)

func (c *CoolSystem) Drain() DrainState {
	if c.draining {
		return Draining
	}
	if c.drainOpen {
		return DrainOpen
	}
	return DrainClosed
}

func (c *CoolSystem) OpenDrain() {
	c.drainOpen = true
	c.draining = true
}

func (c *CoolSystem) CloseDrain() {
	c.drainOpen = false
	c.draining = false
}
