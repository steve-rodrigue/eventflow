package pages

import (
	"strings"
	"testing"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

func TestRendererRenderWithRenderableParams(t *testing.T) {
	page := mustPage(t)
	pageAssets := mustAssets(t)

	result := NewRenderer(
		templates.NewMustacheRenderer(),
		heads.NewRenderer(templates.NewMustacheRenderer(), assets.NewRenderer()),
		components.NewRenderer(templates.NewMustacheRenderer()),
	).Render(page, renderables.Params{
		"head": renderables.Params{
			"title":       "Home",
			"description": "Welcome",
		},
		"body": renderables.Params{
			"message": "Hello",
		},
	}, pageAssets)

	expected := strings.Join([]string{
		`<!doctype html>`,
		`<html lang="en">`,
		`<head>`,
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Home</title>`,
		`<meta name="description" content="Welcome">`,
		`<link rel="stylesheet" href="/assets/home.css">`,
		`<script src="/assets/home.js" defer></script>`,
		`</head>`,
		`<body><main>Hello</main></body>`,
		`</html>`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderWithMapParams(t *testing.T) {
	page := mustPage(t)

	result := NewRenderer(
		templates.NewMustacheRenderer(),
		heads.NewRenderer(templates.NewMustacheRenderer(), nil),
		components.NewRenderer(templates.NewMustacheRenderer()),
	).Render(page, renderables.Params{
		"head": map[string]any{
			"title":       "Home",
			"description": "Welcome",
		},
		"body": map[string]any{
			"message": "Hello",
		},
	}, nil)

	expected := strings.Join([]string{
		`<!doctype html>`,
		`<html lang="en">`,
		`<head>`,
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Home</title>`,
		`<meta name="description" content="Welcome">`,
		`</head>`,
		`<body><main>Hello</main></body>`,
		`</html>`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderWithMissingParams(t *testing.T) {
	page := mustPage(t)

	result := NewRenderer(
		templates.NewMustacheRenderer(),
		heads.NewRenderer(templates.NewMustacheRenderer(), nil),
		components.NewRenderer(templates.NewMustacheRenderer()),
	).Render(page, renderables.Params{}, nil)

	expected := strings.Join([]string{
		`<!doctype html>`,
		`<html lang="en">`,
		`<head>`,
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | </title>`,
		`<meta name="description" content="">`,
		`</head>`,
		`<body><main></main></body>`,
		`</html>`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderWithScalarParams(t *testing.T) {
	page := mustPage(t)

	result := NewRenderer(
		templates.NewMustacheRenderer(),
		heads.NewRenderer(templates.NewMustacheRenderer(), nil),
		components.NewRenderer(templates.NewMustacheRenderer()),
	).Render(page, renderables.Params{
		"head": "Home",
		"body": "Hello",
	}, nil)

	expected := strings.Join([]string{
		`<!doctype html>`,
		`<html lang="en">`,
		`<head>`,
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | </title>`,
		`<meta name="description" content="">`,
		`</head>`,
		`<body><main></main></body>`,
		`</html>`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderStyle(t *testing.T) {
	page := mustStyledPage(t)

	result := NewRenderer(
		templates.NewMustacheRenderer(),
		heads.NewRenderer(templates.NewMustacheRenderer(), nil),
		components.NewRenderer(templates.NewMustacheRenderer()),
	).RenderStyle(page, renderables.Params{
		"color": "red",
	})

	expected := `.body { color: red; }`

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
