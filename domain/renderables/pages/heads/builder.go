package heads

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type builder struct {
	title       templates.Template
	description templates.Template
	language    string
	meta        []Meta
	links       []Link
	openGraph   OpenGraph
	twitterCard TwitterCard
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) WithTitle(title templates.Template) Builder {
	b.title = title
	return b
}

func (b *builder) WithDescription(description templates.Template) Builder {
	b.description = description
	return b
}

func (b *builder) WithLanguage(language string) Builder {
	b.language = strings.TrimSpace(language)
	return b
}

func (b *builder) AddMeta(meta Meta) Builder {
	if meta != nil {
		b.meta = append(b.meta, meta)
	}
	return b
}

func (b *builder) WithMeta(meta []Meta) Builder {
	b.meta = meta
	return b
}

func (b *builder) AddLink(link Link) Builder {
	if link != nil {
		b.links = append(b.links, link)
	}
	return b
}

func (b *builder) WithLinks(links []Link) Builder {
	b.links = links
	return b
}

func (b *builder) WithOpenGraph(openGraph OpenGraph) Builder {
	b.openGraph = openGraph
	return b
}

func (b *builder) WithTwitterCard(twitterCard TwitterCard) Builder {
	b.twitterCard = twitterCard
	return b
}

func (b *builder) Now() (Head, error) {
	if b.title == nil {
		return nil, errors.New("head title is required")
	}

	if b.description == nil {
		return nil, errors.New("head description is required")
	}

	if b.language == "" {
		return nil, errors.New("head language is required")
	}

	return &head{
		title:       b.title,
		description: b.description,
		language:    b.language,
		meta:        b.meta,
		links:       b.links,
		openGraph:   b.openGraph,
		twitterCard: b.twitterCard,
	}, nil
}
