package trees

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages"

type fallback struct {
	httpCode int
	page     pages.Page
}

func (f *fallback) HttpCode() int {
	return f.httpCode
}

func (f *fallback) Page() pages.Page {
	return f.page
}
