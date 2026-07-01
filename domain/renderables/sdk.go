package renderables

import "github.com/steve-rodrigue/eventflow/domain/templates"

// Params represents params
type Params map[string]any

// Renderer represents a renderable renderer.
type Renderer interface {
	Render(renderable Renderable, params Params) string
	RenderList(renderable Renderable, params []Params) string
	RenderWithStyle(renderable StylableRenderable, template Params, style Params) (string, string)
	RenderListWithStyle(renderable StylableRenderable, template []Params, style []Params) (string, string)
}

// Renderable represents an object that can be rendered with template.
type Renderable interface {
	Template() templates.Template
}

// StylableRenderable represents an object that can be rendered with style and template.
type StylableRenderable interface {
	Renderable

	HasStyle() bool
	Style() templates.Template
}
