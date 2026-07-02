package heads

import (
	"fmt"
	"html"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type renderer struct {
	templateRenderer templates.Renderer
	assetsRenderer   assets.Renderer
}

func (r *renderer) Render(head Head, params renderables.Params, assets assets.Assets) string {
	values := toTemplateValues(params)

	parts := []string{
		`<meta charset="utf-8">`,
		`<meta name="viewport" content="width=device-width, initial-scale=1">`,
		fmt.Sprintf(`<title>%s</title>`, r.render(head.Title(), values)),
		fmt.Sprintf(`<meta name="description" content="%s">`, r.render(head.Description(), values)),
	}

	for _, meta := range head.Meta() {
		content := r.render(meta.Content(), values)

		if meta.HasName() {
			parts = append(parts, fmt.Sprintf(
				`<meta name="%s" content="%s">`,
				html.EscapeString(meta.Name()),
				content,
			))
		}

		if meta.HasProperty() {
			parts = append(parts, fmt.Sprintf(
				`<meta property="%s" content="%s">`,
				html.EscapeString(meta.Property()),
				content,
			))
		}
	}

	for _, link := range head.Links() {
		href := r.render(link.Href(), values)

		if link.HasType() {
			parts = append(parts, fmt.Sprintf(
				`<link rel="%s" href="%s" type="%s">`,
				html.EscapeString(link.Rel()),
				href,
				html.EscapeString(link.Type()),
			))
			continue
		}

		parts = append(parts, fmt.Sprintf(
			`<link rel="%s" href="%s">`,
			html.EscapeString(link.Rel()),
			href,
		))
	}

	if head.HasOpenGraph() {
		openGraph := head.OpenGraph()

		parts = append(parts,
			fmt.Sprintf(`<meta property="og:title" content="%s">`, r.render(openGraph.Title(), values)),
			fmt.Sprintf(`<meta property="og:description" content="%s">`, r.render(openGraph.Description(), values)),
			fmt.Sprintf(`<meta property="og:image" content="%s">`, r.render(openGraph.Image(), values)),
			fmt.Sprintf(`<meta property="og:url" content="%s">`, r.render(openGraph.URL(), values)),
			fmt.Sprintf(`<meta property="og:type" content="%s">`, html.EscapeString(openGraph.Type())),
			fmt.Sprintf(`<meta property="og:site_name" content="%s">`, r.render(openGraph.SiteName(), values)),
		)
	}

	if head.HasTwitterCard() {
		twitterCard := head.TwitterCard()

		parts = append(parts,
			fmt.Sprintf(`<meta name="twitter:card" content="%s">`, html.EscapeString(twitterCard.Card())),
			fmt.Sprintf(`<meta name="twitter:title" content="%s">`, r.render(twitterCard.Title(), values)),
			fmt.Sprintf(`<meta name="twitter:description" content="%s">`, r.render(twitterCard.Description(), values)),
			fmt.Sprintf(`<meta name="twitter:image" content="%s">`, r.render(twitterCard.Image(), values)),
		)
	}

	if r.assetsRenderer != nil {
		renderedAssets := r.assetsRenderer.Render(assets)
		if renderedAssets != "" {
			parts = append(parts, renderedAssets)
		}
	}

	return strings.Join(parts, "\n")
}

func (r *renderer) render(template templates.Template, values map[string]string) string {
	return html.EscapeString(r.templateRenderer.Render(template, values))
}

func toTemplateValues(params renderables.Params) map[string]string {
	values := make(map[string]string, len(params))

	for key, value := range params {
		values[key] = fmt.Sprint(value)
	}

	return values
}
