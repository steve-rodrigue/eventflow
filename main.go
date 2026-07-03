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
	"github.com/steve-rodrigue/eventflow/applications"
	eventapps "github.com/steve-rodrigue/eventflow/applications/events"
	applicationhttps "github.com/steve-rodrigue/eventflow/applications/servers/https"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
	"github.com/steve-rodrigue/eventflow/infrastructure/https"
)

func WebSocketHandler(app applications.Application) http.Handler {
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
			var message eventapps.IncomingMessage

			if err := conn.ReadJSON(&message); err != nil {
				log.Println("browser disconnected:", err)
				return
			}

			outgoing, err := app.Trigger(message)
			if err != nil {
				_ = conn.WriteJSON(eventapps.OutgoingMessage{
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

func IndexHandler(app applications.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rendered, err := app.Route(applications.RouteRequest{
			Path:   r.URL.Path,
			Method: r.Method,
			Locale: "en",
			Target: "desktop",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		for _, header := range rendered.Headers {
			w.Header().Set(header.Name, header.Value)
		}

		w.WriteHeader(rendered.HttpCode)
		_, _ = w.Write([]byte(rendered.Body))
	})
}

func AssetsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/assets/home.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			_, _ = w.Write([]byte(""))

		case "/assets/home.js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			_, _ = w.Write([]byte(EventFlowRuntimeJS()))

		default:
			http.NotFound(w, r)
		}
	})
}

func EventFlowRuntimeJS() string {
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

func RenderUserCard(userID string) string {
	return fmt.Sprintf(
		`<div id="user-card"><strong>User %s</strong> was updated from the server.</div>`,
		template.HTMLEscapeString(userID),
	)
}

func BuildTree() applications.Tree {
	return applications.Tree{
		Keyname: "main",
		Nodes: []applications.Node{
			{
				Targets: []applications.Target{
					{
						Keyname: "desktop",
						Groups: []applications.Group{
							{
								Keyname: "default",
								Resources: []applications.Resource{
									{
										Locale: "en",
										Route: applications.Route{
											Pattern: "/",
										},
										Page: applications.Page{
											Keyname:  "home",
											Language: "en",
											Renderable: applications.Renderable{
												Template: applications.Template{
													Keyname: "page",
													Code: `<!doctype html>
<html lang="{language}">
<head>
{head}
</head>
<body>
{body}
</body>
</html>`,
												},
											},
											Head: applications.Head{
												Title: applications.Template{
													Keyname: "title",
													Code:    "EventFlow WebSocket Test",
												},
												Description: applications.Template{
													Keyname: "description",
													Code:    "EventFlow server-driven UI test page.",
												},
												Language: "en",
											},
											Body: applications.Component{
												Keyname: "body",
												StylableRenderable: applications.StylableRenderable{
													Renderable: applications.Renderable{
														Template: applications.Template{
															Keyname: "body",
															Code: `
<main id="app">
	<h1>EventFlow WebSocket Test</h1>

	<div id="user-card">
		No user updated yet.
	</div>

	<button data-event="counter.increment" data-amount="1">
		Increment
	</button>

	<button data-event="message.send" data-text="Hello from browser">
		Send Message
	</button>

	<button data-event="user_updated" data-user-id="123">
		Update User
	</button>

	<section id="content">
		<p>Waiting for operations...</p>
	</section>
</main>
`,
														},
													},
												},
												Events: []applications.Event{
													{
														Keyname: "user_updated",
														Action: func(ctx applications.Context) (*applications.Result, error) {
															userID, _ := ctx.Payload["userId"].(string)

															return &applications.Result{
																Operations: []applications.Operation{
																	{
																		Type: results.OperationTypeReplace,
																		Action: applications.Action{
																			Target: &applications.ActionTarget{
																				Element: "#user-card",
																				HTML:    RenderUserCard(userID),
																			},
																		},
																	},
																	{
																		Type: results.OperationTypeAppend,
																		Action: applications.Action{
																			Target: &applications.ActionTarget{
																				Element: "#content",
																				HTML:    `<p>Server executed event: user_updated</p>`,
																			},
																		},
																	},
																},
															}, nil
														},
													},
													{
														Keyname: "counter.increment",
														Action: func(ctx applications.Context) (*applications.Result, error) {
															amount := ctx.Payload["amount"]

															return &applications.Result{
																Operations: []applications.Operation{
																	{
																		Type: results.OperationTypeAppend,
																		Action: applications.Action{
																			Target: &applications.ActionTarget{
																				Element: "#content",
																				HTML:    fmt.Sprintf(`<p>Server incremented counter by %v</p>`, amount),
																			},
																		},
																	},
																},
															}, nil
														},
													},
													{
														Keyname: "message.send",
														Action: func(ctx applications.Context) (*applications.Result, error) {
															text, _ := ctx.Payload["text"].(string)

															return &applications.Result{
																Operations: []applications.Operation{
																	{
																		Type: results.OperationTypeAppend,
																		Action: applications.Action{
																			Target: &applications.ActionTarget{
																				Element: "#content",
																				HTML:    fmt.Sprintf(`<p>Server received message: %s</p>`, template.HTMLEscapeString(text)),
																			},
																		},
																	},
																},
															}, nil
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func main() {
	app := applications.NewDefaultApplication()

	if err := app.Execute(BuildTree()); err != nil {
		log.Fatal(err)
	}

	apiHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/api").
		WithHandle(WebSocketHandler(app)).
		Now()
	if err != nil {
		log.Fatal(err)
	}

	assetsHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/assets/").
		WithHandle(AssetsHandler()).
		Now()
	if err != nil {
		log.Fatal(err)
	}

	indexHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/").
		WithHandle(IndexHandler(app)).
		Now()
	if err != nil {
		log.Fatal(err)
	}

	server := https.NewServer(":8080", []applicationhttps.Handler{
		apiHandler,
		assetsHandler,
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
