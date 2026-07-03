package trees

import (
	"errors"
	"net/http"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/renderables"
	renderable_pages "github.com/steve-rodrigue/eventflow/domain/renderables/pages"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	rendered_pages "github.com/steve-rodrigue/eventflow/domain/trees/pages"
)

type renderer struct {
	pageRenderer  renderable_pages.Renderer
	pageBuilder   rendered_pages.Builder
	headerBuilder rendered_pages.HeaderBuilder
	assetsBuilder assets.Builder
	assetBuilder  assets.AssetBuilder
}

func (r *renderer) Render(tree Tree, request Request) (rendered_pages.Page, error) {
	if tree == nil {
		return nil, errors.New("tree is required")
	}

	if request == nil {
		return nil, errors.New("request is required")
	}

	page, ok := r.findPage(tree, request)
	if !ok {
		return nil, errors.New("page not found")
	}

	pageAssets, err := r.buildAssets(page)
	if err != nil {
		return nil, err
	}

	body := r.pageRenderer.Render(page, renderables.Params{}, pageAssets)

	contentType, err := r.headerBuilder.
		Create().
		WithName("Content-Type").
		WithValue("text/html; charset=utf-8").
		Now()

	if err != nil {
		return nil, err
	}

	return r.pageBuilder.
		Create().
		WithHttpCode(http.StatusOK).
		AddHeader(contentType).
		WithBody(body).
		Now()
}

func (r *renderer) findPage(tree Tree, request Request) (renderable_pages.Page, bool) {
	for _, node := range tree.Nodes() {
		for _, target := range node.Targets() {
			if target.Keyname() != request.Target() {
				continue
			}

			for _, group := range target.Groups() {
				resource, ok := group.Resource(request.Locale())
				if !ok {
					continue
				}

				if r.matches(resource.Route(), request.Path()) {
					return resource.Page(), true
				}

				if resource.HasChildren() {
					page, ok := r.findPage(resource.Children(), request)
					if ok {
						return page, true
					}
				}
			}
		}

		if fallback, ok := r.fallback(node, http.StatusNotFound); ok {
			return fallback.Page(), true
		}
	}

	return nil, false
}

func (r *renderer) matches(route Route, path string) bool {
	pattern := strings.Trim(route.Pattern(), "/")
	candidate := strings.Trim(path, "/")

	patternParts := splitPath(pattern)
	candidateParts := splitPath(candidate)

	if len(patternParts) != len(candidateParts) {
		return false
	}

	for i, patternPart := range patternParts {
		if isParam(patternPart) {
			continue
		}

		if patternPart != candidateParts[i] {
			return false
		}
	}

	return true
}

func (r *renderer) fallback(node Node, httpCode int) (Fallback, bool) {
	for _, fallback := range node.Fallbacks() {
		if fallback.HttpCode() == httpCode {
			return fallback, true
		}
	}

	return nil, false
}

func (r *renderer) buildAssets(page renderable_pages.Page) (assets.Assets, error) {
	css, err := r.assetBuilder.
		Create().
		WithType(assets.AssetTypeCSS).
		WithURL("/assets/" + page.Keyname() + ".css").
		Now()

	if err != nil {
		return nil, err
	}

	js, err := r.assetBuilder.
		Create().
		WithType(assets.AssetTypeJavaScript).
		WithURL("/assets/" + page.Keyname() + ".js").
		Now()

	if err != nil {
		return nil, err
	}

	return r.assetsBuilder.
		Create().
		AddAsset(css).
		AddAsset(js).
		Now()
}

func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

func isParam(value string) bool {
	return strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")
}
