package events

import (
	"github.com/steve-rodrigue/eventflow/domain/events/contexts"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
)

// BrowserType represents a DOM event supported by the browser.
type BrowserType string

const (
	BrowserTypeClick       BrowserType = "click"
	BrowserTypeDoubleClick BrowserType = "dblclick"

	BrowserTypeMouseDown  BrowserType = "mousedown"
	BrowserTypeMouseUp    BrowserType = "mouseup"
	BrowserTypeMouseEnter BrowserType = "mouseenter"
	BrowserTypeMouseLeave BrowserType = "mouseleave"

	BrowserTypeInput  BrowserType = "input"
	BrowserTypeChange BrowserType = "change"
	BrowserTypeSubmit BrowserType = "submit"

	BrowserTypeFocus BrowserType = "focus"
	BrowserTypeBlur  BrowserType = "blur"

	BrowserTypeKeyDown BrowserType = "keydown"
	BrowserTypeKeyUp   BrowserType = "keyup"

	BrowserTypeScroll BrowserType = "scroll"
	BrowserTypeLoad   BrowserType = "load"
)

// ActionFn represents an action func
type ActionFn func(ctx contexts.Context) (results.Result, error)

// NewRegistry creates a new registry
func NewRegistry() Registry {
	return &registry{
		events: map[string]Event{},
	}
}

// NewBuilder creates a new builder
func NewBuilder() Builder {
	return &builder{}
}

// Registry represents an event registry
type Registry interface {
	Listen(event Event) error
	ListenAll(events []Event) error

	Unlisten(keyname string) error
	UnlistenAll(keynames []string) error
	Clear()

	Trigger(ctx contexts.Context) (results.Result, error)
}

// Builder represents an event builder.
type Builder interface {
	Create() Builder
	WithKeyname(keyname string) Builder
	WithBrowserType(browserType BrowserType) Builder
	WithEventName(eventName string) Builder
	WithAction(action ActionFn) Builder
	Now() (Event, error)
}

// Event represents an event
type Event interface {
	Keyname() string
	Action() ActionFn
}

// BrowserEvent represents a browser DOM event.
type BrowserEvent interface {
	Event
	BrowserType() BrowserType
}

// CustomEvent represents an application-defined event.
type CustomEvent interface {
	Event
	EventName() string
}
