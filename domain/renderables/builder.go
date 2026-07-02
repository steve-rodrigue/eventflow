package renderables

import (
	"errors"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type builder struct {
	template templates.Template
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) WithTemplate(template templates.Template) Builder {
	b.template = template
	return b
}

func (b *builder) Now() (Renderable, error) {
	if b.template == nil {
		return nil, errors.New("renderable template is required")
	}

	return &renderable{
		template: b.template,
	}, nil
}
