package assets

import (
	"fmt"
	"html"
	"strings"
)

type renderer struct{}

func (r *renderer) Render(assets Assets) string {
	if assets == nil || !assets.HasAssets() {
		return ""
	}

	output := make([]string, 0, len(assets.Assets()))

	for _, asset := range assets.Assets() {
		switch asset.Type() {
		case AssetTypeCSS:
			output = append(output, fmt.Sprintf(
				`<link rel="stylesheet" href="%s">`,
				html.EscapeString(asset.URL()),
			))

		case AssetTypeJavaScript:
			output = append(output, fmt.Sprintf(
				`<script src="%s" defer></script>`,
				html.EscapeString(asset.URL()),
			))
		}
	}

	return strings.Join(output, "\n")
}
