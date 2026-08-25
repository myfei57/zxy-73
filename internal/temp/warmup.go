package temp

type SteamWarmupPort interface {
	Trap() error
	OpenValve() error
}

func (c *Controller) Warmup(port SteamWarmupPort) error {
	if err := port.Trap(); err != nil {
		return err
	}
	if err := port.OpenValve(); err != nil {
		return err
	}
	c.mu.Lock()
	c.current = c.setpoint
	c.mu.Unlock()
	return nil
}
