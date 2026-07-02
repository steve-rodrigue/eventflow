package https

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/steve-rodrigue/eventflow/applications/servers/https"
)

type server struct {
	addr       string
	httpServer *http.Server
}

func (s *server) Addr() string {
	return s.addr
}

func (s *server) Start() error {
	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err

}

func (s *server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// StopWithTimeout stops the server with a timeout.
func StopWithTimeout(server https.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return server.Stop(ctx)

}
