package pages

import (
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/components"
)

// Page represents a page
type Page interface {
	renderables.Renderable

	Keyname() string
	Head() Head
	Body() Body
}

// Head represents a page head
type Head interface {
	renderables.Renderable

	Title() string
	Description() string
}

// Body represents a page body
type Body interface {
	renderables.StylableRenderable

	Components() []components.Component
}
