package trees

import (
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
