package https

import (
	"net/http"

	"github.com/steve-rodrigue/eventflow/applications/servers/https"
)

// NewServer creates a new HTTP server.
func NewServer(addr string, handlers []https.Handler) https.Server {
	mux := http.NewServeMux()
	for _, handler := range handlers {
		mux.Handle(handler.Path(), handler.Handle())
	}

	return &server{
		addr: addr,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}

}
