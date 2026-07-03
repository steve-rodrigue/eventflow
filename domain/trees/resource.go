package trees

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages"

type resource struct {
	locale   string
	route    Route
	page     pages.Page
	children Tree
}

func (r *resource) Locale() string {
	return r.locale
}

func (r *resource) Route() Route {
	return r.route
}

func (r *resource) Page() pages.Page {
	return r.page
}

func (r *resource) HasChildren() bool {
	return r.children != nil
}

func (r *resource) Children() Tree {
	return r.children
}
