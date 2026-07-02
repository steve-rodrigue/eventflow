package assets

type builder struct {
	assets []Asset
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) AddAsset(asset Asset) Builder {
	if asset != nil {
		b.assets = append(b.assets, asset)
	}

	return b
}

func (b *builder) WithAssets(assets []Asset) Builder {
	b.assets = assets
	return b
}

func (b *builder) Now() (Assets, error) {
	return &assets{
		assets: b.assets,
	}, nil
}
