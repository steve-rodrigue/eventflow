package heads

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"

type openGraph struct {
	title       templates.Template
	description templates.Template
	image       templates.Template
	url         templates.Template
	graphType   string
	siteName    templates.Template
}

func (o *openGraph) Title() templates.Template {
	return o.title
}

func (o *openGraph) Description() templates.Template {
	return o.description
}

func (o *openGraph) Image() templates.Template {
	return o.image
}

func (o *openGraph) URL() templates.Template {
	return o.url
}

func (o *openGraph) Type() string {
	return o.graphType
}

func (o *openGraph) SiteName() templates.Template {
	return o.siteName
}
