package components

import (
	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

// NewRenderer creates a new component renderer
func NewRenderer(templateRenderer templates.Renderer) Renderer {
	return &renderer{
		templateRenderer: templateRenderer,
	}
}

// NewBuilder creates a new component builder
func NewBuilder(stylableBuilder renderables.StylableBuilder) Builder {
	return &builder{
		stylableBuilder: stylableBuilder,
	}
}

// Renderer represents a renderable renderer.
type Renderer interface {
	Render(component Component, params renderables.Params) string
	RenderList(components []Component, params []renderables.Params) string

	RenderStyle(component Component, params renderables.Params) string
	RenderStyleList(components []Component, params []renderables.Params) string
}

// Builder represents a component builder.
type Builder interface {
	Create() Builder

	WithKeyname(keyname string) Builder
	WithStyle(style templates.Template) Builder
	WithTemplate(template templates.Template) Builder

	AddEvent(event events.Event) Builder
	WithEvents(events []events.Event) Builder

	AddChild(child Component) Builder
	WithChildren(children []Component) Builder

	Now() (Component, error)
}

// Component represents a component.
type Component interface {
	renderables.StylableRenderable

	Keyname() string

	HasEvents() bool
	Events() []events.Event
	HasChildren() bool
	Children() []Component
}
