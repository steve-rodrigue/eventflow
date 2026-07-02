package assets

type asset struct {
	assetType AssetType
	url       string
}

func (a *asset) Type() AssetType {
	return a.assetType
}

func (a *asset) URL() string {
	return a.url
}
