package applications

import (
	"errors"
	"strings"

	eventapps "github.com/steve-rodrigue/eventflow/applications/events"
	domainevents "github.com/steve-rodrigue/eventflow/domain/events"
	domaincontexts "github.com/steve-rodrigue/eventflow/domain/events/contexts"
	domainresults "github.com/steve-rodrigue/eventflow/domain/events/results"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	renderablepages "github.com/steve-rodrigue/eventflow/domain/renderables/pages"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
	"github.com/steve-rodrigue/eventflow/domain/routers"
	"github.com/steve-rodrigue/eventflow/domain/trees"
	treepages "github.com/steve-rodrigue/eventflow/domain/trees/pages"
)

type application struct {
	tree trees.Tree

	templateBuilder templates.Builder

	eventRegistry    domainevents.Registry
	eventApplication eventapps.Application
	eventBuilder     domainevents.Builder

	resultBuilder    domainresults.Builder
	operationBuilder domainresults.OperationBuilder
	actionBuilder    domainresults.ActionBuilder
	targetBuilder    domainresults.TargetBuilder
	navigateBuilder  domainresults.NavigateBuilder

	router               routers.Router
	routerRequestBuilder routers.Builder

	treeRenderer       trees.Renderer
	treeRequestBuilder trees.RequestBuilder
	pageRenderer       renderablepages.Renderer

	treeBuilder       trees.Builder
	nodeBuilder       trees.NodeBuilder
	targetTreeBuilder trees.TargetBuilder
	groupBuilder      trees.GroupBuilder
	fallbackBuilder   trees.FallbackBuilder
	resourceBuilder   trees.ResourceBuilder
	routeBuilder      trees.RouteBuilder
	paramBuilder      trees.ParamBuilder

	renderableBuilder renderables.Builder
	stylableBuilder   renderables.StylableBuilder

	pageBuilder        renderablepages.Builder
	headBuilder        heads.Builder
	metaBuilder        heads.MetaBuilder
	linkBuilder        heads.LinkBuilder
	openGraphBuilder   heads.OpenGraphBuilder
	twitterCardBuilder heads.TwitterCardBuilder
	componentBuilder   components.Builder
}

func (app *application) Initialize(tree Tree) error {
	domainTree, err := app.toTree(tree)
	if err != nil {
		return err
	}

	app.tree = domainTree
	if err := app.registerTreeEvents(tree); err != nil {
		return err
	}
	return nil
}

func (app *application) URI(request URIRequest) (string, error) {
	if app.tree == nil {
		return "", errors.New("application tree is not initialized")
	}

	params, err := app.toParams(request.Params)
	if err != nil {
		return "", err
	}

	routeRequest, err := app.routerRequestBuilder.
		Create().
		WithTargetKeyname(request.TargetKeyname).
		WithGroupKeyname(request.GroupKeyname).
		WithLocale(request.Locale).
		WithParams(params).
		Now()

	if err != nil {
		return "", err
	}

	return app.router.URI(app.tree, routeRequest)
}

func (app *application) ResolveURI(request URIRequest) (string, error) {
	if app.tree == nil {
		return "", errors.New("application tree is not initialized")
	}

	params, err := app.toParams(request.Params)
	if err != nil {
		return "", err
	}

	routeRequest, err := app.routerRequestBuilder.
		Create().
		WithTargetKeyname(request.TargetKeyname).
		WithGroupKeyname(request.GroupKeyname).
		WithLocale(request.Locale).
		WithParams(params).
		Now()

	if err != nil {
		return "", err
	}

	return app.router.URI(app.tree, routeRequest)
}

func (app *application) RenderPage(request RouteRequest) (*RenderedPage, error) {
	page, err := app.renderRoute(request)
	if err != nil {
		return nil, err
	}

	return app.toRenderedPage(page), nil
}

func (app *application) RenderStyle(uri string) (*RenderedPage, error) {
	if app.tree == nil {
		return nil, errors.New("application tree is not initialized")
	}

	page, err := app.treeRenderer.Style(app.tree, uri)
	if err != nil {
		return nil, err
	}

	return app.toRenderedPage(page), nil
}

func (app *application) RenderJavascript(uri string) (*RenderedPage, error) {
	if app.tree == nil {
		return nil, errors.New("application tree is not initialized")
	}

	page, err := app.treeRenderer.Javascript(app.tree, uri)
	if err != nil {
		return nil, err
	}

	return app.toRenderedPage(page), nil
}

func (app *application) toRenderedPage(page treepages.Page) *RenderedPage {
	headers := make([]RenderedHeader, 0, len(page.Headers()))

	for _, header := range page.Headers() {
		headers = append(headers, RenderedHeader{
			Name:  header.Name(),
			Value: header.Value(),
		})
	}

	return &RenderedPage{
		HttpCode: page.HttpCode(),
		Headers:  headers,
		Body:     page.Body(),
	}
}

func (app *application) Trigger(msg eventapps.IncomingMessage) (*eventapps.OutgoingMessage, error) {
	return app.eventApplication.Execute(msg)
}

func (app *application) renderRoute(request RouteRequest) (treepages.Page, error) {
	if app.tree == nil {
		return nil, errors.New("application tree is not initialized")
	}

	routeRequest, err := app.toTreeRouteRequest(request)
	if err != nil {
		return nil, err
	}

	return app.treeRenderer.Render(app.tree, routeRequest)
}

func (app *application) resolvePage(request RouteRequest) (renderablepages.Page, error) {
	if app.tree == nil {
		return nil, errors.New("application tree is not initialized")
	}

	routeRequest, err := app.toTreeRouteRequest(request)
	if err != nil {
		return nil, err
	}

	return app.findPage(app.tree, routeRequest)
}

func (app *application) toTreeRouteRequest(request RouteRequest) (trees.Request, error) {
	return app.treeRequestBuilder.
		Create().
		WithPath(request.Path).
		WithMethod(request.Method).
		WithLocale(request.Locale).
		WithTarget(request.Target).
		Now()
}

func (app *application) findPage(tree trees.Tree, request trees.Request) (renderablepages.Page, error) {
	for _, node := range tree.Nodes() {
		for _, target := range node.Targets() {
			if target.Keyname() != request.Target() {
				continue
			}

			for _, group := range target.Groups() {
				resource, ok := group.Resource(request.Locale())
				if !ok {
					continue
				}

				if app.routeMatches(resource.Route(), request.Path()) {
					return resource.Page(), nil
				}

				if resource.HasChildren() {
					page, err := app.findPage(resource.Children(), request)
					if err == nil {
						return page, nil
					}
				}
			}
		}
	}

	return nil, errors.New("page not found")
}

func (app *application) routeMatches(route trees.Route, path string) bool {
	pattern := strings.Trim(route.Pattern(), "/")
	candidate := strings.Trim(path, "/")

	patternParts := splitPath(pattern)
	candidateParts := splitPath(candidate)

	if len(patternParts) != len(candidateParts) {
		return false
	}

	for index, patternPart := range patternParts {
		if isRouteParam(patternPart) {
			continue
		}

		if patternPart != candidateParts[index] {
			return false
		}
	}

	return true
}

func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

func isRouteParam(value string) bool {
	return strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")
}

func (app *application) registerTreeEvents(tree Tree) error {
	for _, node := range tree.Nodes {
		for _, target := range node.Targets {
			for _, group := range target.Groups {
				for _, resource := range group.Resources {
					if err := app.registerPageEvents(resource.Page); err != nil {
						return err
					}

					if resource.Children != nil {
						if err := app.registerTreeEvents(*resource.Children); err != nil {
							return err
						}
					}
				}
			}
		}

		for _, fallback := range node.Fallbacks {
			if err := app.registerPageEvents(fallback.Page); err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *application) registerPageEvents(page Page) error {
	return app.registerComponentEvents(page.Body)
}

func (app *application) registerComponentEvents(component Component) error {
	for _, event := range component.Events {
		domainEvent, err := app.toEvent(event)
		if err != nil {
			return err
		}

		if err := app.eventRegistry.Listen(domainEvent); err != nil {
			return err
		}
	}

	for _, child := range component.Children {
		if err := app.registerComponentEvents(child); err != nil {
			return err
		}
	}

	return nil
}

func (app *application) toTree(src Tree) (trees.Tree, error) {
	builder := app.treeBuilder.
		Create().
		WithKeyname(src.Keyname)

	for _, node := range src.Nodes {
		domainNode, err := app.toNode(node)
		if err != nil {
			return nil, err
		}

		builder.AddNode(domainNode)
	}

	return builder.Now()
}

func (app *application) toNode(src Node) (trees.Node, error) {
	builder := app.nodeBuilder.Create()

	for _, target := range src.Targets {
		domainTarget, err := app.toTarget(target)
		if err != nil {
			return nil, err
		}

		builder.AddTarget(domainTarget)
	}

	for _, fallback := range src.Fallbacks {
		domainFallback, err := app.toFallback(fallback)
		if err != nil {
			return nil, err
		}

		builder.AddFallback(domainFallback)
	}

	return builder.Now()
}

func (app *application) toTarget(src Target) (trees.Target, error) {
	builder := app.targetTreeBuilder.
		Create().
		WithKeyname(src.Keyname)

	for _, group := range src.Groups {
		domainGroup, err := app.toGroup(group)
		if err != nil {
			return nil, err
		}

		builder.AddGroup(domainGroup)
	}

	return builder.Now()
}

func (app *application) toGroup(src Group) (trees.Group, error) {
	builder := app.groupBuilder.
		Create().
		WithKeyname(src.Keyname)

	for _, resource := range src.Resources {
		domainResource, err := app.toResource(resource)
		if err != nil {
			return nil, err
		}

		builder.AddResource(domainResource)
	}

	return builder.Now()
}

func (app *application) toFallback(src Fallback) (trees.Fallback, error) {
	page, err := app.toPage(src.Page)
	if err != nil {
		return nil, err
	}

	return app.fallbackBuilder.
		Create().
		WithHttpCode(src.HttpCode).
		WithPage(page).
		Now()
}

func (app *application) toResource(src Resource) (trees.Resource, error) {
	route, err := app.toRoute(src.Route)
	if err != nil {
		return nil, err
	}

	page, err := app.toPage(src.Page)
	if err != nil {
		return nil, err
	}

	builder := app.resourceBuilder.
		Create().
		WithLocale(src.Locale).
		WithRoute(route).
		WithPage(page)

	if src.Children != nil {
		children, err := app.toTree(*src.Children)
		if err != nil {
			return nil, err
		}

		builder.WithChildren(children)
	}

	return builder.Now()
}

func (app *application) toRoute(src Route) (trees.Route, error) {
	params, err := app.toParams(src.Params)
	if err != nil {
		return nil, err
	}

	return app.routeBuilder.
		Create().
		WithPattern(src.Pattern).
		WithParams(params).
		Now()
}

func (app *application) toParams(src []Param) ([]trees.Param, error) {
	params := make([]trees.Param, 0, len(src))

	for _, param := range src {
		domainParam, err := app.paramBuilder.
			Create().
			WithKeyname(param.Keyname).
			WithValue(param.Value).
			Now()

		if err != nil {
			return nil, err
		}

		params = append(params, domainParam)
	}

	return params, nil
}

func (app *application) toPage(src Page) (renderablepages.Page, error) {
	head, err := app.toHead(src.Head)
	if err != nil {
		return nil, err
	}

	body, err := app.toComponent(src.Body)
	if err != nil {
		return nil, err
	}

	template, err := app.toTemplate(src.Template)
	if err != nil {
		return nil, err
	}

	return app.pageBuilder.
		Create().
		WithLanguage(src.Language).
		WithKeyname(src.Keyname).
		WithTemplate(template).
		WithHead(head).
		WithBody(body).
		Now()
}

func (app *application) toHead(src Head) (heads.Head, error) {
	title, err := app.toTemplate(src.Title)
	if err != nil {
		return nil, err
	}

	description, err := app.toTemplate(src.Description)
	if err != nil {
		return nil, err
	}

	builder := app.headBuilder.
		Create().
		WithTitle(title).
		WithDescription(description).
		WithLanguage(src.Language)

	for _, meta := range src.Meta {
		domainMeta, err := app.toMeta(meta)
		if err != nil {
			return nil, err
		}

		builder.AddMeta(domainMeta)
	}

	for _, link := range src.Links {
		domainLink, err := app.toLink(link)
		if err != nil {
			return nil, err
		}

		builder.AddLink(domainLink)
	}

	if src.OpenGraph != nil {
		openGraph, err := app.toOpenGraph(*src.OpenGraph)
		if err != nil {
			return nil, err
		}

		builder.WithOpenGraph(openGraph)
	}

	if src.TwitterCard != nil {
		twitterCard, err := app.toTwitterCard(*src.TwitterCard)
		if err != nil {
			return nil, err
		}

		builder.WithTwitterCard(twitterCard)
	}

	return builder.Now()
}

func (app *application) toMeta(src Meta) (heads.Meta, error) {
	content, err := app.toTemplate(src.Content)
	if err != nil {
		return nil, err
	}

	builder := app.metaBuilder.
		Create().
		WithContent(content)

	if src.Name != "" {
		builder.WithName(src.Name)
	}

	if src.Property != "" {
		builder.WithProperty(src.Property)
	}

	return builder.Now()
}

func (app *application) toLink(src Link) (heads.Link, error) {
	href, err := app.toTemplate(src.Href)
	if err != nil {
		return nil, err
	}

	return app.linkBuilder.
		Create().
		WithRel(src.Rel).
		WithHref(href).
		WithType(src.Type).
		Now()
}

func (app *application) toOpenGraph(src OpenGraph) (heads.OpenGraph, error) {
	title, err := app.toTemplate(src.Title)
	if err != nil {
		return nil, err
	}

	description, err := app.toTemplate(src.Description)
	if err != nil {
		return nil, err
	}

	image, err := app.toTemplate(src.Image)
	if err != nil {
		return nil, err
	}

	url, err := app.toTemplate(src.URL)
	if err != nil {
		return nil, err
	}

	siteName, err := app.toTemplate(src.SiteName)
	if err != nil {
		return nil, err
	}

	return app.openGraphBuilder.
		Create().
		WithTitle(title).
		WithDescription(description).
		WithImage(image).
		WithURL(url).
		WithType(src.Type).
		WithSiteName(siteName).
		Now()
}

func (app *application) toTwitterCard(src TwitterCard) (heads.TwitterCard, error) {
	title, err := app.toTemplate(src.Title)
	if err != nil {
		return nil, err
	}

	description, err := app.toTemplate(src.Description)
	if err != nil {
		return nil, err
	}

	image, err := app.toTemplate(src.Image)
	if err != nil {
		return nil, err
	}

	return app.twitterCardBuilder.
		Create().
		WithCard(src.Card).
		WithTitle(title).
		WithDescription(description).
		WithImage(image).
		Now()
}

func (app *application) toComponent(src Component) (components.Component, error) {
	template, err := app.toTemplate(src.Template)
	if err != nil {
		return nil, err
	}

	builder := app.componentBuilder.
		Create().
		WithKeyname(src.Keyname).
		WithTemplate(template)

	if src.Style != nil {
		style, err := app.toTemplate(*src.Style)
		if err != nil {
			return nil, err
		}

		builder.WithStyle(style)
	}

	for _, event := range src.Events {
		domainEvent, err := app.toEvent(event)
		if err != nil {
			return nil, err
		}

		builder.AddEvent(domainEvent)
	}

	for _, child := range src.Children {
		domainChild, err := app.toComponent(child)
		if err != nil {
			return nil, err
		}

		builder.AddChild(domainChild)
	}

	return builder.Now()
}

func (app *application) toEvent(src Event) (domainevents.Event, error) {
	return app.eventBuilder.
		Create().
		WithKeyname(src.Keyname).
		WithEventName(src.Keyname).
		WithAction(app.toActionFn(src.Action)).
		Now()
}

func (app *application) toActionFn(action ActionFn) domainevents.ActionFn {
	return func(ctx domaincontexts.Context) (domainresults.Result, error) {
		result, err := action(Context{
			EventName: ctx.EventName(),
			Payload:   ctx.Payload(),
		})

		if err != nil {
			return nil, err
		}

		if result == nil {
			return app.resultBuilder.Create().Now()
		}

		return app.toResult(*result)
	}
}

func (app *application) toResult(src Result) (domainresults.Result, error) {
	builder := app.resultBuilder.Create()

	for _, operation := range src.Operations {
		domainOperation, err := app.toOperation(operation)
		if err != nil {
			return nil, err
		}

		builder.AddOperation(domainOperation)
	}

	return builder.Now()
}

func (app *application) toOperation(src Operation) (domainresults.Operation, error) {
	action, err := app.toAction(src.Action)
	if err != nil {
		return nil, err
	}

	return app.operationBuilder.
		Create().
		WithType(src.Type).
		WithAction(action).
		Now()
}

func (app *application) toAction(src Action) (domainresults.Action, error) {
	builder := app.actionBuilder.Create()

	if src.Target != nil {
		target, err := app.targetBuilder.
			Create().
			WithElement(src.Target.Element).
			WithHTML(src.Target.HTML).
			Now()

		if err != nil {
			return nil, err
		}

		builder.WithTarget(target)
	}

	if src.Navigate != nil {
		navigate, err := app.navigateBuilder.
			Create().
			WithURL(src.Navigate.URL).
			Now()

		if err != nil {
			return nil, err
		}

		builder.WithNavigate(navigate)
	}

	return builder.Now()
}

func (app *application) toTemplate(src Template) (templates.Template, error) {
	return app.templateBuilder.
		Create().
		WithKeyname(src.Keyname).
		WithCode(src.Code).
		Now()
}
