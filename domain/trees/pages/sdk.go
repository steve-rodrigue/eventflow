package pages

// NewBuilder creates a new page builder.
func NewBuilder() Builder {
	return &builder{}
}

// NewHeaderBuilder creates a new header builder.
func NewHeaderBuilder() HeaderBuilder {
	return &headerBuilder{}
}

// Builder represents a rendered page builder.
type Builder interface {
	Create() Builder
	WithHttpCode(httpCode int) Builder
	AddHeader(header Header) Builder
	WithHeaders(headers []Header) Builder
	WithBody(body string) Builder
	Now() (Page, error)
}

// Page represents a fully rendered page.
type Page interface {
	HttpCode() int
	Headers() []Header
	Body() string
}

// HeaderBuilder represents a header builder.
type HeaderBuilder interface {
	Create() HeaderBuilder
	WithName(name string) HeaderBuilder
	WithValue(value string) HeaderBuilder
	Now() (Header, error)
}

// Header represents a page header
type Header interface {
	Name() string
	Value() string
}
