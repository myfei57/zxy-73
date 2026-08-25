package temp

import "errors"

func (c *Controller) Ramp(target float64, steps int) error {
	if steps <= 0 {
		return errors.New("ramp steps must be positive")
	}
	c.mu.Lock()
	start := c.current
	c.mu.Unlock()
	delta := (target - start) / float64(steps)
	for i := 1; i <= steps; i++ {
		c.mu.Lock()
		c.current = start + delta*float64(i)
		c.mu.Unlock()
	}
	return nil
}

func (c *Controller) SetpointReached() bool {
	return c.Current() >= c.Setpoint()
}
