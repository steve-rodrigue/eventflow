package heads

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"

type meta struct {
	name     string
	property string
	content  templates.Template
}

func (m *meta) HasName() bool {
	return m.name != ""
}

func (m *meta) Name() string {
	return m.name
}

func (m *meta) HasProperty() bool {
	return m.property != ""
}

func (m *meta) Property() string {
	return m.property
}

func (m *meta) Content() templates.Template {
	return m.content
}
