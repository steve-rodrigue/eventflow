package assets

import "testing"

func mustAsset(t *testing.T, assetType AssetType, url string) Asset {
	t.Helper()

	asset, err := NewAssetBuilder().
		Create().
		WithType(assetType).
		WithURL(url).
		Now()

	if err != nil {
		t.Fatalf("expected asset, got error: %v", err)
	}

	return asset
}
