package pages

import (
	"errors"
)

type builder struct {
	httpCode int
	headers  []Header
	body     string
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) WithHttpCode(httpCode int) Builder {
	b.httpCode = httpCode
	return b
}

func (b *builder) AddHeader(header Header) Builder {
	if header != nil {
		b.headers = append(b.headers, header)
	}

	return b
}

func (b *builder) WithHeaders(headers []Header) Builder {
	b.headers = headers
	return b
}

func (b *builder) WithBody(body string) Builder {
	b.body = body
	return b
}

func (b *builder) Now() (Page, error) {
	if b.httpCode < 100 || b.httpCode > 599 {
		return nil, errors.New("page http code is invalid")
	}

	return &page{
		httpCode: b.httpCode,
		headers:  b.headers,
		body:     b.body,
	}, nil
}
