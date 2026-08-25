package cool

type CoolSystem struct {
	running   bool
	stopped   bool
	speed     int
	drainOpen bool
	draining  bool
	alarm     bool
}

func NewCoolSystem() *CoolSystem {
	return &CoolSystem{running: true, stopped: false}
}

func (c *CoolSystem) Running() bool {
	return c.running && !c.stopped
}

func (c *CoolSystem) Stopped() bool {
	return c.stopped
}
