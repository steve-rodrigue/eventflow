package assets

import "slices"

type assets struct {
	assets []Asset
}

func (a *assets) HasAssets() bool {
	return len(a.assets) > 0
}

func (a *assets) Assets() []Asset {
	return slices.Clone(a.assets)
}
