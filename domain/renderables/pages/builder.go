package pages

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type builder struct {
	renderableBuilder renderables.Builder

	language string
	keyname  string
	template templates.Template
	head     heads.Head
	body     components.Component
}

func (b *builder) Create() Builder {
	return &builder{
		renderableBuilder: b.renderableBuilder,
	}
}

func (b *builder) WithLanguage(language string) Builder {
	b.language = strings.TrimSpace(language)
	return b
}

func (b *builder) WithKeyname(keyname string) Builder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}

func (b *builder) WithTemplate(template templates.Template) Builder {
	b.template = template
	return b
}

func (b *builder) WithHead(head heads.Head) Builder {
	b.head = head
	return b
}

func (b *builder) WithBody(body components.Component) Builder {
	b.body = body
	return b
}

func (b *builder) Now() (Page, error) {
	if b.renderableBuilder == nil {
		return nil, errors.New("page renderable builder is required")
	}

	if b.language == "" {
		return nil, errors.New("page language is required")
	}

	if b.keyname == "" {
		return nil, errors.New("page keyname is required")
	}

	if b.head == nil {
		return nil, errors.New("page head is required")
	}

	if b.body == nil {
		return nil, errors.New("page body is required")
	}

	renderable, err := b.renderableBuilder.
		Create().
		WithTemplate(b.template).
		Now()

	if err != nil {
		return nil, err
	}

	return &page{
		Renderable: renderable,
		language:   b.language,
		keyname:    b.keyname,
		head:       b.head,
		body:       b.body,
	}, nil
}
