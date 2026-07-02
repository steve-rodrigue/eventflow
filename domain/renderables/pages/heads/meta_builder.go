package heads

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type metaBuilder struct {
	name     string
	property string
	content  templates.Template
}

func (b *metaBuilder) Create() MetaBuilder {
	return &metaBuilder{}
}

func (b *metaBuilder) WithName(name string) MetaBuilder {
	b.name = strings.TrimSpace(name)
	return b
}

func (b *metaBuilder) WithProperty(property string) MetaBuilder {
	b.property = strings.TrimSpace(property)
	return b
}

func (b *metaBuilder) WithContent(content templates.Template) MetaBuilder {
	b.content = content
	return b
}

func (b *metaBuilder) Now() (Meta, error) {
	if b.name == "" && b.property == "" {
		return nil, errors.New("meta name or property is required")
	}

	if b.name != "" && b.property != "" {
		return nil, errors.New("meta cannot have both name and property")
	}

	if b.content == nil {
		return nil, errors.New("meta content is required")
	}

	return &meta{
		name:     b.name,
		property: b.property,
		content:  b.content,
	}, nil
}
