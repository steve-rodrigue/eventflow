package trees

import (
	"errors"
	"strings"
)

type requestBuilder struct {
	path   string
	method string
	locale string
	target string
}

func (b *requestBuilder) Create() RequestBuilder { return &requestBuilder{} }
func (b *requestBuilder) WithPath(path string) RequestBuilder {
	b.path = strings.TrimSpace(path)
	return b
}
func (b *requestBuilder) WithMethod(method string) RequestBuilder {
	b.method = strings.TrimSpace(method)
	return b
}
func (b *requestBuilder) WithLocale(locale string) RequestBuilder {
	b.locale = strings.TrimSpace(locale)
	return b
}
func (b *requestBuilder) WithTarget(target string) RequestBuilder {
	b.target = strings.TrimSpace(target)
	return b
}

func (b *requestBuilder) Now() (Request, error) {
	if b.path == "" {
		return nil, errors.New("request path is required")
	}
	if b.method == "" {
		return nil, errors.New("request method is required")
	}
	if b.locale == "" {
		return nil, errors.New("request locale is required")
	}
	if b.target == "" {
		return nil, errors.New("request target is required")
	}

	return &request{path: b.path, method: b.method, locale: b.locale, target: b.target}, nil
}
