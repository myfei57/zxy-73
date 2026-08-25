package cool

func (c *CoolSystem) State() string {
	if c.stopped {
		return "stopped"
	}
	if c.running {
		return "running"
	}
	return "idle"
}
