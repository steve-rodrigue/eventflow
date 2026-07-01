package templates

// NewMustacheRenderer creates a new mustache template renderer.
func NewMustacheRenderer() Renderer {
	return &mustacheRenderer{}
}

// NewMustacheBuilder creates a new mustache template builder.
func NewMustacheBuilder() Builder {
	return &mustacheBuilder{}
}

// Renderer represents a template renderer
type Renderer interface {
	Render(template Template, values map[string]string) string
}

// Builder represents a template builder
type Builder interface {
	Create() Builder
	WithKeyname(keyname string) Builder
	WithCode(code string) Builder
	Now() (Template, error)
}

// Template represents a template
type Template interface {
	Keyname() string
	Code() string
	Contains(name string) bool
	HasParams() bool
	Params() []string
}
