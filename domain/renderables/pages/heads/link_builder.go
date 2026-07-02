package heads

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type linkBuilder struct {
	rel      string
	href     templates.Template
	linkType string
}

func (b *linkBuilder) Create() LinkBuilder {
	return &linkBuilder{}
}

func (b *linkBuilder) WithRel(rel string) LinkBuilder {
	b.rel = strings.TrimSpace(rel)
	return b
}

func (b *linkBuilder) WithHref(href templates.Template) LinkBuilder {
	b.href = href
	return b
}

func (b *linkBuilder) WithType(linkType string) LinkBuilder {
	b.linkType = strings.TrimSpace(linkType)
	return b
}

func (b *linkBuilder) Now() (Link, error) {
	if b.rel == "" {
		return nil, errors.New("link rel is required")
	}

	if b.href == nil {
		return nil, errors.New("link href is required")
	}

	return &link{
		rel:      b.rel,
		href:     b.href,
		linkType: b.linkType,
	}, nil
}
