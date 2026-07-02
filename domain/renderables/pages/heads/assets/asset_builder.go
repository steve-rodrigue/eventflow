package assets

import (
	"errors"
	"strings"
)

type assetBuilder struct {
	assetType AssetType
	url       string
}

func (b *assetBuilder) Create() AssetBuilder {
	return &assetBuilder{}
}

func (b *assetBuilder) WithType(assetType AssetType) AssetBuilder {
	b.assetType = assetType
	return b
}

func (b *assetBuilder) WithURL(url string) AssetBuilder {
	b.url = strings.TrimSpace(url)
	return b
}

func (b *assetBuilder) Now() (Asset, error) {
	switch b.assetType {
	case AssetTypeCSS, AssetTypeJavaScript:
		// valid
	default:
		return nil, errors.New("asset type is required")
	}

	if b.url == "" {
		return nil, errors.New("asset url is required")
	}

	return &asset{
		assetType: b.assetType,
		url:       b.url,
	}, nil
}
