package trees

import (
	"strings"
	"testing"
)

func TestRendererRenderStaticRoute(t *testing.T) {
	page := mustRenderablePage(t, "home")

	rendered, err := newRealRenderer().Render(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/home", page))))),
		mustRequest(t, "/home", "GET", "en", "desktop"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertRenderedPage(t, rendered, 200, `<body><main></main></body>`)
}

func TestRendererRenderParamRoute(t *testing.T) {
	page := mustRenderablePage(t, "article")

	rendered, err := newRealRenderer().Render(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/articles/{slug}", page))))),
		mustRequest(t, "/articles/hello", "GET", "en", "desktop"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertRenderedPage(t, rendered, 200, `<body><main></main></body>`)
}

func TestRendererRenderChildTree(t *testing.T) {
	childPage := mustRenderablePage(t, "child")
	childTree := mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/parent/child", childPage)))))

	parentResource := mustResourceWithChildren(t, "en", "/parent", mustRenderablePage(t, "parent"), childTree)

	rendered, err := newRealRenderer().Render(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, parentResource)))),
		mustRequest(t, "/parent/child", "GET", "en", "desktop"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertRenderedPage(t, rendered, 200, `<body><main></main></body>`)
}

func TestRendererRenderFallbackWhenRouteDoesNotMatch(t *testing.T) {
	fallbackPage := mustRenderablePage(t, "not_found")

	rendered, err := newRealRenderer().Render(
		mustTree(t,
			mustNodeWithFallback(
				t,
				mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/home", mustRenderablePage(t, "home")))),
				mustFallback(t, 404, fallbackPage),
			),
		),
		mustRequest(t, "/missing", "GET", "en", "desktop"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertRenderedPage(t, rendered, 200, `<body><main></main></body>`)
}

func TestRendererRenderReturnsErrorWhenTreeIsNil(t *testing.T) {
	_, err := newRealRenderer().Render(nil, mustRequest(t, "/home", "GET", "en", "desktop"))

	if err == nil || err.Error() != "tree is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererRenderReturnsErrorWhenRequestIsNil(t *testing.T) {
	_, err := newRealRenderer().Render(mustTree(t), nil)

	if err == nil || err.Error() != "request is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererRenderReturnsPageNotFound(t *testing.T) {
	_, err := newRealRenderer().Render(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/home", mustRenderablePage(t, "home")))))),
		mustRequest(t, "/missing", "GET", "en", "desktop"),
	)

	if err == nil || err.Error() != "page not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererRenderSkipsDifferentTarget(t *testing.T) {
	_, err := newRealRenderer().Render(
		mustTree(t, mustNode(t, mustTarget(t, "mobile", mustGroup(t, mustResource(t, "en", "/home", mustRenderablePage(t, "home")))))),
		mustRequest(t, "/home", "GET", "en", "desktop"),
	)

	if err == nil || err.Error() != "page not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererRenderSkipsDifferentLocale(t *testing.T) {
	_, err := newRealRenderer().Render(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "fr", "/home", mustRenderablePage(t, "home")))))),
		mustRequest(t, "/home", "GET", "en", "desktop"),
	)

	if err == nil || err.Error() != "page not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererRenderPassesAssetsToPageRenderer(t *testing.T) {
	page := mustRenderablePage(t, "home")

	rendered, err := newRealRenderer().Render(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/", page))))),
		mustRequest(t, "/", "GET", "en", "desktop"),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(rendered.Body(), `<link rel="stylesheet" href="/assets/home.css">`) {
		t.Fatalf("expected css asset in body, got %q", rendered.Body())
	}

	if !strings.Contains(rendered.Body(), `<script src="/assets/home.js" defer></script>`) {
		t.Fatalf("expected js asset in body, got %q", rendered.Body())
	}
}

func TestRendererRenderUsesCustomAssetsBasePath(t *testing.T) {
	page := mustRenderablePage(t, "home")

	rendered, err := newRealRendererWithAssetsBasePath("/static/generated").Render(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/", page))))),
		mustRequest(t, "/", "GET", "en", "desktop"),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(rendered.Body(), `<link rel="stylesheet" href="/static/generated/home.css">`) {
		t.Fatalf("expected custom css asset path, got %q", rendered.Body())
	}

	if !strings.Contains(rendered.Body(), `<script src="/static/generated/home.js" defer></script>`) {
		t.Fatalf("expected custom js asset path, got %q", rendered.Body())
	}
}

func TestRendererStyleRendersCSSForAssetURI(t *testing.T) {
	page := mustStyledRenderablePage(t, "home")

	rendered, err := newRealRenderer().Style(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/", page))))),
		"/assets/home.css",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rendered.HttpCode() != 200 {
		t.Fatalf("expected http code 200, got %d", rendered.HttpCode())
	}

	if !strings.Contains(rendered.Body(), "body") {
		t.Fatalf("expected CSS body selector, got %q", rendered.Body())
	}

	assertHeader(t, rendered, "Content-Type", "text/css; charset=utf-8")
}

func TestRendererStyleUsesCustomAssetsBasePath(t *testing.T) {
	page := mustRenderablePage(t, "home")

	rendered, err := newRealRendererWithAssetsBasePath("/static/generated").Style(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/", page))))),
		"/static/generated/home.css",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertHeader(t, rendered, "Content-Type", "text/css; charset=utf-8")
}

func TestRendererJavascriptRendersRuntimeForAssetURI(t *testing.T) {
	page := mustRenderablePage(t, "home")

	rendered, err := newRealRenderer().Javascript(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/", page))))),
		"/assets/home.js",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rendered.HttpCode() != 200 {
		t.Fatalf("expected http code 200, got %d", rendered.HttpCode())
	}

	if !strings.Contains(rendered.Body(), `const socket = new WebSocket`) {
		t.Fatalf("expected runtime JS, got %q", rendered.Body())
	}

	assertHeader(t, rendered, "Content-Type", "application/javascript; charset=utf-8")
}

func TestRendererJavascriptUsesCustomAssetsBasePath(t *testing.T) {
	page := mustRenderablePage(t, "home")

	rendered, err := newRealRendererWithAssetsBasePath("/static/generated").Javascript(
		mustTree(t, mustNode(t, mustTarget(t, "desktop", mustGroup(t, mustResource(t, "en", "/", page))))),
		"/static/generated/home.js",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertHeader(t, rendered, "Content-Type", "application/javascript; charset=utf-8")
}

func TestRendererStyleReturnsErrorWhenTreeIsNil(t *testing.T) {
	_, err := newRealRenderer().Style(nil, "/assets/home.css")

	if err == nil || err.Error() != "tree is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererStyleReturnsPageNotFound(t *testing.T) {
	_, err := newRealRenderer().Style(
		mustTree(t),
		"/assets/home.css",
	)

	if err == nil || err.Error() != "page not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererJavascriptReturnsErrorWhenTreeIsNil(t *testing.T) {
	_, err := newRealRenderer().Javascript(nil, "/assets/home.js")

	if err == nil || err.Error() != "tree is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRendererJavascriptReturnsPageNotFound(t *testing.T) {
	_, err := newRealRenderer().Javascript(
		mustTree(t),
		"/assets/home.js",
	)

	if err == nil || err.Error() != "page not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNormalizeAssetsBasePath(t *testing.T) {
	tests := map[string]string{
		"":                  "/assets",
		"/":                 "/assets",
		"assets":            "/assets",
		"/assets":           "/assets",
		"/assets/":          "/assets",
		"static/generated":  "/static/generated",
		"/static/generated": "/static/generated",
	}

	for input, expected := range tests {
		result := normalizeAssetsBasePath(input)

		if result != expected {
			t.Fatalf("expected %q for %q, got %q", expected, input, result)
		}
	}
}
