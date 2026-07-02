package heads

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type openGraphBuilder struct {
	title       templates.Template
	description templates.Template
	image       templates.Template
	url         templates.Template
	graphType   string
	siteName    templates.Template
}

func (b *openGraphBuilder) Create() OpenGraphBuilder {
	return &openGraphBuilder{}
}

func (b *openGraphBuilder) WithTitle(title templates.Template) OpenGraphBuilder {
	b.title = title
	return b
}

func (b *openGraphBuilder) WithDescription(description templates.Template) OpenGraphBuilder {
	b.description = description
	return b
}

func (b *openGraphBuilder) WithImage(image templates.Template) OpenGraphBuilder {
	b.image = image
	return b
}

func (b *openGraphBuilder) WithURL(url templates.Template) OpenGraphBuilder {
	b.url = url
	return b
}

func (b *openGraphBuilder) WithType(openGraphType string) OpenGraphBuilder {
	b.graphType = strings.TrimSpace(openGraphType)
	return b
}

func (b *openGraphBuilder) WithSiteName(siteName templates.Template) OpenGraphBuilder {
	b.siteName = siteName
	return b
}

func (b *openGraphBuilder) Now() (OpenGraph, error) {
	if b.title == nil {
		return nil, errors.New("open graph title is required")
	}

	if b.description == nil {
		return nil, errors.New("open graph description is required")
	}

	if b.image == nil {
		return nil, errors.New("open graph image is required")
	}

	if b.url == nil {
		return nil, errors.New("open graph url is required")
	}

	if b.graphType == "" {
		return nil, errors.New("open graph type is required")
	}

	if b.siteName == nil {
		return nil, errors.New("open graph site name is required")
	}

	return &openGraph{
		title:       b.title,
		description: b.description,
		image:       b.image,
		url:         b.url,
		graphType:   b.graphType,
		siteName:    b.siteName,
	}, nil
}
