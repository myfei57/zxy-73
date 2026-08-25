package temp

type persistedReading struct {
	PressID  string  `json:"press_id"`
	Setpoint float64 `json:"setpoint"`
	Current  float64 `json:"current"`
}

func (c *Controller) Persist() error {
	if c.store == nil {
		return nil
	}
	c.mu.Lock()
	reading := persistedReading{
		PressID:  c.pressID,
		Setpoint: c.setpoint,
		Current:  c.current,
	}
	c.mu.Unlock()
	return c.store.SaveJSON("temp:"+c.pressID, reading)
}

func (c *Controller) Read() (persistedReading, bool, error) {
	if c.store == nil {
		return persistedReading{}, false, nil
	}
	var reading persistedReading
	if err := c.store.LoadJSON("temp:"+c.pressID, &reading); err != nil {
		return persistedReading{}, false, err
	}
	return reading, true, nil
}
