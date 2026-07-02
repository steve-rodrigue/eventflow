package pages

import (
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
)

type page struct {
	renderables.Renderable

	language string
	keyname  string
	head     heads.Head
	body     components.Component
}

func (p *page) Language() string {
	return p.language
}

func (p *page) Keyname() string {
	return p.keyname
}

func (p *page) Head() heads.Head {
	return p.head
}

func (p *page) Body() components.Component {
	return p.body
}
