package heads

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"

type link struct {
	rel      string
	href     templates.Template
	linkType string
}

func (l *link) Rel() string {
	return l.rel
}

func (l *link) Href() templates.Template {
	return l.href
}

func (l *link) HasType() bool {
	return l.linkType != ""
}

func (l *link) Type() string {
	return l.linkType
}
