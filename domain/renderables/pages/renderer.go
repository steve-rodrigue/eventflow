package pages

import (
	"fmt"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type renderer struct {
	templateRenderer  templates.Renderer
	headRenderer      heads.Renderer
	componentRenderer components.Renderer
}

func (r *renderer) Render(page Page, params renderables.Params, pageAssets assets.Assets) string {
	headParams := paramsForKey(params, "head")
	bodyParams := paramsForKey(params, "body")

	head := r.headRenderer.Render(page.Head(), headParams, pageAssets)
	body := r.componentRenderer.Render(page.Body(), bodyParams)

	return r.templateRenderer.Render(page.Template(), map[string]string{
		"language": page.Language(),
		"head":     head,
		"body":     body,
	})
}

func (r *renderer) RenderStyle(page Page, params renderables.Params) string {
	return r.componentRenderer.RenderStyle(page.Body(), params)
}

func paramsForKey(params renderables.Params, key string) renderables.Params {
	value, ok := params[key]
	if !ok {
		return renderables.Params{}
	}

	switch typed := value.(type) {
	case renderables.Params:
		return typed
	case map[string]any:
		return renderables.Params(typed)
	default:
		return renderables.Params{
			key: fmt.Sprint(typed),
		}
	}
}
