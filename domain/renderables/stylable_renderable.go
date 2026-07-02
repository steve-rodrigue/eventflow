package renderables

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"

type stylableRenderable struct {
	renderable

	style templates.Template
}

func (r *stylableRenderable) HasStyle() bool {
	return r.style != nil
}

func (r *stylableRenderable) Style() templates.Template {
	return r.style
}
