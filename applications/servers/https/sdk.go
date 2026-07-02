package https

import (
	"context"
	"net/http"
)

// NewHandlerBuilder creates a new handler builder
func NewHandlerBuilder() HandlerBuilder {
	return &handlerBuilder{
		path:   "",
		handle: nil,
	}
}

// Server represents an HTTP server.
type Server interface {
	Start() error
	Stop(ctx context.Context) error
	Addr() string
}

// HandlerBuilder creates an handler builder
type HandlerBuilder interface {
	Create() HandlerBuilder
	WithPath(path string) HandlerBuilder
	WithHandle(handle http.Handler) HandlerBuilder
	Now() (Handler, error)
}

// Handler represents an HTTP route handler.
type Handler interface {
	Path() string
	Handle() http.Handler
}
