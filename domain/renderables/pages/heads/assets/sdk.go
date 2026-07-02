package assets

// AssetType represents the type of asset.
type AssetType string

const (
	AssetTypeCSS        AssetType = "css"
	AssetTypeJavaScript AssetType = "javascript"
)

// NewRenderer creates a new assets renderer.
func NewRenderer() Renderer {
	return &renderer{}
}

// NewBuilder creates a new assets builder.
func NewBuilder() Builder {
	return &builder{}
}

// NewAssetBuilder creates a new asset builder.
func NewAssetBuilder() AssetBuilder {
	return &assetBuilder{}
}

// Renderer represents an assets renderer.
type Renderer interface {
	Render(assets Assets) string
}

// Builder represents an assets builder.
type Builder interface {
	Create() Builder

	AddAsset(asset Asset) Builder
	WithAssets(assets []Asset) Builder

	Now() (Assets, error)
}

// Assets represents the generated assets for a page.
type Assets interface {
	HasAssets() bool
	Assets() []Asset
}

// AssetBuilder represents an asset builder.
type AssetBuilder interface {
	Create() AssetBuilder
	WithType(assetType AssetType) AssetBuilder
	WithURL(url string) AssetBuilder
	Now() (Asset, error)
}

// Asset represents a generated asset.
type Asset interface {
	Type() AssetType
	URL() string
}
