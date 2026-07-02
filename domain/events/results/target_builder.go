package results

import (
	"errors"
	"strings"
)

type targetBuilder struct {
	element string
	html    string
}

func (b *targetBuilder) Create() TargetBuilder {
	return &targetBuilder{}
}

func (b *targetBuilder) WithElement(element string) TargetBuilder {
	b.element = strings.TrimSpace(element)
	return b
}

func (b *targetBuilder) WithHTML(html string) TargetBuilder {
	b.html = html
	return b
}

func (b *targetBuilder) Now() (Target, error) {
	if b.element == "" {
		return nil, errors.New("target element is required")
	}

	return &target{
		element: b.element,
		html:    b.html,
	}, nil
}
