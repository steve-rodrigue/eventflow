package assets

import "testing"

func TestRendererRenderNilAssets(t *testing.T) {
	result := NewRenderer().Render(nil)

	if result != "" {
		t.Fatalf("expected empty result, got %q", result)
	}
}

func TestRendererRenderEmptyAssets(t *testing.T) {
	assets, err := NewBuilder().Create().Now()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewRenderer().Render(assets)

	if result != "" {
		t.Fatalf("expected empty result, got %q", result)
	}
}

func TestRendererRenderCSSAsset(t *testing.T) {
	css := mustAsset(t, AssetTypeCSS, "/assets/app.css")

	assets, err := NewBuilder().
		Create().
		AddAsset(css).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewRenderer().Render(assets)
	expected := `<link rel="stylesheet" href="/assets/app.css">`

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderJavaScriptAsset(t *testing.T) {
	js := mustAsset(t, AssetTypeJavaScript, "/assets/app.js")

	assets, err := NewBuilder().
		Create().
		AddAsset(js).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewRenderer().Render(assets)
	expected := `<script src="/assets/app.js" defer></script>`

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderMultipleAssets(t *testing.T) {
	css := mustAsset(t, AssetTypeCSS, "/assets/app.css")
	js := mustAsset(t, AssetTypeJavaScript, "/assets/app.js")

	assets, err := NewBuilder().
		Create().
		WithAssets([]Asset{css, js}).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewRenderer().Render(assets)
	expected := `<link rel="stylesheet" href="/assets/app.css">` + "\n" +
		`<script src="/assets/app.js" defer></script>`

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererEscapesAssetURLs(t *testing.T) {
	css := mustAsset(t, AssetTypeCSS, `/assets/app.css?name=<bad>&v=1`)

	assets, err := NewBuilder().
		Create().
		AddAsset(css).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewRenderer().Render(assets)
	expected := `<link rel="stylesheet" href="/assets/app.css?name=&lt;bad&gt;&amp;v=1">`

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestBuilderAddAssetIgnoresNil(t *testing.T) {
	assets, err := NewBuilder().
		Create().
		AddAsset(nil).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if assets.HasAssets() {
		t.Fatal("expected no assets")
	}
}

func TestBuilderWithAssets(t *testing.T) {
	css := mustAsset(t, AssetTypeCSS, "/assets/app.css")

	assets, err := NewBuilder().
		Create().
		WithAssets([]Asset{css}).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !assets.HasAssets() {
		t.Fatal("expected assets")
	}

	if len(assets.Assets()) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(assets.Assets()))
	}
}

func TestAssetsReturnsClone(t *testing.T) {
	css := mustAsset(t, AssetTypeCSS, "/assets/app.css")
	js := mustAsset(t, AssetTypeJavaScript, "/assets/app.js")

	assets, err := NewBuilder().
		Create().
		WithAssets([]Asset{css}).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cloned := assets.Assets()
	cloned[0] = js

	if assets.Assets()[0].URL() != "/assets/app.css" {
		t.Fatal("expected assets to be cloned")
	}
}

func TestAssetBuilderNowCSS(t *testing.T) {
	asset, err := NewAssetBuilder().
		Create().
		WithType(AssetTypeCSS).
		WithURL("/assets/app.css").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if asset.Type() != AssetTypeCSS {
		t.Fatalf("expected css type, got %s", asset.Type())
	}

	if asset.URL() != "/assets/app.css" {
		t.Fatalf("expected url /assets/app.css, got %s", asset.URL())
	}
}

func TestAssetBuilderNowJavaScript(t *testing.T) {
	asset, err := NewAssetBuilder().
		Create().
		WithType(AssetTypeJavaScript).
		WithURL("/assets/app.js").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if asset.Type() != AssetTypeJavaScript {
		t.Fatalf("expected javascript type, got %s", asset.Type())
	}

	if asset.URL() != "/assets/app.js" {
		t.Fatalf("expected url /assets/app.js, got %s", asset.URL())
	}
}

func TestAssetBuilderTrimsURL(t *testing.T) {
	asset, err := NewAssetBuilder().
		Create().
		WithType(AssetTypeCSS).
		WithURL("  /assets/app.css  ").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if asset.URL() != "/assets/app.css" {
		t.Fatalf("expected trimmed url, got %q", asset.URL())
	}
}

func TestAssetBuilderNowReturnsErrorWhenTypeIsMissing(t *testing.T) {
	_, err := NewAssetBuilder().
		Create().
		WithURL("/assets/app.css").
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "asset type is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssetBuilderNowReturnsErrorWhenTypeIsInvalid(t *testing.T) {
	_, err := NewAssetBuilder().
		Create().
		WithType(AssetType("image")).
		WithURL("/assets/app.png").
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "asset type is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssetBuilderNowReturnsErrorWhenURLIsMissing(t *testing.T) {
	_, err := NewAssetBuilder().
		Create().
		WithType(AssetTypeCSS).
		WithURL("   ").
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "asset url is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}
