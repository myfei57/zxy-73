package cool

func (c *CoolSystem) SetAlarm() {
	c.alarm = true
}

func (c *CoolSystem) ClearAlarm() {
	c.alarm = false
}

func (c *CoolSystem) Alarm() bool {
	return c.alarm
}
