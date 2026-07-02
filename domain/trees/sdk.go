package trees

import (
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages"
	rendered_pages "github.com/steve-rodrigue/eventflow/domain/trees/pages"
)

// Renderer resolves a request into a rendered page.
type Renderer interface {
	Render(tree Tree, request Request) (rendered_pages.Page, error)
}

// RequestBuilder represents a request builder.
type RequestBuilder interface {
	Create() RequestBuilder
	WithPath(path string) RequestBuilder
	WithMethod(method string) RequestBuilder
	WithLocale(locale string) RequestBuilder
	WithTarget(target string) RequestBuilder
	Now() (Request, error)
}

// Request represents an incoming request.
type Request interface {
	Path() string
	Method() string
	Locale() string
	Target() string
}

// Builder represents a page tree builder.
type Builder interface {
	Create() Builder
	WithKeyname(keyname string) Builder
	AddNode(node Node) Builder
	WithNodes(nodes []Node) Builder
	Now() (Tree, error)
}

// Tree represents a page tree.
type Tree interface {
	Keyname() string
	Nodes() []Node
}

// NodeBuilder represents a node builder.
type NodeBuilder interface {
	Create() NodeBuilder
	WithKeyname(keyname string) NodeBuilder

	AddTarget(target Target) NodeBuilder
	WithTargets(targets []Target) NodeBuilder

	AddFallback(fallback Fallback) NodeBuilder
	WithFallbacks(fallbacks []Fallback) NodeBuilder

	Now() (Node, error)
}

// Node represents a logical page/resource group.
type Node interface {
	Keyname() string

	HasTargets() bool
	Targets() []Target

	HasFallbacks() bool
	Fallbacks() []Fallback

	Resource(locale string, targetKeyname string, groupKeyname string) (Resource, bool)
}

// TargetBuilder represents a target builder.
type TargetBuilder interface {
	Create() TargetBuilder
	WithKeyname(keyname string) TargetBuilder

	AddGroup(group Group) TargetBuilder
	WithGroups(groups []Group) TargetBuilder

	Now() (Target, error)
}

// Target represents a rendering target (desktop, mobile, print, amp...).
type Target interface {
	Keyname() string

	HasGroups() bool
	Groups() []Group

	Group(keyname string) (Group, bool)
}

// GroupBuilder represents a group builder.
type GroupBuilder interface {
	Create() GroupBuilder
	WithKeyname(keyname string) GroupBuilder

	AddResource(resource Resource) GroupBuilder
	WithResources(resources []Resource) GroupBuilder

	Now() (Group, error)
}

// Group represents a localized resource group.
type Group interface {
	Keyname() string

	HasResources() bool
	Resources() []Resource

	Resource(locale string) (Resource, bool)
}

// FallbackBuilder represents a fallback builder.
type FallbackBuilder interface {
	Create() FallbackBuilder
	WithHttpCode(httpCode int) FallbackBuilder
	WithPage(page pages.Page) FallbackBuilder
	Now() (Fallback, error)
}

// Fallback represents a non-localized fallback response.
type Fallback interface {
	HttpCode() int
	Page() pages.Page
}

// ResourceBuilder represents a resource builder.
type ResourceBuilder interface {
	Create() ResourceBuilder
	WithLocale(locale string) ResourceBuilder
	WithRoute(route Route) ResourceBuilder
	WithPage(page pages.Page) ResourceBuilder
	WithChildren(children Tree) ResourceBuilder
	Now() (Resource, error)
}

// Resource represents a localized resource.
type Resource interface {
	Locale() string

	Route() Route
	Page() pages.Page

	HasChildren() bool
	Children() Tree
}

// RouteBuilder represents a route builder.
type RouteBuilder interface {
	Create() RouteBuilder
	WithPattern(pattern string) RouteBuilder
	AddParam(param Param) RouteBuilder
	WithParams(params []Param) RouteBuilder
	Now() (Route, error)
}

// Route represents a route.
type Route interface {
	Pattern() string
	HasParams() bool
	Params() []Param
}

// ParamBuilder represents a route param builder.
type ParamBuilder interface {
	Create() ParamBuilder
	WithKeyname(keyname string) ParamBuilder
	WithValue(value string) ParamBuilder
	Now() (Param, error)
}

// Param represents a route param.
type Param interface {
	Keyname() string
	Value() string
}
