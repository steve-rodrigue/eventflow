package components

import (
	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/templates"
)

// Builder represents a component builder.
type Builder interface {
	Create() Builder

	WithKeyname(keyname string) Builder
	WithStyle(style templates.Template) Builder
	WithTemplate(template templates.Template) Builder

	AddEvent(event events.Event) Builder
	WithEvents(events []events.Event) Builder

	Now() (Component, error)
}

// Component represents a component.
type Component interface {
	renderables.StylableRenderable

	Keyname() string

	HasEvents() bool
	Events() []events.Event
}
