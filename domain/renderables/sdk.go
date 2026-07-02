package renderables

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"

// Params represents params
type Params map[string]any

// NewBuilder creates a new renderable builder.
func NewBuilder() Builder {
	return &builder{}
}

// NewStylableBuilder creates a new stylable renderable builder.
func NewStylableBuilder() StylableBuilder {
	return &stylableBuilder{}
}

// Builder represents a renderable builder.
type Builder interface {
	Create() Builder
	WithTemplate(template templates.Template) Builder
	Now() (Renderable, error)
}

// Renderable represents an object that can be rendered with template.
type Renderable interface {
	Template() templates.Template
}

// StylableBuilder represents a stylable renderable builder.
type StylableBuilder interface {
	Create() StylableBuilder
	WithTemplate(template templates.Template) StylableBuilder
	WithStyle(style templates.Template) StylableBuilder
	Now() (StylableRenderable, error)
}

// StylableRenderable represents an object that can be rendered with style and template.
type StylableRenderable interface {
	Renderable

	HasStyle() bool
	Style() templates.Template
}
