package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/steve-rodrigue/eventflow/applications/events"
	applicationhttps "github.com/steve-rodrigue/eventflow/applications/servers/https"
	domainevents "github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/events/contexts"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
	"github.com/steve-rodrigue/eventflow/infrastructure/https"
)

func NewTargetResult(operationType results.OperationType, element string, html string) (results.Result, error) {
	operation, err := NewTargetOperation(operationType, element, html)
	if err != nil {
		return nil, err
	}

	return results.NewBuilder().
		Create().
		AddOperation(operation).
		Now()
}

func NewResultWithOperations(operations ...results.Operation) (results.Result, error) {
	return results.NewBuilder().
		Create().
		WithOperations(operations).
		Now()
}

func NewTargetOperation(operationType results.OperationType, element string, html string) (results.Operation, error) {
	target, err := results.NewTargetBuilder().
		Create().
		WithElement(element).
		WithHTML(html).
		Now()

	if err != nil {
		return nil, err
	}

	action, err := results.NewActionBuilder().
		Create().
		WithTarget(target).
		Now()

	if err != nil {
		return nil, err
	}

	return results.NewOperationBuilder().
		Create().
		WithType(operationType).
		WithAction(action).
		Now()
}

func RenderUserCard(userID string) string {
	return fmt.Sprintf(
		`<div id="user-card"><strong>User %s</strong> was updated from the server.</div>`,
		template.HTMLEscapeString(userID),
	)
}

func WebSocketHandler(app events.Application) http.Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("websocket upgrade error:", err)
			return
		}
		defer conn.Close()

		log.Println("browser connected")

		for {
			var message events.IncomingMessage

			if err := conn.ReadJSON(&message); err != nil {
				log.Println("browser disconnected:", err)
				return
			}

			outgoing, err := app.Execute(message)
			if err != nil {
				_ = conn.WriteJSON(events.OutgoingMessage{
					Type:  "error",
					Error: err.Error(),
				})
				continue
			}

			if err := conn.WriteJSON(outgoing); err != nil {
				log.Println("websocket write error:", err)
				return
			}
		}
	})
}

func IndexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "index.html")
	})
}

func MustEvent(event domainevents.Event, err error) domainevents.Event {
	if err != nil {
		log.Fatal(err)
	}

	return event
}

func main() {
	registry := domainevents.NewRegistry()

	err := registry.ListenAll([]domainevents.Event{
		MustEvent(domainevents.NewBuilder().
			Create().
			WithKeyname("user_updated").
			WithEventName("user_updated").
			WithAction(func(ctx contexts.Context) (results.Result, error) {
				userID, _ := ctx.Value("userId").(string)

				replaceUserCard, err := NewTargetOperation(
					results.OperationTypeReplace,
					"#user-card",
					RenderUserCard(userID),
				)
				if err != nil {
					return nil, err
				}

				appendMessage, err := NewTargetOperation(
					results.OperationTypeAppend,
					"#content",
					`<p>Server executed event: user_updated</p>`,
				)
				if err != nil {
					return nil, err
				}

				return NewResultWithOperations(replaceUserCard, appendMessage)
			}).
			Now()),

		MustEvent(domainevents.NewBuilder().
			Create().
			WithKeyname("counter.increment").
			WithEventName("counter.increment").
			WithAction(func(ctx contexts.Context) (results.Result, error) {
				amount := ctx.Value("amount")

				return NewTargetResult(
					results.OperationTypeAppend,
					"#content",
					fmt.Sprintf(`<p>Server incremented counter by %v</p>`, amount),
				)
			}).
			Now()),

		MustEvent(domainevents.NewBuilder().
			Create().
			WithKeyname("message.send").
			WithEventName("message.send").
			WithAction(func(ctx contexts.Context) (results.Result, error) {
				text, _ := ctx.Value("text").(string)

				return NewTargetResult(
					results.OperationTypeAppend,
					"#content",
					fmt.Sprintf(`<p>Server received message: %s</p>`, template.HTMLEscapeString(text)),
				)
			}).
			Now()),
	})

	if err != nil {
		log.Fatal(err)
	}

	app := events.NewApplication(
		contexts.NewBuilder(),
		registry,
	)

	apiHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/api").
		WithHandle(WebSocketHandler(app)).
		Now()
	if err != nil {
		log.Fatal(err)
	}

	indexHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/").
		WithHandle(IndexHandler()).
		Now()
	if err != nil {
		log.Fatal(err)
	}

	server := https.NewServer(":8080", []applicationhttps.Handler{
		apiHandler,
		indexHandler,
	})

	go func() {
		log.Println("server started at http://localhost:8080")

		if err := server.Start(); err != nil {
			log.Fatal("server error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatal("server shutdown error:", err)
	}

	log.Println("server stopped")
}
