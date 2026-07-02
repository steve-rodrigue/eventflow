package renderables

import (
	"errors"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type stylableBuilder struct {
	template templates.Template
	style    templates.Template
}

func (b *stylableBuilder) Create() StylableBuilder {
	return &stylableBuilder{}
}

func (b *stylableBuilder) WithTemplate(template templates.Template) StylableBuilder {
	b.template = template
	return b
}

func (b *stylableBuilder) WithStyle(style templates.Template) StylableBuilder {
	b.style = style
	return b
}

func (b *stylableBuilder) Now() (StylableRenderable, error) {
	if b.template == nil {
		return nil, errors.New("renderable template is required")
	}

	return &stylableRenderable{
		renderable: renderable{
			template: b.template,
		},
		style: b.style,
	}, nil
}
