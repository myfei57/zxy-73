package temp

func (c *Controller) SetPoint(value float64) float64 {
	c.mu.Lock()
	c.setpoint = value
	c.mu.Unlock()
	return c.header.Set(c.pressID, value)
}

func (c *Controller) EffectiveHeader() float64 {
	return c.header.Effective()
}
