package temp

func (c *Controller) WithinTolerance(tolerance float64) bool {
	if tolerance < 0 {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	diff := c.current - c.setpoint
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
