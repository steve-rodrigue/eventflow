package trees

import (
	"errors"
	"strings"
)

type routeBuilder struct {
	pattern string
	params  []Param
}

func (b *routeBuilder) Create() RouteBuilder {
	return &routeBuilder{}
}

func (b *routeBuilder) WithPattern(pattern string) RouteBuilder {
	b.pattern = strings.TrimSpace(pattern)
	return b
}
func (b *routeBuilder) AddParam(param Param) RouteBuilder {
	if param != nil {
		b.params = append(b.params, param)
	}
	return b
}
func (b *routeBuilder) WithParams(params []Param) RouteBuilder {
	b.params = params
	return b
}

func (b *routeBuilder) Now() (Route, error) {
	if b.pattern == "" {
		return nil, errors.New("route pattern is required")
	}
	return &route{pattern: b.pattern, params: b.params}, nil
}
