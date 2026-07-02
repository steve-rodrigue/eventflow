package https

import (
	"errors"
	"net/http"
)

type handlerBuilder struct {
	path   string
	handle http.Handler
}

func (app *handlerBuilder) Create() HandlerBuilder {
	return &handlerBuilder{
		path:   "",
		handle: nil,
	}
}

func (app *handlerBuilder) WithPath(path string) HandlerBuilder {
	app.path = path
	return app
}

func (app *handlerBuilder) WithHandle(handle http.Handler) HandlerBuilder {
	app.handle = handle
	return app
}

func (app *handlerBuilder) Now() (Handler, error) {
	if app.path == "" {
		return nil, errors.New("the path is mandatory in order to build an Handler")
	}

	if app.handle == nil {
		return nil, errors.New("the handle is mandatory in order to build an Handler")
	}

	return &handler{
		path:   app.path,
		handle: app.handle,
	}, nil
}
