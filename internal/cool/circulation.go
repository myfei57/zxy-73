package cool

type PumpState struct {
	Running bool    `json:"running"`
	Flow    float64 `json:"flow"`
	Speed   int     `json:"speed"`
}

func (c *CoolSystem) SetSpeed(speed int) {
	if speed < 0 {
		speed = 0
	}
	if speed > 100 {
		speed = 100
	}
	c.speed = speed
}

func (c *CoolSystem) Speed() int {
	return c.speed
}

func (c *CoolSystem) Flow() float64 {
	if !c.Running() {
		return 0
	}
	return float64(c.speed) * 0.35
}

func (c *CoolSystem) PumpState() PumpState {
	return PumpState{Running: c.Running(), Flow: c.Flow(), Speed: c.speed}
}
