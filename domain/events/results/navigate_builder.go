package results

import (
	"errors"
	"strings"
)

type navigateBuilder struct {
	url string
}

func (b *navigateBuilder) Create() NavigateBuilder {
	return &navigateBuilder{}
}

func (b *navigateBuilder) WithURL(url string) NavigateBuilder {
	b.url = strings.TrimSpace(url)
	return b
}

func (b *navigateBuilder) Now() (Navigate, error) {
	if b.url == "" {
		return nil, errors.New("navigate url is required")
	}

	return &navigate{
		url: b.url,
	}, nil
}
