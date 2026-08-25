package temp

func (c *Controller) Durable() bool {
	_, ok, err := c.Read()
	if err != nil {
		return false
	}
	return ok
}
