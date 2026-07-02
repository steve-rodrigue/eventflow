package heads

import (
	"strings"
	"testing"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

func TestRendererRenderMinimalHead(t *testing.T) {
	head := mustHead(t)

	result := NewRenderer(templates.NewMustacheRenderer(), nil).
		Render(head, renderables.Params{
			"article_title": "Hello World",
			"description":   "My description",
		}, nil)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Hello World</title>`,
		`<meta name="description" content="My description">`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererEscapesTemplateValues(t *testing.T) {
	head := mustHead(t)

	result := NewRenderer(templates.NewMustacheRenderer(), nil).
		Render(head, renderables.Params{
			"article_title": `<script>alert("x")</script>`,
			"description":   `A & B`,
		}, nil)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | &lt;script&gt;alert(&#34;x&#34;)&lt;/script&gt;</title>`,
		`<meta name="description" content="A &amp; B">`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderMetaNameAndProperty(t *testing.T) {
	metaName := mustMetaName(t, "robots", "index,{follow}")
	metaProperty := mustMetaProperty(t, `custom"<property>`, "{description}")

	head := mustHeadWithOptions(t, []Meta{metaName, metaProperty}, nil, nil, nil)

	result := NewRenderer(templates.NewMustacheRenderer(), nil).
		Render(head, renderables.Params{
			"article_title": "Article",
			"description":   `A "quoted" description`,
			"follow":        "follow",
		}, nil)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Article</title>`,
		`<meta name="description" content="A &#34;quoted&#34; description">`,
		`<meta name="robots" content="index,follow">`,
		`<meta property="custom&#34;&lt;property&gt;" content="A &#34;quoted&#34; description">`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderLinksWithAndWithoutType(t *testing.T) {
	canonical := mustLink(t, "canonical", "https://example.com/{slug}", "")
	feed := mustLink(t, `alternate"<bad>`, "https://example.com/feed.xml", `application/rss+xml"<bad>`)

	head := mustHeadWithOptions(t, nil, []Link{canonical, feed}, nil, nil)

	result := NewRenderer(templates.NewMustacheRenderer(), nil).
		Render(head, renderables.Params{
			"article_title": "Article",
			"description":   "Description",
			"slug":          `hello?x=<bad>&y=1`,
		}, nil)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Article</title>`,
		`<meta name="description" content="Description">`,
		`<link rel="canonical" href="https://example.com/hello?x=&lt;bad&gt;&amp;y=1">`,
		`<link rel="alternate&#34;&lt;bad&gt;" href="https://example.com/feed.xml" type="application/rss+xml&#34;&lt;bad&gt;">`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderOpenGraph(t *testing.T) {
	openGraph := mustOpenGraph(t)

	head := mustHeadWithOptions(t, nil, nil, openGraph, nil)

	result := NewRenderer(templates.NewMustacheRenderer(), nil).
		Render(head, renderables.Params{
			"article_title": "Article",
			"description":   "Description",
			"image":         "https://example.com/image.png",
			"url":           "https://example.com/article",
			"site_name":     "SteveCare",
		}, nil)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Article</title>`,
		`<meta name="description" content="Description">`,
		`<meta property="og:title" content="Article">`,
		`<meta property="og:description" content="Description">`,
		`<meta property="og:image" content="https://example.com/image.png">`,
		`<meta property="og:url" content="https://example.com/article">`,
		`<meta property="og:type" content="article">`,
		`<meta property="og:site_name" content="SteveCare">`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderTwitterCard(t *testing.T) {
	twitterCard := mustTwitterCard(t)

	head := mustHeadWithOptions(t, nil, nil, nil, twitterCard)

	result := NewRenderer(templates.NewMustacheRenderer(), nil).
		Render(head, renderables.Params{
			"article_title": "Article",
			"description":   "Description",
			"image":         "https://example.com/image.png",
		}, nil)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Article</title>`,
		`<meta name="description" content="Description">`,
		`<meta name="twitter:card" content="summary_large_image">`,
		`<meta name="twitter:title" content="Article">`,
		`<meta name="twitter:description" content="Description">`,
		`<meta name="twitter:image" content="https://example.com/image.png">`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderAssets(t *testing.T) {
	css := mustAsset(t, assets.AssetTypeCSS, "/assets/app.css")
	js := mustAsset(t, assets.AssetTypeJavaScript, "/assets/app.js")

	pageAssets, err := assets.NewBuilder().
		Create().
		WithAssets([]assets.Asset{css, js}).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	head := mustHead(t)

	result := NewRenderer(templates.NewMustacheRenderer(), assets.NewRenderer()).
		Render(head, renderables.Params{
			"article_title": "Article",
			"description":   "Description",
		}, pageAssets)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Article</title>`,
		`<meta name="description" content="Description">`,
		`<link rel="stylesheet" href="/assets/app.css">`,
		`<script src="/assets/app.js" defer></script>`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererDoesNotRenderEmptyAssets(t *testing.T) {
	pageAssets, err := assets.NewBuilder().Create().Now()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	head := mustHead(t)

	result := NewRenderer(templates.NewMustacheRenderer(), assets.NewRenderer()).
		Render(head, renderables.Params{
			"article_title": "Article",
			"description":   "Description",
		}, pageAssets)

	expected := strings.Join([]string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		`<title>SteveCare | Article</title>`,
		`<meta name="description" content="Description">`,
	}, "\n")

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
