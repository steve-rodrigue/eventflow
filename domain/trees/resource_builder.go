package trees

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages"
)

type resourceBuilder struct {
	locale   string
	route    Route
	page     pages.Page
	children Tree
}

func (b *resourceBuilder) Create() ResourceBuilder {
	return &resourceBuilder{}
}

func (b *resourceBuilder) WithLocale(locale string) ResourceBuilder {
	b.locale = strings.TrimSpace(locale)
	return b
}

func (b *resourceBuilder) WithRoute(route Route) ResourceBuilder {
	b.route = route
	return b
}

func (b *resourceBuilder) WithPage(page pages.Page) ResourceBuilder {
	b.page = page
	return b
}

func (b *resourceBuilder) WithChildren(children Tree) ResourceBuilder {
	b.children = children
	return b
}

func (b *resourceBuilder) Now() (Resource, error) {
	if b.locale == "" {
		return nil, errors.New("resource locale is required")
	}
	if b.route == nil {
		return nil, errors.New("resource route is required")
	}
	if b.page == nil {
		return nil, errors.New("resource page is required")
	}

	return &resource{locale: b.locale, route: b.route, page: b.page, children: b.children}, nil
}
