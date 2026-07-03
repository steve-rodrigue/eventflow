package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/steve-rodrigue/eventflow/applications"
	eventflow "github.com/steve-rodrigue/eventflow/cmd/eventflow"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
)

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
													Style: &applications.Template{
														Keyname: "body_style",
														Code: `
													body {
														margin: 0;
														font-family: system-ui, sans-serif;
														background: #111827;
														color: #f9fafb;
													}
													
													#app {
														max-width: 760px;
														margin: 64px auto;
														padding: 32px;
														background: #1f2937;
														border-radius: 16px;
													}
													
													button {
														margin: 8px 8px 8px 0;
														padding: 10px 14px;
														border: 0;
														border-radius: 8px;
														cursor: pointer;
														font-weight: 600;
													}
													
													#user-card,
													#content {
														margin-top: 20px;
														padding: 16px;
														background: #374151;
														border-radius: 10px;
													}
													
													#content p {
														margin: 8px 0;
													}
													`,
													},
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
	assetsBasePath := "/assets"

	command, err := eventflow.New(
		":8080",
		assetsBasePath,
		BuildTree(),
		log.Default(),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := command.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
