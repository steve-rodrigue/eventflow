package pages

import (
	"errors"
	"strings"
)

type headerBuilder struct {
	name  string
	value string
}

func (b *headerBuilder) Create() HeaderBuilder {
	return &headerBuilder{}
}

func (b *headerBuilder) WithName(name string) HeaderBuilder {
	b.name = strings.TrimSpace(name)
	return b
}

func (b *headerBuilder) WithValue(value string) HeaderBuilder {
	b.value = strings.TrimSpace(value)
	return b
}

func (b *headerBuilder) Now() (Header, error) {
	if b.name == "" {
		return nil, errors.New("header name is required")
	}

	return &header{
		name:  b.name,
		value: b.value,
	}, nil
}
