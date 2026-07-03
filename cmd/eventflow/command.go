package eventflow

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/steve-rodrigue/eventflow/applications"
	eventapps "github.com/steve-rodrigue/eventflow/applications/events"
	"github.com/steve-rodrigue/eventflow/applications/servers/https"
	applicationhttps "github.com/steve-rodrigue/eventflow/applications/servers/https"
)

type Command struct {
	app            applications.Application
	server         https.Server
	address        string
	assetsBasePath string
	logger         *log.Logger
}

func (c *Command) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		c.logger.Printf("server started at http://localhost%s", c.address)
		errCh <- c.server.Start()
	}()

	select {
	case <-ctx.Done():
		c.logger.Println("stopping server...")
		return c.server.Stop(context.Background())

	case err := <-errCh:
		return err
	}
}

func (c *Command) handlers() ([]applicationhttps.Handler, error) {
	apiHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/api").
		WithHandle(c.webSocketHandler()).
		Now()
	if err != nil {
		return nil, err
	}

	assetsHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath(fmt.Sprintf("%s/", c.assetsBasePath)).
		WithHandle(c.assetsHandler()).
		Now()
	if err != nil {
		return nil, err
	}

	pageHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/").
		WithHandle(c.pageHandler()).
		Now()
	if err != nil {
		return nil, err
	}

	return []applicationhttps.Handler{
		apiHandler,
		assetsHandler,
		pageHandler,
	}, nil
}

func (c *Command) webSocketHandler() http.Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			c.logger.Printf("websocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		c.logger.Printf("websocket connected: %s", r.RemoteAddr)

		for {
			var message eventapps.IncomingMessage

			if err := conn.ReadJSON(&message); err != nil {
				c.logger.Printf("websocket disconnected: %s: %v", r.RemoteAddr, err)
				return
			}

			outgoing, err := c.app.Trigger(message)
			if err != nil {
				c.logger.Printf("event failed: event=%q remote=%s error=%v", message.Event, r.RemoteAddr, err)

				_ = conn.WriteJSON(eventapps.OutgoingMessage{
					Type:  "error",
					Error: err.Error(),
				})

				continue
			}

			if err := conn.WriteJSON(outgoing); err != nil {
				c.logger.Printf("websocket write failed: remote=%s error=%v", r.RemoteAddr, err)
				return
			}
		}
	})
}

func (c *Command) pageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rendered, err := c.app.RenderPage(applications.RouteRequest{
			Path:   r.URL.Path,
			Method: r.Method,
			Locale: "en",
			Target: "desktop",
		})
		if err != nil {
			c.logger.Printf("page render failed: method=%s path=%s remote=%s error=%v", r.Method, r.URL.Path, r.RemoteAddr, err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		writeRenderedPage(w, rendered)
	})
}

func (c *Command) assetsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			rendered *applications.RenderedPage
			err      error
		)

		switch {
		case hasSuffix(r.URL.Path, ".css"):
			rendered, err = c.app.RenderStyle(r.URL.Path)

		case hasSuffix(r.URL.Path, ".js"):
			rendered, err = c.app.RenderJavascript(r.URL.Path)

		default:
			http.NotFound(w, r)
			return
		}

		if err != nil {
			c.logger.Printf("asset render failed: path=%s remote=%s error=%v", r.URL.Path, r.RemoteAddr, err)
			http.NotFound(w, r)
			return
		}

		writeRenderedPage(w, rendered)
	})
}

func writeRenderedPage(w http.ResponseWriter, rendered *applications.RenderedPage) {
	for _, header := range rendered.Headers {
		w.Header().Set(header.Name, header.Value)
	}

	w.WriteHeader(rendered.HttpCode)
	_, _ = w.Write([]byte(rendered.Body))
}

func hasSuffix(value string, suffix string) bool {
	return len(value) >= len(suffix) && value[len(value)-len(suffix):] == suffix
}
