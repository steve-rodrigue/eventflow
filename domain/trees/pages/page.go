package pages

import "slices"

type page struct {
	httpCode int
	headers  []Header
	body     string
}

func (p *page) HttpCode() int {
	return p.httpCode
}

func (p *page) Headers() []Header {
	return slices.Clone(p.headers)
}

func (p *page) Body() string {
	return p.body
}
