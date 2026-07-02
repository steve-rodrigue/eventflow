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

func mustPage(t *testing.T) Page {
	t.Helper()

	pageTemplate := mustTemplate(t, "page", strings.Join([]string{
		`<!doctype html>`,
		`<html lang="{language}">`,
		`<head>`,
		`{head}`,
		`</head>`,
		`<body>{body}</body>`,
		`</html>`,
	}, "\n"))

	page, err := NewBuilder(renderables.NewBuilder()).
		Create().
		WithLanguage("en").
		WithKeyname("home").
		WithTemplate(pageTemplate).
		WithHead(mustHead(t)).
		WithBody(mustBody(t)).
		Now()

	if err != nil {
		t.Fatalf("expected page, got error: %v", err)
	}

	return page
}

func mustStyledPage(t *testing.T) Page {
	t.Helper()

	page, err := NewBuilder(renderables.NewBuilder()).
		Create().
		WithLanguage("en").
		WithKeyname("styled").
		WithTemplate(mustTemplate(t, "page", `{body}`)).
		WithHead(mustHead(t)).
		WithBody(mustStyledBody(t)).
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
		WithTitle(mustTemplate(t, "title", "SteveCare | {title}")).
		WithDescription(mustTemplate(t, "description", "{description}")).
		WithLanguage("en").
		Now()

	if err != nil {
		t.Fatalf("expected head, got error: %v", err)
	}

	return head
}

func mustBody(t *testing.T) components.Component {
	t.Helper()

	component, err := components.NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("body").
		WithTemplate(mustTemplate(t, "body", "<main>{message}</main>")).
		Now()

	if err != nil {
		t.Fatalf("expected body component, got error: %v", err)
	}

	return component
}

func mustStyledBody(t *testing.T) components.Component {
	t.Helper()

	component, err := components.NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("body").
		WithTemplate(mustTemplate(t, "body", "<main>{message}</main>")).
		WithStyle(mustTemplate(t, "style", ".body { color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected body component, got error: %v", err)
	}

	return component
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

func mustAssets(t *testing.T) assets.Assets {
	t.Helper()

	css, err := assets.NewAssetBuilder().
		Create().
		WithType(assets.AssetTypeCSS).
		WithURL("/assets/home.css").
		Now()
	if err != nil {
		t.Fatalf("expected css asset, got error: %v", err)
	}

	js, err := assets.NewAssetBuilder().
		Create().
		WithType(assets.AssetTypeJavaScript).
		WithURL("/assets/home.js").
		Now()
	if err != nil {
		t.Fatalf("expected js asset, got error: %v", err)
	}

	pageAssets, err := assets.NewBuilder().
		Create().
		WithAssets([]assets.Asset{css, js}).
		Now()
	if err != nil {
		t.Fatalf("expected assets, got error: %v", err)
	}

	return pageAssets
}
