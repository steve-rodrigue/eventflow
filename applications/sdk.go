package applications

import (
	eventapps "github.com/steve-rodrigue/eventflow/applications/events"
	"github.com/steve-rodrigue/eventflow/domain/events"
	domainevents "github.com/steve-rodrigue/eventflow/domain/events"
	domaincontexts "github.com/steve-rodrigue/eventflow/domain/events/contexts"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
	domainresults "github.com/steve-rodrigue/eventflow/domain/events/results"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	renderablepages "github.com/steve-rodrigue/eventflow/domain/renderables/pages"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
	"github.com/steve-rodrigue/eventflow/domain/routers"
	"github.com/steve-rodrigue/eventflow/domain/trees"
	treepages "github.com/steve-rodrigue/eventflow/domain/trees/pages"
)

// NewDefaultApplication creates a new default application
func NewDefaultApplication() Application {
	templateRenderer := templates.NewMustacheRenderer()
	templateBuilder := templates.NewMustacheBuilder()

	eventRegistry := domainevents.NewRegistry()

	eventApplication := eventapps.NewApplication(
		domaincontexts.NewBuilder(),
		eventRegistry,
	)

	assetsRenderer := assets.NewRenderer()

	headRenderer := heads.NewRenderer(
		templateRenderer,
		assetsRenderer,
	)

	componentRenderer := components.NewRenderer(
		templateRenderer,
	)

	pageRenderer := renderablepages.NewRenderer(
		templateRenderer,
		headRenderer,
		componentRenderer,
	)

	treeRenderer := trees.NewRenderer(
		pageRenderer,
		treepages.NewBuilder(),
		treepages.NewHeaderBuilder(),
		assets.NewBuilder(),
		assets.NewAssetBuilder(),
	)

	return NewApplication(
		templateBuilder,

		eventRegistry,
		eventApplication,
		domainevents.NewBuilder(),

		domainresults.NewBuilder(),
		domainresults.NewOperationBuilder(),
		domainresults.NewActionBuilder(),
		domainresults.NewTargetBuilder(),
		domainresults.NewNavigateBuilder(),

		routers.NewRouter(),
		routers.NewBuilder(),

		treeRenderer,
		trees.NewRequestBuilder(),
		trees.NewBuilder(),
		trees.NewNodeBuilder(),
		trees.NewTargetBuilder(),
		trees.NewGroupBuilder(),
		trees.NewFallbackBuilder(),
		trees.NewResourceBuilder(),
		trees.NewRouteBuilder(),
		trees.NewParamBuilder(),

		renderables.NewBuilder(),
		renderables.NewStylableBuilder(),

		renderablepages.NewBuilder(renderables.NewBuilder()),
		heads.NewBuilder(),
		heads.NewMetaBuilder(),
		heads.NewLinkBuilder(),
		heads.NewOpenGraphBuilder(),
		heads.NewTwitterCardBuilder(),
		components.NewBuilder(renderables.NewStylableBuilder()),
	)
}

// NewApplication creates a new application
func NewApplication(
	templateBuilder templates.Builder,

	eventRegistry domainevents.Registry,
	eventApplication eventapps.Application,
	eventBuilder domainevents.Builder,

	resultBuilder domainresults.Builder,
	operationBuilder domainresults.OperationBuilder,
	actionBuilder domainresults.ActionBuilder,
	targetActionBuilder domainresults.TargetBuilder,
	navigateBuilder domainresults.NavigateBuilder,

	router routers.Router,
	routerRequestBuilder routers.Builder,

	treeRenderer trees.Renderer,
	treeRequestBuilder trees.RequestBuilder,
	treeBuilder trees.Builder,
	nodeBuilder trees.NodeBuilder,
	targetBuilder trees.TargetBuilder,
	groupBuilder trees.GroupBuilder,
	fallbackBuilder trees.FallbackBuilder,
	resourceBuilder trees.ResourceBuilder,
	routeBuilder trees.RouteBuilder,
	paramBuilder trees.ParamBuilder,

	renderableBuilder renderables.Builder,
	stylableBuilder renderables.StylableBuilder,

	pageBuilder renderablepages.Builder,
	headBuilder heads.Builder,
	metaBuilder heads.MetaBuilder,
	linkBuilder heads.LinkBuilder,
	openGraphBuilder heads.OpenGraphBuilder,
	twitterCardBuilder heads.TwitterCardBuilder,
	componentBuilder components.Builder,
) Application {
	return &application{
		templateBuilder: templateBuilder,

		eventRegistry:    eventRegistry,
		eventApplication: eventApplication,
		eventBuilder:     eventBuilder,
		resultBuilder:    resultBuilder,
		operationBuilder: operationBuilder,
		actionBuilder:    actionBuilder,
		targetBuilder:    targetActionBuilder,
		navigateBuilder:  navigateBuilder,

		router:               router,
		routerRequestBuilder: routerRequestBuilder,

		treeRenderer:       treeRenderer,
		treeRequestBuilder: treeRequestBuilder,
		treeBuilder:        treeBuilder,
		nodeBuilder:        nodeBuilder,
		targetTreeBuilder:  targetBuilder,
		groupBuilder:       groupBuilder,
		fallbackBuilder:    fallbackBuilder,
		resourceBuilder:    resourceBuilder,
		routeBuilder:       routeBuilder,
		paramBuilder:       paramBuilder,

		renderableBuilder: renderableBuilder,
		stylableBuilder:   stylableBuilder,

		pageBuilder:        pageBuilder,
		headBuilder:        headBuilder,
		metaBuilder:        metaBuilder,
		linkBuilder:        linkBuilder,
		openGraphBuilder:   openGraphBuilder,
		twitterCardBuilder: twitterCardBuilder,
		componentBuilder:   componentBuilder,
	}
}

/*
Tree represents the root of an application page hierarchy.
*/
type Tree struct {
	Keyname string
	Nodes   []Node
}

/*
Node represents a logical section of the application.

A node can expose multiple rendering targets (desktop, mobile,
print, pdf, etc.) and optional fallback pages (404, 500, etc.).
*/
type Node struct {
	Targets   []Target
	Fallbacks []Fallback
}

/*
Target represents a rendering target.

Typical keynames include:
- desktop
- mobile
- tablet
- print
- pdf
*/
type Target struct {
	Keyname string
	Groups  []Group
}

/*
Group represents the same logical resource in multiple locales.

For example, a "posts" group may contain:
- en -> /posts
- fr -> /articles
- es -> /articulos

A group cannot contain two resources with the same locale.
*/
type Group struct {
	Keyname   string
	Resources []Resource
}

/*
Fallback represents a page returned for a specific HTTP status.

Typical status codes include:
- 404
- 401
- 403
- 500
*/
type Fallback struct {
	HttpCode int
	Page     Page
}

/*
Resource represents a localized page.

Each resource belongs to a locale and may optionally contain a child
tree, allowing applications to build hierarchical navigation.
*/
type Resource struct {
	Locale string

	Route    Route
	Page     Page
	Children *Tree
}

/*
Route represents the URL pattern used to access a resource.

Examples:
- /
- /posts
- /posts/{id}
- /users/{userId}/posts/{postId}
*/
type Route struct {
	Pattern string
	Params  []Param
}

/*
Param represents a route parameter.
*/
type Param struct {
	Keyname string
	Value   string
}

/*
Page represents a renderable HTML page.
*/
type Page struct {
	Renderable

	Language string
	Head     Head
	Keyname  string
	Body     Component
}

/*
Head represents the HTML <head> section.
*/
type Head struct {
	Title       Template
	Description Template
	Language    string

	Meta        []Meta
	Links       []Link
	OpenGraph   *OpenGraph
	TwitterCard *TwitterCard
}

/*
Meta represents an HTML <meta> tag.

Specify either Name or Property depending on the type of metadata.
*/
type Meta struct {
	Name     string
	Property string
	Content  Template
}

/*
Link represents an HTML <link> tag.

Common relations include:
- stylesheet
- icon
- canonical
- preload
- preconnect
*/
type Link struct {
	Rel  string
	Href Template
	Type string
}

/*
OpenGraph represents metadata used by social platforms that support
the Open Graph protocol.
*/
type OpenGraph struct {
	Title       Template
	Description Template
	Image       Template
	URL         Template

	Type     string
	SiteName Template
}

/*
TwitterCard represents metadata used by Twitter/X.
*/
type TwitterCard struct {
	Card string

	Title       Template
	Description Template
	Image       Template
}

/*
Component represents a reusable UI component.

Components may contain child components and browser events.
*/
type Component struct {
	StylableRenderable

	Keyname  string
	Events   []Event
	Children []Component
}

/*
StylableRenderable represents a renderable object that optionally
contains CSS styles.
*/
type StylableRenderable struct {
	Renderable
	Style *Template
}

/*
Renderable represents an object rendered from a template.
*/
type Renderable struct {
	Template Template
}

/*
Template represents a parsed template.

Params contains the template variables referenced by the template.
*/
type Template struct {
	Keyname string
	Code    string
	Params  []string
}

/*
ActionFn represents a server-side event handler.

The function receives the event execution context and returns one or
more DOM operations that should be applied by the browser.
*/
type ActionFn func(ctx Context) (*Result, error)

/*
Result represents the response returned after executing an event.

A result may contain multiple operations that are executed by the
browser in the order they are defined.
*/
type Result struct {
	Operations []Operation
}

/*
Context represents the execution context of an event.

Payload contains the values sent by the browser when the event was
triggered.
*/
type Context struct {
	EventName string
	Payload   map[string]any
}

/*
Operation represents a single browser operation.

The operation type determines which action is executed (replace,
append, prepend, remove, navigate, etc.).
*/
type Operation struct {
	Type   results.OperationType
	Action Action
}

/*
Action represents the details of an operation.

Only one action should be specified depending on the operation type.
*/
type Action struct {
	Navigate *ActionNavigate
	Target   *ActionTarget
}

/*
ActionTarget represents an operation targeting a DOM element.

Element is a CSS selector used to locate the element.
HTML contains the HTML fragment used by the operation.
*/
type ActionTarget struct {
	Element string
	HTML    string
}

/*
ActionNavigate represents a browser navigation.

When executed, the browser navigates to the specified URL.
*/
type ActionNavigate struct {
	URL string
}

/*
Event represents executable server-side behavior.
*/
type Event struct {
	Keyname string
	Action  ActionFn
}

/*
BrowserEvent represents an event triggered directly by browser
interactions.

Typical browser events include:
- click
- input
- submit
- change
- keydown
- keyup
*/
type BrowserEvent struct {
	Event
	BrowserType events.BrowserType
}

/*
CustomEvent represents an application-defined event that can be
triggered programmatically.
*/
type CustomEvent struct {
	Event
	EventName string
}

/*
URIRequest represents a request to generate a URI from the application.

It identifies a localized resource by its rendering target, resource
group, locale, and optional route parameters.
*/
type URIRequest struct {
	TargetKeyname string
	GroupKeyname  string
	Locale        string
	Params        []Param
}

/*
RouteRequest represents an incoming HTTP request.

It is used to resolve a request path into a rendered page.
*/
type RouteRequest struct {
	Path   string
	Method string
	Locale string
	Target string
}

/*
RenderedPage represents the final HTTP response generated by the
application.

It contains the HTTP status code, response headers and fully rendered
HTML body.
*/
type RenderedPage struct {
	HttpCode int
	Headers  []RenderedHeader
	Body     string
}

/*
RenderedHeader represents an HTTP response header.
*/
type RenderedHeader struct {
	Name  string
	Value string
}

/*
Application executes an EventFlow application.

Execute initializes the application and prepares all internal
renderers, routers and event handlers.

URI generates the URI of a localized resource.

Route resolves an incoming HTTP request into a rendered page.

Trigger executes a browser or custom event and returns the DOM
operations that should be applied by the client.
*/
type Application interface {
	Execute(tree Tree) error
	URI(request URIRequest) (string, error)
	Route(request RouteRequest) (*RenderedPage, error)
	Trigger(msg eventapps.IncomingMessage) (*eventapps.OutgoingMessage, error)
}
