package trees

import "github.com/steve-rodrigue/eventflow/domain/renderables/pages"

// Tree represents a page tree
type Tree interface {
	Keyname() string
	Nodes() []Node
}

// Node represents a tree node
type Node interface {
	Keyname() string
	Generics() []Generic
}

// Generic represents a generic node
type Generic interface {
	HttpCode() int
	Page() pages.Page
	Resources() []Resource
}

// Resource represents a resource
type Resource interface {
	Link() Link
	Page() pages.Page
	HasChildren() bool
	Children() Tree
}

// Link represents a link
type Link interface {
	Keyname() string
	RoutePattern() string
	HasParams() bool
	Params() []Param
}

// Param represents a link param
type Param interface {
	Keyname() string
	Value() string
}
