package routers

import (
	"testing"
)

func TestRouterURIReturnsErrorWhenTreeIsNil(t *testing.T) {
	_, err := NewRouter().URI(nil, mustRequest(t, "desktop", "default", "en"))

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "tree is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRouterURIReturnsErrorWhenRequestIsNil(t *testing.T) {
	_, err := NewRouter().URI(mustTree(t), nil)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "route request is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRouterURIReturnsStaticURI(t *testing.T) {
	uri, err := NewRouter().URI(
		mustTree(t,
			mustNode(t,
				mustTarget(t, "desktop",
					mustGroup(t, "default",
						mustResource(t, "en", "/articles"),
					),
				),
			),
		),
		mustRequest(t, "desktop", "default", "en"),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if uri != "/articles" {
		t.Fatalf("expected /articles, got %s", uri)
	}
}

func TestRouterURIRendersParams(t *testing.T) {
	uri, err := NewRouter().URI(
		mustTree(t,
			mustNode(t,
				mustTarget(t, "desktop",
					mustGroup(t, "default",
						mustResource(t, "fr-CA", "/fr/articles/{slug}/{id}"),
					),
				),
			),
		),
		mustRequest(
			t,
			"desktop",
			"default",
			"fr-CA",
			mustParam(t, "slug", "bonjour"),
			mustParam(t, "id", "123"),
		),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if uri != "/fr/articles/bonjour/123" {
		t.Fatalf("expected rendered URI, got %s", uri)
	}
}

func TestRouterURIFindsResourceInChildTree(t *testing.T) {
	childTree := mustTree(t,
		mustNode(t,
			mustTarget(t, "desktop",
				mustGroup(t, "default",
					mustResource(t, "en", "/parent/child/{slug}"),
				),
			),
		),
	)

	parent := mustResourceWithChildren(t, "en", "/parent", childTree)

	uri, err := NewRouter().URI(
		mustTree(t,
			mustNode(t,
				mustTarget(t, "desktop",
					mustGroup(t, "default", parent),
				),
			),
		),
		mustRequest(t, "desktop", "default", "en", mustParam(t, "slug", "hello")),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if uri != "/parent/child/hello" {
		t.Fatalf("expected child URI, got %s", uri)
	}
}

func TestRouterURIReturnsErrorWhenTargetDoesNotMatch(t *testing.T) {
	_, err := NewRouter().URI(
		mustTree(t,
			mustNode(t,
				mustTarget(t, "mobile",
					mustGroup(t, "default",
						mustResource(t, "en", "/articles"),
					),
				),
			),
		),
		mustRequest(t, "desktop", "default", "en"),
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "resource not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRouterURIReturnsErrorWhenGroupDoesNotMatch(t *testing.T) {
	_, err := NewRouter().URI(
		mustTree(t,
			mustNode(t,
				mustTarget(t, "desktop",
					mustGroup(t, "premium",
						mustResource(t, "en", "/articles"),
					),
				),
			),
		),
		mustRequest(t, "desktop", "default", "en"),
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "resource not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRouterURIReturnsErrorWhenLocaleDoesNotMatch(t *testing.T) {
	_, err := NewRouter().URI(
		mustTree(t,
			mustNode(t,
				mustTarget(t, "desktop",
					mustGroup(t, "default",
						mustResource(t, "fr", "/articles"),
					),
				),
			),
		),
		mustRequest(t, "desktop", "default", "en"),
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "resource not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}
