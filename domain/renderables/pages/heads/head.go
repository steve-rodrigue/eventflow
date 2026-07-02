package heads

import (
	"slices"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type head struct {
	renderables.Renderable

	title       templates.Template
	description templates.Template
	language    string
	meta        []Meta
	links       []Link
	openGraph   OpenGraph
	twitterCard TwitterCard
}

func (h *head) Title() templates.Template {
	return h.title
}

func (h *head) Description() templates.Template {
	return h.description
}

func (h *head) Language() string {
	return h.language
}

func (h *head) Meta() []Meta {
	return slices.Clone(h.meta)
}

func (h *head) Links() []Link {
	return slices.Clone(h.links)
}

func (h *head) HasOpenGraph() bool {
	return h.openGraph != nil
}

func (h *head) OpenGraph() OpenGraph {
	return h.openGraph
}

func (h *head) HasTwitterCard() bool {
	return h.twitterCard != nil
}

func (h *head) TwitterCard() TwitterCard {
	return h.twitterCard
}
