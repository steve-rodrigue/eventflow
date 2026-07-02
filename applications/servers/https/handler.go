package https

import "net/http"

type handler struct {
	path   string
	handle http.Handler
}

func NewHandler(path string, handle http.Handler) Handler {
	return &handler{
		path:   path,
		handle: handle,
	}
}

func (h *handler) Path() string {
	return h.path
}

func (h *handler) Handle() http.Handler {
	return h.handle
}
