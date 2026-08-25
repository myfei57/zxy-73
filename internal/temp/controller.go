package temp

import (
	"sync"

	"rubbercure/internal/steam"
	"rubbercure/internal/store"
)

type Controller struct {
	mu        sync.Mutex
	pressID   string
	header    *steam.Header
	store     *store.Store
	setpoint  float64
	current   float64
	integral  float64
	lowLimit  float64
	highLimit float64
}

func NewController(pressID string, header *steam.Header, st *store.Store) *Controller {
	if header == nil {
		header = steam.NewHeader()
	}
	return &Controller{pressID: pressID, header: header, store: st}
}

func (c *Controller) PressID() string {
	return c.pressID
}

func (c *Controller) Setpoint() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.setpoint
}

func (c *Controller) Current() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current
}
