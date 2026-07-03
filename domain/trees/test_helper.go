package trees

import (
	"testing"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	renderable_pages "github.com/steve-rodrigue/eventflow/domain/renderables/pages"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
	rendered_pages "github.com/steve-rodrigue/eventflow/domain/trees/pages"
)

func newRealRenderer() Renderer {
	return newRealRendererWithAssetsBasePath("/assets")
}

func newRealRendererWithAssetsBasePath(assetsBasePath string) Renderer {
	templateRenderer := templates.NewMustacheRenderer()
	return NewRenderer(
		renderable_pages.NewRenderer(
			templateRenderer,
			heads.NewRenderer(templateRenderer, assets.NewRenderer()),
			components.NewRenderer(templateRenderer),
		),
		rendered_pages.NewBuilder(),
		rendered_pages.NewHeaderBuilder(),
		assets.NewBuilder(),
		assets.NewAssetBuilder(),
		assetsBasePath,
	)

}

func mustRenderablePage(t *testing.T, keyname string) renderable_pages.Page {
	t.Helper()

	page, err := renderable_pages.NewBuilder(renderables.NewBuilder()).
		Create().
		WithLanguage("en").
		WithKeyname(keyname).
		WithTemplate(mustTemplate(t, "page_"+keyname, `<html lang="{language}"><head>{head}</head><body>{body}</body></html>`)).
		WithHead(mustHead(t)).
		WithBody(mustBody(t)).
		Now()

	if err != nil {
		t.Fatalf("expected page, got error: %v", err)
	}

	return page
}

func mustHead(t *testing.T) heads.Head {
	t.Helper()

	head, err := heads.NewBuilder().
		Create().
		WithTitle(mustTemplate(t, "title", "Title")).
		WithDescription(mustTemplate(t, "description", "Description")).
		WithLanguage("en").
		Now()

	if err != nil {
		t.Fatalf("expected head, got error: %v", err)
	}

	return head
}

func mustBody(t *testing.T) components.Component {
	t.Helper()

	body, err := components.NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("body").
		WithTemplate(mustTemplate(t, "body", "<main>{message}</main>")).
		Now()

	if err != nil {
		t.Fatalf("expected body, got error: %v", err)
	}

	return body
}

func mustTemplate(t *testing.T, keyname string, code string) templates.Template {
	t.Helper()

	template, err := templates.NewMustacheBuilder().
		Create().
		WithKeyname(keyname).
		WithCode(code).
		Now()

	if err != nil {
		t.Fatalf("expected template, got error: %v", err)
	}

	return template
}

func mustRequest(t *testing.T, path string, method string, locale string, target string) Request {
	t.Helper()

	request, err := NewRequestBuilder().
		Create().
		WithPath(path).
		WithMethod(method).
		WithLocale(locale).
		WithTarget(target).
		Now()

	if err != nil {
		t.Fatalf("expected request, got error: %v", err)
	}

	return request
}

func mustTree(t *testing.T, nodes ...Node) Tree {
	t.Helper()

	tree, err := NewBuilder().
		Create().
		WithKeyname("main").
		WithNodes(nodes).
		Now()

	if err != nil {
		t.Fatalf("expected tree, got error: %v", err)
	}

	return tree
}

func mustNode(t *testing.T, targets ...Target) Node {
	t.Helper()

	node, err := NewNodeBuilder().
		Create().
		WithTargets(targets).
		Now()

	if err != nil {
		t.Fatalf("expected node, got error: %v", err)
	}

	return node
}

func mustNodeWithFallback(t *testing.T, target Target, fallback Fallback) Node {
	t.Helper()

	node, err := NewNodeBuilder().
		Create().
		AddTarget(target).
		AddFallback(fallback).
		Now()

	if err != nil {
		t.Fatalf("expected node, got error: %v", err)
	}

	return node
}

func mustTarget(t *testing.T, keyname string, groups ...Group) Target {
	t.Helper()

	target, err := NewTargetBuilder().
		Create().
		WithKeyname(keyname).
		WithGroups(groups).
		Now()

	if err != nil {
		t.Fatalf("expected target, got error: %v", err)
	}

	return target
}

func mustGroup(t *testing.T, resources ...Resource) Group {
	t.Helper()

	group, err := NewGroupBuilder().
		Create().
		WithKeyname("default").
		WithResources(resources).
		Now()

	if err != nil {
		t.Fatalf("expected group, got error: %v", err)
	}

	return group
}

func mustResource(t *testing.T, locale string, pattern string, page renderable_pages.Page) Resource {
	t.Helper()

	resource, err := NewResourceBuilder().
		Create().
		WithLocale(locale).
		WithRoute(mustRoute(t, pattern)).
		WithPage(page).
		Now()

	if err != nil {
		t.Fatalf("expected resource, got error: %v", err)
	}

	return resource
}

func mustResourceWithChildren(t *testing.T, locale string, pattern string, page renderable_pages.Page, children Tree) Resource {
	t.Helper()

	resource, err := NewResourceBuilder().
		Create().
		WithLocale(locale).
		WithRoute(mustRoute(t, pattern)).
		WithPage(page).
		WithChildren(children).
		Now()

	if err != nil {
		t.Fatalf("expected resource, got error: %v", err)
	}

	return resource
}

func mustRoute(t *testing.T, pattern string) Route {
	t.Helper()

	route, err := NewRouteBuilder().
		Create().
		WithPattern(pattern).
		Now()

	if err != nil {
		t.Fatalf("expected route, got error: %v", err)
	}

	return route
}

func mustFallback(t *testing.T, httpCode int, page renderable_pages.Page) Fallback {
	t.Helper()

	fallback, err := NewFallbackBuilder().
		Create().
		WithHttpCode(httpCode).
		WithPage(page).
		Now()

	if err != nil {
		t.Fatalf("expected fallback, got error: %v", err)
	}

	return fallback
}

func mustStyledRenderablePage(t *testing.T, keyname string) renderable_pages.Page {
	t.Helper()

	style := mustTemplate(t, "body_style", "body { color: red; }")

	body, err := components.NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("body").
		WithTemplate(mustTemplate(t, "body", "<main></main>")).
		WithStyle(style).
		Now()
	if err != nil {
		t.Fatalf("expected body, got error: %v", err)
	}

	page, err := renderable_pages.NewBuilder(renderables.NewBuilder()).
		Create().
		WithLanguage("en").
		WithKeyname(keyname).
		WithTemplate(mustTemplate(t, "page_"+keyname, `<html lang="{language}"><head>{head}</head><body>{body}</body></html>`)).
		WithHead(mustHead(t)).
		WithBody(body).
		Now()
	if err != nil {
		t.Fatalf("expected page, got error: %v", err)
	}

	return page
}

func assertRenderedPage(t *testing.T, page rendered_pages.Page, httpCode int, bodyContains string) {
	t.Helper()

	if page.HttpCode() != httpCode {
		t.Fatalf("expected http code %d, got %d", httpCode, page.HttpCode())
	}

	if page.Body() == "" {
		t.Fatal("expected body")
	}

	if !contains(page.Body(), bodyContains) {
		t.Fatalf("expected body to contain %q, got %q", bodyContains, page.Body())
	}

	if len(page.Headers()) != 1 {
		t.Fatalf("expected 1 header, got %d", len(page.Headers()))
	}

	if page.Headers()[0].Name() != "Content-Type" {
		t.Fatalf("expected Content-Type header, got %s", page.Headers()[0].Name())
	}

	if page.Headers()[0].Value() != "text/html; charset=utf-8" {
		t.Fatalf("unexpected header value: %s", page.Headers()[0].Value())
	}
}

func assertHeader(t *testing.T, page rendered_pages.Page, name string, value string) {
	t.Helper()

	for _, header := range page.Headers() {
		if header.Name() == name && header.Value() == value {
			return
		}
	}

	t.Fatalf("expected header %s=%s, got %#v", name, value, page.Headers())
}

func contains(value string, expected string) bool {
	return len(expected) == 0 || len(value) >= len(expected) && (value == expected || contains(value[1:], expected) || value[:len(expected)] == expected)
}
