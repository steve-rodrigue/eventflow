package pages

import (
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

// NewRenderer creates a new page renderer
func NewRenderer(
	templateRenderer templates.Renderer,
	headRenderer heads.Renderer,
	componentRenderer components.Renderer,
) Renderer {
	return &renderer{
		templateRenderer:  templateRenderer,
		headRenderer:      headRenderer,
		componentRenderer: componentRenderer,
	}
}

// NewBuilder creates a new page builder
func NewBuilder(renderableBuilder renderables.Builder) Builder {
	return &builder{
		renderableBuilder: renderableBuilder,
	}
}

// Renderer renders the full page.
type Renderer interface {
	Render(page Page, params renderables.Params, assets assets.Assets) string
	RenderStyle(page Page, params renderables.Params) string
}

// Builder represents a page builder.
type Builder interface {
	Create() Builder
	WithLanguage(language string) Builder
	WithKeyname(keyname string) Builder
	WithTemplate(template templates.Template) Builder
	WithHead(head heads.Head) Builder
	WithBody(body components.Component) Builder
	Now() (Page, error)
}

// Page represents a page.
type Page interface {
	renderables.Renderable

	Language() string
	Head() heads.Head
	Keyname() string
	Body() components.Component
}
