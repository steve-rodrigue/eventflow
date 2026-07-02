package heads

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type twitterCardBuilder struct {
	card        string
	title       templates.Template
	description templates.Template
	image       templates.Template
}

func (b *twitterCardBuilder) Create() TwitterCardBuilder {
	return &twitterCardBuilder{}
}

func (b *twitterCardBuilder) WithCard(card string) TwitterCardBuilder {
	b.card = strings.TrimSpace(card)
	return b
}

func (b *twitterCardBuilder) WithTitle(title templates.Template) TwitterCardBuilder {
	b.title = title
	return b
}

func (b *twitterCardBuilder) WithDescription(description templates.Template) TwitterCardBuilder {
	b.description = description
	return b
}

func (b *twitterCardBuilder) WithImage(image templates.Template) TwitterCardBuilder {
	b.image = image
	return b
}

func (b *twitterCardBuilder) Now() (TwitterCard, error) {
	if b.card == "" {
		return nil, errors.New("twitter card is required")
	}

	if b.title == nil {
		return nil, errors.New("twitter title is required")
	}

	if b.description == nil {
		return nil, errors.New("twitter description is required")
	}

	if b.image == nil {
		return nil, errors.New("twitter image is required")
	}

	return &twitterCard{
		card:        b.card,
		title:       b.title,
		description: b.description,
		image:       b.image,
	}, nil
}
