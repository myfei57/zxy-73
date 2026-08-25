package temp

import "errors"

type RegulationState struct {
	Setpoint   float64
	Current    float64
	Error      float64
	Integral   float64
	Derivative float64
	Output     float64
}

func (c *Controller) Regulate(target float64, gain float64, reset float64) (RegulationState, error) {
	if gain <= 0 {
		return RegulationState{}, errors.New("gain must be positive")
	}
	if reset < 0 {
		return RegulationState{}, errors.New("reset must be non-negative")
	}
	c.mu.Lock()
	current := c.current
	c.integral += (target - current) * reset
	integral := c.integral
	c.mu.Unlock()
	errValue := target - current
	output := gain*errValue + integral
	state := RegulationState{
		Setpoint:   target,
		Current:    current,
		Error:      errValue,
		Integral:   integral,
		Derivative: 0,
		Output:     output,
	}
	return state, nil
}

func (c *Controller) Integral() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.integral
}
