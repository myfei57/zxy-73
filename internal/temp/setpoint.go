package temp

func (c *Controller) SetPoint(value float64) float64 {
	c.setpoint = value
	return c.header.Set(c.pressID, value)
}

func (c *Controller) EffectiveHeader() float64 {
	return c.header.Effective()
}
