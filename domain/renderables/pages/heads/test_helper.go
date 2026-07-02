package heads

import (
	"testing"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

func mustHead(t *testing.T) Head {
	t.Helper()

	return mustHeadWithOptions(t, nil, nil, nil, nil)
}

func mustHeadWithOptions(t *testing.T, meta []Meta, links []Link, openGraph OpenGraph, twitterCard TwitterCard) Head {
	t.Helper()

	head, err := NewBuilder().
		Create().
		WithTitle(mustTemplate(t, "title", "SteveCare | {article_title}")).
		WithDescription(mustTemplate(t, "description", "{description}")).
		WithLanguage("en").
		WithMeta(meta).
		WithLinks(links).
		WithOpenGraph(openGraph).
		WithTwitterCard(twitterCard).
		Now()

	if err != nil {
		t.Fatalf("expected head, got error: %v", err)
	}

	return head
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

func mustMetaName(t *testing.T, name string, content string) Meta {
	t.Helper()

	meta, err := NewMetaBuilder().
		Create().
		WithName(name).
		WithContent(mustTemplate(t, "meta_"+name, content)).
		Now()

	if err != nil {
		t.Fatalf("expected meta, got error: %v", err)
	}

	return meta
}

func mustMetaProperty(t *testing.T, property string, content string) Meta {
	t.Helper()

	meta, err := NewMetaBuilder().
		Create().
		WithProperty(property).
		WithContent(mustTemplate(t, "meta_property", content)).
		Now()

	if err != nil {
		t.Fatalf("expected meta, got error: %v", err)
	}

	return meta
}

func mustLink(t *testing.T, rel string, href string, linkType string) Link {
	t.Helper()

	builder := NewLinkBuilder().
		Create().
		WithRel(rel).
		WithHref(mustTemplate(t, "link_"+rel, href))

	if linkType != "" {
		builder = builder.WithType(linkType)
	}

	link, err := builder.Now()
	if err != nil {
		t.Fatalf("expected link, got error: %v", err)
	}

	return link
}

func mustOpenGraph(t *testing.T) OpenGraph {
	t.Helper()

	openGraph, err := NewOpenGraphBuilder().
		Create().
		WithTitle(mustTemplate(t, "og_title", "{article_title}")).
		WithDescription(mustTemplate(t, "og_description", "{description}")).
		WithImage(mustTemplate(t, "og_image", "{image}")).
		WithURL(mustTemplate(t, "og_url", "{url}")).
		WithType("article").
		WithSiteName(mustTemplate(t, "og_site_name", "{site_name}")).
		Now()

	if err != nil {
		t.Fatalf("expected open graph, got error: %v", err)
	}

	return openGraph
}

func mustTwitterCard(t *testing.T) TwitterCard {
	t.Helper()

	twitterCard, err := NewTwitterCardBuilder().
		Create().
		WithCard("summary_large_image").
		WithTitle(mustTemplate(t, "twitter_title", "{article_title}")).
		WithDescription(mustTemplate(t, "twitter_description", "{description}")).
		WithImage(mustTemplate(t, "twitter_image", "{image}")).
		Now()

	if err != nil {
		t.Fatalf("expected twitter card, got error: %v", err)
	}

	return twitterCard
}

func mustAsset(t *testing.T, assetType assets.AssetType, url string) assets.Asset {
	t.Helper()

	asset, err := assets.NewAssetBuilder().
		Create().
		WithType(assetType).
		WithURL(url).
		Now()

	if err != nil {
		t.Fatalf("expected asset, got error: %v", err)
	}

	return asset
}
