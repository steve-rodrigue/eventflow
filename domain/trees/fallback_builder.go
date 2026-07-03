package trees

import (
	"errors"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages"
)

type fallbackBuilder struct {
	httpCode int
	page     pages.Page
}

func (b *fallbackBuilder) Create() FallbackBuilder {
	return &fallbackBuilder{}
}

func (b *fallbackBuilder) WithHttpCode(httpCode int) FallbackBuilder {
	b.httpCode = httpCode
	return b
}

func (b *fallbackBuilder) WithPage(page pages.Page) FallbackBuilder {
	b.page = page
	return b
}

func (b *fallbackBuilder) Now() (Fallback, error) {
	if b.httpCode < 100 || b.httpCode > 599 {
		return nil, errors.New("fallback http code is invalid")
	}
	if b.page == nil {
		return nil, errors.New("fallback page is required")
	}
	return &fallback{httpCode: b.httpCode, page: b.page}, nil
}
