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
	pageRenderer   renderable_pages.Renderer
	pageBuilder    rendered_pages.Builder
	headerBuilder  rendered_pages.HeaderBuilder
	assetsBuilder  assets.Builder
	assetBuilder   assets.AssetBuilder
	assetsBasePath string
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

func (r *renderer) Style(tree Tree, uri string) (rendered_pages.Page, error) {
	if tree == nil {
		return nil, errors.New("tree is required")
	}

	pagePath := r.assetURIToPagePath(uri, ".css")

	request, err := NewRequestBuilder().
		Create().
		WithPath(pagePath).
		WithMethod(http.MethodGet).
		WithLocale("en").
		WithTarget("desktop").
		Now()
	if err != nil {
		return nil, err
	}

	page, ok := r.findPage(tree, request)
	if !ok {
		return nil, errors.New("page not found")
	}

	body := r.pageRenderer.RenderStyle(page, renderables.Params{})

	contentType, err := r.headerBuilder.
		Create().
		WithName("Content-Type").
		WithValue("text/css; charset=utf-8").
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

func (r *renderer) Javascript(tree Tree, uri string) (rendered_pages.Page, error) {
	if tree == nil {
		return nil, errors.New("tree is required")
	}

	pagePath := r.assetURIToPagePath(uri, ".js")

	request, err := NewRequestBuilder().
		Create().
		WithPath(pagePath).
		WithMethod(http.MethodGet).
		WithLocale("en").
		WithTarget("desktop").
		Now()
	if err != nil {
		return nil, err
	}

	if _, ok := r.findPage(tree, request); !ok {
		return nil, errors.New("page not found")
	}

	contentType, err := r.headerBuilder.
		Create().
		WithName("Content-Type").
		WithValue("application/javascript; charset=utf-8").
		Now()
	if err != nil {
		return nil, err
	}

	return r.pageBuilder.
		Create().
		WithHttpCode(http.StatusOK).
		AddHeader(contentType).
		WithBody(runtimeJS()).
		Now()
}

func (r *renderer) buildAssets(page renderable_pages.Page) (assets.Assets, error) {
	css, err := r.assetBuilder.
		Create().
		WithType(assets.AssetTypeCSS).
		WithURL(r.assetURL(page, ".css")).
		Now()
	if err != nil {
		return nil, err
	}

	js, err := r.assetBuilder.
		Create().
		WithType(assets.AssetTypeJavaScript).
		WithURL(r.assetURL(page, ".js")).
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

func (r *renderer) assetURL(page renderable_pages.Page, extension string) string {
	return r.assetsBasePath + "/" + page.Keyname() + extension
}

func (r *renderer) assetURIToPagePath(uri string, extension string) string {
	path := strings.TrimPrefix(uri, r.assetsBasePath+"/")
	path = strings.TrimSuffix(path, extension)

	if path == "" || path == "home" {
		return "/"
	}

	return "/" + path
}

func normalizeAssetsBasePath(path string) string {
	path = strings.TrimSpace(path)
	path = "/" + strings.Trim(path, "/")

	if path == "/" {
		return "/assets"
	}

	return path
}

func runtimeJS() string {
	return `
const socket = new WebSocket("ws://localhost:8080/api");

function trigger(eventName, payload = {}) {
	if (socket.readyState !== WebSocket.OPEN) {
		console.warn("WebSocket is not connected yet");
		return;
	}

	socket.send(JSON.stringify({
		type: "event",
		event: eventName,
		payload
	}));
}

socket.addEventListener("open", () => {
	console.log("WebSocket connected");
});

socket.addEventListener("close", () => {
	console.log("WebSocket disconnected");
});

socket.addEventListener("error", (error) => {
	console.error("WebSocket error", error);
});

socket.addEventListener("message", (message) => {
	const data = JSON.parse(message.data);

	if (data.type === "operations") {
		applyOperations(data.operations);
	}
});

function applyOperations(operations) {
	for (const op of operations) {
		const target = op.target ? document.querySelector(op.target) : null;

		if (op.type === "replace" && target) {
			target.outerHTML = op.html;
		}

		if (op.type === "remove" && target) {
			target.remove();
		}

		if (op.type === "append" && target) {
			target.insertAdjacentHTML("beforeend", op.html);
		}

		if (op.type === "prepend" && target) {
			target.insertAdjacentHTML("afterbegin", op.html);
		}

		if (op.type === "navigate") {
			window.location.href = op.url;
		}
	}
}

document.addEventListener("click", (event) => {
	const element = event.target.closest("[data-event]");

	if (!element) {
		return;
	}

	const payload = {};

	for (const [key, value] of Object.entries(element.dataset)) {
		if (key === "event") {
			continue;
		}

		payload[key] = value;
	}

	trigger(element.dataset.event, payload);
});
`
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

func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

func isParam(value string) bool {
	return strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")
}
