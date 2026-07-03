package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/steve-rodrigue/eventflow/applications"
	eventflow "github.com/steve-rodrigue/eventflow/cmd/eventflow"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
)

type Channel struct {
	ID          string
	Name        string
	Description string
	Emoji       string
	Path        string
	Online      int
}

type Message struct {
	ChannelID string
	User      string
	Avatar    string
	Text      string
	Time      string
	Own       bool
}

var channels = []Channel{
	{
		ID:          "general",
		Name:        "General",
		Description: "Company-wide discussion and daily updates.",
		Emoji:       "💬",
		Path:        "/",
		Online:      12,
	},
	{
		ID:          "engineering",
		Name:        "Engineering",
		Description: "Architecture, backend, infrastructure and EventFlow.",
		Emoji:       "⚙️",
		Path:        "/channels/engineering",
		Online:      7,
	},
	{
		ID:          "design",
		Name:        "Design",
		Description: "Product design, UI polish and user experience.",
		Emoji:       "🎨",
		Path:        "/channels/design",
		Online:      4,
	},
	{
		ID:          "random",
		Name:        "Random",
		Description: "Casual conversations and fun ideas.",
		Emoji:       "✨",
		Path:        "/channels/random",
		Online:      9,
	},
}

var messages = []Message{
	{
		ChannelID: "general",
		User:      "Maya Chen",
		Avatar:    "MC",
		Text:      "Welcome to the SteveCare live chat. This UI is rendered by EventFlow.",
		Time:      "09:41",
	},
	{
		ChannelID: "general",
		User:      "Steve Rodrigue",
		Avatar:    "SR",
		Text:      "Every message is sent through the event system and updates the DOM in real time.",
		Time:      "09:43",
		Own:       true,
	},
	{
		ChannelID: "engineering",
		User:      "Luc Tremblay",
		Avatar:    "LT",
		Text:      "The channel routing is static for now, but later it can come from the database.",
		Time:      "10:05",
	},
	{
		ChannelID: "design",
		User:      "Maya Chen",
		Avatar:    "MC",
		Text:      "The glass panels and gradients make the chat feel premium without being too heavy.",
		Time:      "10:18",
	},
	{
		ChannelID: "random",
		User:      "Nora",
		Avatar:    "NO",
		Text:      "Random channel is ready for experiments.",
		Time:      "10:27",
	},
}

func BuildTree() applications.Tree {
	groups := []applications.Group{
		group("channel-general",
			resource("en", "/", chatPage(findChannel("general"))),
		),
	}

	for _, channel := range channels {
		if channel.ID == "general" {
			continue
		}

		groups = append(groups, group("channel-"+channel.ID,
			resource("en", channel.Path, chatPage(channel)),
		))
	}

	return applications.Tree{
		Keyname: "stevecare-chat",
		Nodes: []applications.Node{
			{
				Targets: []applications.Target{
					{
						Keyname: "desktop",
						Groups:  groups,
					},
				},
			},
		},
	}
}

func chatEvents(register bool) []applications.Event {
	if !register {
		return nil
	}

	return []applications.Event{
		{
			Keyname: "chat.message.send",
			Action: func(ctx applications.Context) (*applications.Result, error) {
				channelID, _ := ctx.Payload["channel"].(string)
				text, _ := ctx.Payload["text"].(string)

				if channelID == "" {
					channelID = "general"
				}
				if text == "" {
					text = "New message from EventFlow."
				}

				return &applications.Result{
					Operations: []applications.Operation{
						{
							Type: results.OperationTypeAppend,
							Action: applications.Action{
								Target: &applications.ActionTarget{
									Element: "#messages-" + channelID,
									HTML: renderMessage(Message{
										ChannelID: channelID,
										User:      "Steve Rodrigue",
										Avatar:    "SR",
										Text:      text,
										Time:      time.Now().Format("15:04"),
										Own:       true,
									}),
								},
							},
						},
						{
							Type: results.OperationTypeReplace,
							Action: applications.Action{
								Target: &applications.ActionTarget{
									Element: "#chat-input-" + channelID,
									HTML:    fmt.Sprintf(`<input id="chat-input-%s" type="text" placeholder="Write a message..." />`, escape(channelID)),
								},
							},
						},
					},
				}, nil
			},
		},
	}
}

func chatPage(channel Channel) applications.Page {
	return applications.Page{
		Keyname:  pageKeyname(channel),
		Language: "en",
		Renderable: applications.Renderable{
			Template: applications.Template{
				Keyname: "layout_" + channel.ID,
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
				Keyname: "title_" + channel.ID,
				Code:    channel.Name + " · SteveCare Chat",
			},
			Description: applications.Template{
				Keyname: "description_" + channel.ID,
				Code:    channel.Description,
			},
			Language: "en",
			Meta: []applications.Meta{
				{
					Name: "theme-color",
					Content: applications.Template{
						Keyname: "theme_color_" + channel.ID,
						Code:    "#0f172a",
					},
				},
			},
			OpenGraph: &applications.OpenGraph{
				Title: applications.Template{
					Keyname: "og_title_" + channel.ID,
					Code:    channel.Name + " · SteveCare Chat",
				},
				Description: applications.Template{
					Keyname: "og_description_" + channel.ID,
					Code:    channel.Description,
				},
				Image: applications.Template{
					Keyname: "og_image_" + channel.ID,
					Code:    "/assets/stevecare-chat.png",
				},
				URL: applications.Template{
					Keyname: "og_url_" + channel.ID,
					Code:    "https://steve.care",
				},
				Type: "website",
				SiteName: applications.Template{
					Keyname: "og_site_" + channel.ID,
					Code:    "SteveCare Chat",
				},
			},
			TwitterCard: &applications.TwitterCard{
				Card: "summary_large_image",
				Title: applications.Template{
					Keyname: "twitter_title_" + channel.ID,
					Code:    channel.Name + " · SteveCare Chat",
				},
				Description: applications.Template{
					Keyname: "twitter_description_" + channel.ID,
					Code:    channel.Description,
				},
				Image: applications.Template{
					Keyname: "twitter_image_" + channel.ID,
					Code:    "/assets/stevecare-chat.png",
				},
			},
		},
		Body: applications.Component{
			Keyname: "body",
			StylableRenderable: applications.StylableRenderable{
				Style: &applications.Template{
					Keyname: "style_" + channel.ID,
					Code:    siteCSS(),
				},
				Renderable: applications.Renderable{
					Template: applications.Template{
						Keyname: "body_" + channel.ID,
						Code:    chatBody(channel),
					},
				},
			},
			Events: chatEvents(channel.ID == "general"),
		},
	}
}

func chatBody(channel Channel) string {
	return fmt.Sprintf(`
<div class="app-shell">
	<aside class="sidebar">
		<a class="brand" href="/">
			<span class="brand-mark">SC</span>
			<span>
				<strong>SteveCare</strong>
				<small>Live Chat</small>
			</span>
		</a>

		<div class="sidebar-section">
			<div class="section-title">Channels</div>
			<nav class="channel-list">
				%s
			</nav>
		</div>

		<div class="sidebar-section">
			<div class="section-title">Online now</div>
			%s
		</div>
	</aside>

	<main class="chat">
		<header class="chat-header">
			<div>
				<div class="breadcrumb">
					<a href="/">SteveCare</a>
					<span>/</span>
					<span>%s</span>
				</div>

				<h1><span>%s</span> %s</h1>
				<p>%s</p>
			</div>

			<div class="channel-stats">
				<span class="pulse"></span>
				<strong>%d</strong>
				<small>online</small>
			</div>
		</header>

		<section id="messages-%s" class="messages">
			%s
		</section>

		<section class="composer">
			<div class="composer-input">
				<input id="chat-input-%s" type="text" placeholder="Write a message..." />
				<button
					data-event="chat.message.send"
					data-channel="%s"
					data-user="steve"
					data-payload-text="#chat-input-%s">
					Send
				</button>
			</div>
		</section>
	</main>
</div>
`,
		renderChannels(channel.ID),
		renderOnlineUsers(),
		escape(channel.Name),
		escape(channel.Emoji),
		escape(channel.Name),
		escape(channel.Description),
		channel.Online,
		escape(channel.ID), // section id: messages-{id}
		renderMessages(channel.ID),
		escape(channel.ID), // input id: chat-input-{id}
		escape(channel.ID), // data-channel
		escape(channel.ID), // data-payload-text selector
	)
}

func renderChannels(activeID string) string {
	html := ""

	for _, channel := range channels {
		active := ""
		if channel.ID == activeID {
			active = " active"
		}

		html += fmt.Sprintf(`
<a class="channel%s" href="%s">
	<span class="channel-emoji">%s</span>
	<span>
		<strong>%s</strong>
		<small>%d online</small>
	</span>
</a>`,
			active,
			escape(channel.Path),
			escape(channel.Emoji),
			escape(channel.Name),
			channel.Online,
		)
	}

	return html
}

func renderOnlineUsers() string {
	users := []struct {
		Name   string
		Role   string
		Avatar string
	}{
		{"Steve Rodrigue", "Architect", "SR"},
		{"Maya Chen", "Designer", "MC"},
		{"Luc Tremblay", "Engineer", "LT"},
		{"Nora Fields", "Product", "NO"},
	}

	html := `<div class="people">`

	for _, user := range users {
		html += fmt.Sprintf(`
<div class="person">
	<span class="avatar">%s</span>
	<span>
		<strong>%s</strong>
		<small>%s</small>
	</span>
	<span class="dot"></span>
</div>`,
			escape(user.Avatar),
			escape(user.Name),
			escape(user.Role),
		)
	}

	return html + `</div>`
}

func renderMessages(channelID string) string {
	html := ""

	for _, message := range messages {
		if message.ChannelID != channelID {
			continue
		}

		html += renderMessage(message)
	}

	return html
}

func renderMessage(message Message) string {
	own := ""
	if message.Own {
		own = " own"
	}

	return fmt.Sprintf(`
<article class="message%s">
	<span class="avatar">%s</span>
	<div class="bubble">
		<header>
			<strong>%s</strong>
			<time>%s</time>
		</header>
		<p>%s</p>
	</div>
</article>`,
		own,
		escape(message.Avatar),
		escape(message.User),
		escape(message.Time),
		escape(message.Text),
	)
}

func group(keyname string, resources ...applications.Resource) applications.Group {
	return applications.Group{
		Keyname:   keyname,
		Resources: resources,
	}
}

func resource(locale string, pattern string, page applications.Page) applications.Resource {
	return applications.Resource{
		Locale: locale,
		Route: applications.Route{
			Pattern: pattern,
		},
		Page: page,
	}
}

func findChannel(id string) Channel {
	for _, channel := range channels {
		if channel.ID == id {
			return channel
		}
	}

	return channels[0]
}

func pageKeyname(channel Channel) string {
	if channel.ID == "general" {
		return "home"
	}

	return "channel-" + channel.ID
}

func escape(value string) string {
	return template.HTMLEscapeString(value)
}

func siteCSS() string {
	return `
:root {
	color-scheme: dark;
	--bg: #080b12;
	--bg-2: #0f172a;
	--surface: rgba(15, 23, 42, .78);
	--surface-2: rgba(30, 41, 59, .78);
	--surface-3: rgba(51, 65, 85, .75);
	--border: rgba(148, 163, 184, .2);
	--text: #f8fafc;
	--muted: #94a3b8;
	--brand: #38bdf8;
	--brand-2: #a78bfa;
	--green: #22c55e;
	--shadow: 0 28px 90px rgba(0, 0, 0, .45);
}

* {
	box-sizing: border-box;
}

html,
body {
	min-height: 100%;
}

body {
	margin: 0;
	font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
	background:
		radial-gradient(circle at 18% 8%, rgba(56, 189, 248, .24), transparent 28rem),
		radial-gradient(circle at 92% 18%, rgba(167, 139, 250, .22), transparent 30rem),
		linear-gradient(135deg, #080b12, #0f172a 48%, #111827);
	color: var(--text);
}

a {
	color: inherit;
	text-decoration: none;
}

button {
	border: 0;
	cursor: pointer;
}

.app-shell {
	display: grid;
	grid-template-columns: 320px minmax(0, 1fr);
	width: min(1440px, calc(100% - 32px));
	min-height: calc(100vh - 32px);
	margin: 16px auto;
	border: 1px solid var(--border);
	border-radius: 34px;
	overflow: hidden;
	background: rgba(2, 6, 23, .58);
	box-shadow: var(--shadow);
	backdrop-filter: blur(18px);
}

.sidebar {
	padding: 24px;
	background:
		linear-gradient(180deg, rgba(15, 23, 42, .94), rgba(15, 23, 42, .66)),
		rgba(15, 23, 42, .72);
	border-right: 1px solid var(--border);
}

.brand {
	display: flex;
	align-items: center;
	gap: 14px;
	margin-bottom: 34px;
}

.brand-mark {
	display: grid;
	place-items: center;
	width: 52px;
	height: 52px;
	border-radius: 18px;
	background: linear-gradient(135deg, var(--brand), var(--brand-2));
	font-weight: 950;
	letter-spacing: -.04em;
	box-shadow: 0 18px 44px rgba(56, 189, 248, .18);
}

.brand small,
.channel small,
.person small,
.channel-stats small,
.status {
	display: block;
	color: var(--muted);
	margin-top: 3px;
	font-size: .82rem;
}

.sidebar-section {
	margin-top: 28px;
}

.section-title {
	margin-bottom: 12px;
	color: var(--muted);
	font-size: .78rem;
	font-weight: 800;
	letter-spacing: .12em;
	text-transform: uppercase;
}

.channel-list,
.people {
	display: grid;
	gap: 10px;
}

.channel,
.person {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 12px;
	border: 1px solid transparent;
	border-radius: 18px;
	color: var(--text);
	transition: .18s ease;
}

.channel:hover,
.channel.active,
.person {
	background: var(--surface-2);
	border-color: var(--border);
}

.channel.active {
	box-shadow: inset 3px 0 0 var(--brand);
}

.channel-emoji {
	display: grid;
	place-items: center;
	width: 42px;
	height: 42px;
	border-radius: 15px;
	background: rgba(255, 255, 255, .08);
	font-size: 1.2rem;
}

.avatar {
	display: grid;
	place-items: center;
	flex: 0 0 auto;
	width: 42px;
	height: 42px;
	border-radius: 15px;
	background: linear-gradient(135deg, var(--brand), var(--brand-2));
	color: white;
	font-weight: 900;
	font-size: .9rem;
}

.dot,
.pulse {
	width: 10px;
	height: 10px;
	border-radius: 999px;
	background: var(--green);
	box-shadow: 0 0 0 5px rgba(34, 197, 94, .15);
	margin-left: auto;
}

.chat {
	display: grid;
	grid-template-rows: auto minmax(0, 1fr) auto;
	min-height: calc(100vh - 32px);
}

.chat-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 24px;
	padding: 34px 38px 26px;
	border-bottom: 1px solid var(--border);
	background: rgba(15, 23, 42, .42);
}

.breadcrumb {
	display: flex;
	align-items: center;
	gap: 8px;
	color: var(--muted);
	font-size: .9rem;
	margin-bottom: 12px;
}

.chat-header h1 {
	display: flex;
	align-items: center;
	gap: 12px;
	margin: 0;
	font-size: clamp(2rem, 4vw, 4rem);
	line-height: .95;
	letter-spacing: -.06em;
}

.chat-header p {
	max-width: 660px;
	color: var(--muted);
	line-height: 1.7;
	margin: 14px 0 0;
}

.channel-stats {
	display: grid;
	place-items: center;
	min-width: 124px;
	padding: 18px;
	border-radius: 22px;
	background: var(--surface-2);
	border: 1px solid var(--border);
}

.channel-stats strong {
	font-size: 2rem;
	line-height: 1;
	margin-top: 8px;
}

.channel-stats .pulse {
	margin-left: 0;
}

.messages {
	display: flex;
	flex-direction: column;
	gap: 16px;
	padding: 34px 38px;
	overflow: auto;
}

.message {
	display: flex;
	gap: 14px;
	max-width: 760px;
}

.message.own {
	align-self: flex-end;
	flex-direction: row-reverse;
}

.bubble {
	padding: 16px 18px;
	border-radius: 24px;
	background: var(--surface-2);
	border: 1px solid var(--border);
}

.message.own .bubble {
	background: linear-gradient(135deg, rgba(56, 189, 248, .26), rgba(167, 139, 250, .24));
	border-color: rgba(56, 189, 248, .26);
}

.bubble header {
	display: flex;
	gap: 12px;
	align-items: baseline;
	margin-bottom: 8px;
}

.bubble time {
	color: var(--muted);
	font-size: .78rem;
}

.bubble p {
	margin: 0;
	color: #dbeafe;
	line-height: 1.6;
}

.composer {
	padding: 24px 38px;
	border-top: 1px solid var(--border);
	background: rgba(15, 23, 42, .72);
}

.composer-input {
	display: grid;
	grid-template-columns: minmax(0, 1fr) auto;
	gap: 12px;
	width: 100%;
	padding: 10px;
	border: 1px solid var(--border);
	border-radius: 24px;
	background: rgba(2, 6, 23, .46);
	box-shadow: inset 0 1px 0 rgba(255, 255, 255, .05);
}

.composer-input input {
	width: 100%;
	min-height: 52px;
	padding: 0 18px;
	border: 0;
	outline: 0;
	border-radius: 18px;
	background: rgba(15, 23, 42, .86);
	color: var(--text);
	font: inherit;
}

.composer-input input::placeholder {
	color: var(--muted);
}

.composer-input input:focus {
	box-shadow: 0 0 0 3px rgba(56, 189, 248, .18);
}

.composer-input button {
	min-height: 52px;
	padding: 0 22px;
	border-radius: 18px;
	background: linear-gradient(135deg, var(--brand), var(--brand-2));
	color: white;
	font-weight: 900;
	box-shadow: 0 16px 40px rgba(56, 189, 248, .18);
}

.composer-actions {
	display: flex;
	flex-wrap: wrap;
	gap: 10px;
	justify-content: flex-end;
}

.composer button {
	padding: 12px 16px;
	border-radius: 999px;
	background: linear-gradient(135deg, var(--brand), var(--brand-2));
	color: white;
	font-weight: 850;
	box-shadow: 0 16px 40px rgba(56, 189, 248, .18);
}

.status.success {
	color: var(--green);
}

@media (max-width: 900px) {
	.app-shell {
		grid-template-columns: 1fr;
	}

	.sidebar {
		border-right: 0;
		border-bottom: 1px solid var(--border);
	}

	.chat-header,
	.composer {
		align-items: flex-start;
		flex-direction: column;
	}

	.composer-actions {
		justify-content: flex-start;
	}

	.messages,
	.chat-header,
	.composer {
		padding-left: 22px;
		padding-right: 22px;
	}
}
`
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
