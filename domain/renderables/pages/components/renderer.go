package components

import (
	"fmt"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

type renderer struct {
	templateRenderer templates.Renderer
}

func (r *renderer) Render(component Component, params renderables.Params) string {
	values := r.toTemplateValues(params)

	for _, child := range component.Children() {
		childParams, ok := params[child.Keyname()]
		if !ok {
			continue
		}

		values[child.Keyname()] = r.renderChild(child, childParams)
	}

	return r.templateRenderer.Render(component.Template(), values)
}

func (r *renderer) RenderList(components []Component, params []renderables.Params) string {
	rendered := make([]string, 0, len(components))

	for index, component := range components {
		itemParams := renderables.Params{}

		if index < len(params) {
			itemParams = params[index]
		}

		rendered = append(rendered, r.Render(component, itemParams))
	}

	return strings.Join(rendered, "")
}

func (r *renderer) RenderStyle(component Component, params renderables.Params) string {
	return r.renderStyle(component, params, "")
}

func (r *renderer) RenderStyleList(components []Component, params []renderables.Params) string {
	rendered := make([]string, 0, len(components))

	for index, component := range components {
		itemParams := renderables.Params{}

		if index < len(params) {
			itemParams = params[index]
		}

		style := r.RenderStyle(component, itemParams)
		if style != "" {
			rendered = append(rendered, style)
		}
	}

	return strings.Join(rendered, "\n")
}

func (r *renderer) renderStyle(component Component, params renderables.Params, parentScope string) string {
	parts := []string{}

	currentScope := r.selector(parentScope, component.Keyname())

	if component.HasStyle() {
		values := r.toTemplateValues(params)
		style := r.templateRenderer.Render(component.Style(), values)
		parts = append(parts, r.scopeCSS(style, parentScope))
	}

	for _, child := range component.Children() {
		childParams, ok := params[child.Keyname()]
		if !ok {
			continue
		}

		style := r.renderChildStyle(child, childParams, currentScope)
		if style != "" {
			parts = append(parts, style)
		}
	}

	return strings.Join(parts, "\n")
}

func (r *renderer) renderChild(component Component, value any) string {
	switch typed := value.(type) {
	case renderables.Params:
		return r.Render(component, typed)

	case map[string]any:
		return r.Render(component, renderables.Params(typed))

	case []renderables.Params:
		rendered := make([]string, 0, len(typed))

		for _, item := range typed {
			rendered = append(rendered, r.Render(component, item))
		}

		return strings.Join(rendered, "")

	case []map[string]any:
		rendered := make([]string, 0, len(typed))

		for _, item := range typed {
			rendered = append(rendered, r.Render(component, renderables.Params(item)))
		}

		return strings.Join(rendered, "")

	default:
		return fmt.Sprint(typed)
	}
}

func (r *renderer) renderChildStyle(component Component, value any, parentScope string) string {
	switch typed := value.(type) {
	case renderables.Params:
		return r.renderStyle(component, typed, parentScope)

	case map[string]any:
		return r.renderStyle(component, renderables.Params(typed), parentScope)

	case []renderables.Params:
		rendered := make([]string, 0, len(typed))

		for _, item := range typed {
			style := r.renderStyle(component, item, parentScope)
			if style != "" {
				rendered = append(rendered, style)
			}
		}

		return strings.Join(rendered, "\n")

	case []map[string]any:
		rendered := make([]string, 0, len(typed))

		for _, item := range typed {
			style := r.renderStyle(component, renderables.Params(item), parentScope)
			if style != "" {
				rendered = append(rendered, style)
			}
		}

		return strings.Join(rendered, "\n")

	default:
		return fmt.Sprint(typed)
	}
}

func (r *renderer) scopeCSS(css string, parentScope string) string {
	parentScope = strings.TrimSpace(parentScope)

	if parentScope == "" {
		return css
	}

	lines := strings.Split(css, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" ||
			strings.HasPrefix(trimmed, "@") ||
			strings.HasPrefix(trimmed, "}") {
			continue
		}

		if strings.Contains(trimmed, "{") {
			lines[i] = strings.Replace(line, trimmed, parentScope+" "+trimmed, 1)
		}
	}

	return strings.Join(lines, "\n")
}

func (r *renderer) selector(parentScope string, keyname string) string {
	current := "." + keyname

	if parentScope == "" {
		return current
	}

	return parentScope + " " + current
}

func (r *renderer) toTemplateValues(params renderables.Params) map[string]string {
	values := make(map[string]string, len(params))

	for key, value := range params {
		switch value.(type) {
		case renderables.Params, map[string]any, []renderables.Params, []map[string]any:
			continue
		default:
			values[key] = fmt.Sprint(value)
		}
	}

	return values
}
