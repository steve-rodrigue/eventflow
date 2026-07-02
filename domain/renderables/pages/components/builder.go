package components

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type builder struct {
	stylableBuilder renderables.StylableBuilder

	keyname  string
	style    templates.Template
	template templates.Template
	events   []events.Event
	children []Component
}

func (b *builder) Create() Builder {
	return &builder{
		stylableBuilder: b.stylableBuilder,
	}
}

func (b *builder) WithKeyname(keyname string) Builder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}

func (b *builder) WithStyle(style templates.Template) Builder {
	b.style = style
	return b
}

func (b *builder) WithTemplate(template templates.Template) Builder {
	b.template = template
	return b
}

func (b *builder) AddEvent(event events.Event) Builder {
	if event != nil {
		b.events = append(b.events, event)
	}

	return b
}

func (b *builder) WithEvents(events []events.Event) Builder {
	b.events = events
	return b
}

func (b *builder) AddChild(child Component) Builder {
	if child != nil {
		b.children = append(b.children, child)
	}

	return b
}

func (b *builder) WithChildren(children []Component) Builder {
	b.children = children
	return b
}

func (b *builder) Now() (Component, error) {
	if b.keyname == "" {
		return nil, errors.New("component keyname is required")
	}

	if b.stylableBuilder == nil {
		return nil, errors.New("component stylable builder is required")
	}

	stylableRenderable, err := b.stylableBuilder.
		Create().
		WithTemplate(b.template).
		WithStyle(b.style).
		Now()

	if err != nil {
		return nil, err
	}

	return &component{
		StylableRenderable: stylableRenderable,
		keyname:            b.keyname,
		events:             b.events,
		children:           b.children,
	}, nil
}
