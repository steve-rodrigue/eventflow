package renderables

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"

type renderable struct {
	template templates.Template
}

func (r *renderable) Template() templates.Template {
	return r.template
}
