package components

import (
	"slices"

	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
)

type component struct {
	renderables.StylableRenderable

	keyname  string
	events   []events.Event
	children []Component
}

func (c *component) Keyname() string {
	return c.keyname
}

func (c *component) HasEvents() bool {
	return len(c.events) > 0
}

func (c *component) Events() []events.Event {
	return slices.Clone(c.events)
}

func (c *component) HasChildren() bool {
	return len(c.children) > 0
}

func (c *component) Children() []Component {
	return slices.Clone(c.children)
}
