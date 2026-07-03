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
	appevents "github.com/steve-rodrigue/eventflow/applications/events"
	applicationhttps "github.com/steve-rodrigue/eventflow/applications/servers/https"
	domainevents "github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/events/contexts"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	renderablepages "github.com/steve-rodrigue/eventflow/domain/renderables/pages"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/components"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/heads/assets"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
	domaintrees "github.com/steve-rodrigue/eventflow/domain/trees"
	renderedpages "github.com/steve-rodrigue/eventflow/domain/trees/pages"
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

func WebSocketHandler(app appevents.Application) http.Handler {
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
			var message appevents.IncomingMessage

			if err := conn.ReadJSON(&message); err != nil {
				log.Println("browser disconnected:", err)
				return
			}

			outgoing, err := app.Execute(message)
			if err != nil {
				_ = conn.WriteJSON(appevents.OutgoingMessage{
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

func IndexHandler(treeRenderer domaintrees.Renderer, tree domaintrees.Tree) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request, err := domaintrees.NewRequestBuilder().
			Create().
			WithPath(r.URL.Path).
			WithMethod(r.Method).
			WithLocale("en").
			WithTarget("desktop").
			Now()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rendered, err := treeRenderer.Render(tree, request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		for _, header := range rendered.Headers() {
			w.Header().Set(header.Name(), header.Value())
		}

		w.WriteHeader(rendered.HttpCode())
		_, _ = w.Write([]byte(rendered.Body()))
	})
}

func AssetsHandler(pageRenderer renderablepages.Renderer, page renderablepages.Page) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/assets/home.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			_, _ = w.Write([]byte(pageRenderer.RenderStyle(page, renderables.Params{})))

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

func MustEvent(event domainevents.Event, err error) domainevents.Event {
	if err != nil {
		log.Fatal(err)
	}

	return event
}

func MustTemplate(template templates.Template, err error) templates.Template {
	if err != nil {
		log.Fatal(err)
	}

	return template
}

func MustHead(head heads.Head, err error) heads.Head {
	if err != nil {
		log.Fatal(err)
	}

	return head
}

func MustComponent(component components.Component, err error) components.Component {
	if err != nil {
		log.Fatal(err)
	}

	return component
}

func MustPage(page renderablepages.Page, err error) renderablepages.Page {
	if err != nil {
		log.Fatal(err)
	}

	return page
}

func MustRoute(route domaintrees.Route, err error) domaintrees.Route {
	if err != nil {
		log.Fatal(err)
	}

	return route
}

func MustResource(resource domaintrees.Resource, err error) domaintrees.Resource {
	if err != nil {
		log.Fatal(err)
	}

	return resource
}

func MustGroup(group domaintrees.Group, err error) domaintrees.Group {
	if err != nil {
		log.Fatal(err)
	}

	return group
}

func MustTarget(target domaintrees.Target, err error) domaintrees.Target {
	if err != nil {
		log.Fatal(err)
	}

	return target
}

func MustNode(node domaintrees.Node, err error) domaintrees.Node {
	if err != nil {
		log.Fatal(err)
	}

	return node
}

func MustTree(tree domaintrees.Tree, err error) domaintrees.Tree {
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func BuildTree() (domaintrees.Tree, renderablepages.Page) {
	pageTemplate := MustTemplate(templates.NewMustacheBuilder().
		Create().
		WithKeyname("page").
		WithCode(`<!doctype html>
<html lang="{language}">
<head>
{head}
</head>
<body>
{body}
</body>
</html>`).
		Now())

	head := MustHead(heads.NewBuilder().
		Create().
		WithTitle(MustTemplate(templates.NewMustacheBuilder().
			Create().
			WithKeyname("title").
			WithCode("EventFlow WebSocket Test").
			Now())).
		WithDescription(MustTemplate(templates.NewMustacheBuilder().
			Create().
			WithKeyname("description").
			WithCode("EventFlow server-driven UI test page.").
			Now())).
		WithLanguage("en").
		Now())

	body := MustComponent(components.NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("body").
		WithTemplate(MustTemplate(templates.NewMustacheBuilder().
			Create().
			WithKeyname("body").
			WithCode(`
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
`).
			Now())).
		Now())

	page := MustPage(renderablepages.NewBuilder(renderables.NewBuilder()).
		Create().
		WithLanguage("en").
		WithKeyname("home").
		WithTemplate(pageTemplate).
		WithHead(head).
		WithBody(body).
		Now())

	resource := MustResource(domaintrees.NewResourceBuilder().
		Create().
		WithLocale("en").
		WithRoute(MustRoute(domaintrees.NewRouteBuilder().
			Create().
			WithPattern("/").
			Now())).
		WithPage(page).
		Now())

	group := MustGroup(domaintrees.NewGroupBuilder().
		Create().
		WithKeyname("default").
		AddResource(resource).
		Now())

	target := MustTarget(domaintrees.NewTargetBuilder().
		Create().
		WithKeyname("desktop").
		AddGroup(group).
		Now())

	node := MustNode(domaintrees.NewNodeBuilder().
		Create().
		AddTarget(target).
		Now())

	tree := MustTree(domaintrees.NewBuilder().
		Create().
		WithKeyname("main").
		AddNode(node).
		Now())

	return tree, page
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

	app := appevents.NewApplication(
		contexts.NewBuilder(),
		registry,
	)

	templateRenderer := templates.NewMustacheRenderer()
	assetsRenderer := assets.NewRenderer()
	headRenderer := heads.NewRenderer(templateRenderer, assetsRenderer)
	componentRenderer := components.NewRenderer(templateRenderer)

	pageRenderer := renderablepages.NewRenderer(
		templateRenderer,
		headRenderer,
		componentRenderer,
	)

	treeRenderer := domaintrees.NewRenderer(
		pageRenderer,
		renderedpages.NewBuilder(),
		renderedpages.NewHeaderBuilder(),
		assets.NewBuilder(),
		assets.NewAssetBuilder(),
	)

	tree, page := BuildTree()

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
		WithHandle(AssetsHandler(pageRenderer, page)).
		Now()
	if err != nil {
		log.Fatal(err)
	}

	indexHandler, err := applicationhttps.NewHandlerBuilder().
		Create().
		WithPath("/").
		WithHandle(IndexHandler(treeRenderer, tree)).
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
