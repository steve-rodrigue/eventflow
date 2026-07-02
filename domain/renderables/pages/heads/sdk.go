package heads

import (
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

// NewRenderer creates a new head renderer
func NewRenderer(
	templateRenderer templates.Renderer,
	assetsRenderer assets.Renderer,
) Renderer {
	return &renderer{
		templateRenderer: templateRenderer,
		assetsRenderer:   assetsRenderer,
	}
}

// NewOpenGraphBuilder creates a new head builder
func NewBuilder() Builder {
	return &builder{}
}

// NewOpenGraphBuilder creates a new meta builder
func NewMetaBuilder() MetaBuilder {
	return &metaBuilder{}
}

// NewOpenGraphBuilder creates a new link builder
func NewLinkBuilder() LinkBuilder {
	return &linkBuilder{}
}

// NewOpenGraphBuilder creates a new open graph builder
func NewOpenGraphBuilder() OpenGraphBuilder {
	return &openGraphBuilder{}
}

// NewTwitterCardBuilder creates a new twitter card builder
func NewTwitterCardBuilder() TwitterCardBuilder {
	return &twitterCardBuilder{}
}

// Renderer renders the head content.
type Renderer interface {
	Render(head Head, params renderables.Params, assets assets.Assets) string
}

// Builder represents a page head builder.
type Builder interface {
	Create() Builder
	WithTitle(title templates.Template) Builder
	WithDescription(description templates.Template) Builder
	WithLanguage(language string) Builder
	AddMeta(meta Meta) Builder
	WithMeta(meta []Meta) Builder
	AddLink(link Link) Builder
	WithLinks(links []Link) Builder
	WithOpenGraph(openGraph OpenGraph) Builder
	WithTwitterCard(twitterCard TwitterCard) Builder
	Now() (Head, error)
}

// Head represents a page head.
type Head interface {
	Title() templates.Template
	Description() templates.Template
	Language() string

	Meta() []Meta
	Links() []Link

	HasOpenGraph() bool
	OpenGraph() OpenGraph

	HasTwitterCard() bool
	TwitterCard() TwitterCard
}

// MetaBuilder represents a meta builder.
type MetaBuilder interface {
	Create() MetaBuilder
	WithName(name string) MetaBuilder
	WithProperty(property string) MetaBuilder
	WithContent(content templates.Template) MetaBuilder
	Now() (Meta, error)
}

// Meta represents a meta tag.
type Meta interface {
	HasName() bool
	Name() string

	HasProperty() bool
	Property() string

	Content() templates.Template
}

// LinkBuilder represents a link builder.
type LinkBuilder interface {
	Create() LinkBuilder
	WithRel(rel string) LinkBuilder
	WithHref(href templates.Template) LinkBuilder
	WithType(linkType string) LinkBuilder
	Now() (Link, error)
}

// Link represents a link tag.
type Link interface {
	Rel() string
	Href() templates.Template

	HasType() bool
	Type() string
}

// OpenGraphBuilder represents an Open Graph builder.
type OpenGraphBuilder interface {
	Create() OpenGraphBuilder
	WithTitle(title templates.Template) OpenGraphBuilder
	WithDescription(description templates.Template) OpenGraphBuilder
	WithImage(image templates.Template) OpenGraphBuilder
	WithURL(url templates.Template) OpenGraphBuilder
	WithType(openGraphType string) OpenGraphBuilder
	WithSiteName(siteName templates.Template) OpenGraphBuilder
	Now() (OpenGraph, error)
}

// OpenGraph represents Open Graph metadata.
type OpenGraph interface {
	Title() templates.Template
	Description() templates.Template
	Image() templates.Template
	URL() templates.Template

	Type() string
	SiteName() templates.Template
}

// TwitterCardBuilder represents a Twitter/X card builder.
type TwitterCardBuilder interface {
	Create() TwitterCardBuilder
	WithCard(card string) TwitterCardBuilder
	WithTitle(title templates.Template) TwitterCardBuilder
	WithDescription(description templates.Template) TwitterCardBuilder
	WithImage(image templates.Template) TwitterCardBuilder
	Now() (TwitterCard, error)
}

// TwitterCard represents Twitter/X card metadata.
type TwitterCard interface {
	Card() string

	Title() templates.Template
	Description() templates.Template
	Image() templates.Template
}
