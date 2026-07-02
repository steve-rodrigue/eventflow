package heads

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"

type twitterCard struct {
	card        string
	title       templates.Template
	description templates.Template
	image       templates.Template
}

func (t *twitterCard) Card() string {
	return t.card
}

func (t *twitterCard) Title() templates.Template {
	return t.title
}

func (t *twitterCard) Description() templates.Template {
	return t.description
}

func (t *twitterCard) Image() templates.Template {
	return t.image
}
