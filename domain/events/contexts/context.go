package contexts

type context struct {
	eventName string
	payload   map[string]any
}

func (c *context) EventName() string {
	return c.eventName
}

func (c *context) Payload() map[string]any {
	return c.payload
}

func (c *context) HasValue(key string) bool {
	_, ok := c.payload[key]
	return ok
}

func (c *context) Value(key string) any {
	return c.payload[key]
}
