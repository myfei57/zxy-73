package temp

import "errors"

func (c *Controller) SetLimits(low float64, high float64) error {
	if high <= low {
		return errors.New("high limit must be greater than low limit")
	}
	c.mu.Lock()
	c.lowLimit = low
	c.highLimit = high
	c.mu.Unlock()
	return nil
}

func (c *Controller) LowLimit() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lowLimit
}

func (c *Controller) HighLimit() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.highLimit
}

func (c *Controller) WithinLimits() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current >= c.lowLimit && c.current <= c.highLimit
}
