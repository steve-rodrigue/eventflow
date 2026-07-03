package routers

import "github.com/steve-rodrigue/eventflow/domain/trees"

// NewRouter creates a new router.
func NewRouter() Router {
	return &router{}
}

// NewBuilder creates a new route request builder.
func NewBuilder() Builder {
	return &builder{}
}

// Router generates URIs from the tree.
type Router interface {
	URI(tree trees.Tree, request Request) (string, error)
}

// Builder represents a route request builder.
type Builder interface {
	Create() Builder

	WithTargetKeyname(targetKeyname string) Builder
	WithGroupKeyname(groupKeyname string) Builder
	WithLocale(locale string) Builder

	AddParam(param trees.Param) Builder
	WithParams(params []trees.Param) Builder

	Now() (Request, error)
}

// Request represents a request to generate a URI.
type Request interface {
	TargetKeyname() string
	GroupKeyname() string
	Locale() string

	HasParams() bool
	Params() []trees.Param
}
