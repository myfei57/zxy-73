package cool

func (c *CoolSystem) Stop() {
	c.running = false
	c.stopped = true
}
