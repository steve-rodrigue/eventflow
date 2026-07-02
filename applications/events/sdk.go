package events

import (
	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/events/contexts"
)

// IncomingMessage represents an incoming message
type IncomingMessage struct {
	Type    string         `json:"type"`
	Event   string         `json:"event"`
	Payload map[string]any `json:"payload"`
}

// OutgoingOperation represents an outgoing operation
type OutgoingOperation struct {
	Type   string `json:"type"`
	Target string `json:"target,omitempty"`
	HTML   string `json:"html,omitempty"`
	URL    string `json:"url,omitempty"`
}

// OutgoingMessage represents an outgoing message
type OutgoingMessage struct {
	Type       string              `json:"type"`
	Operations []OutgoingOperation `json:"operations,omitempty"`
	Error      string              `json:"error,omitempty"`
}

// NewApplication creates a new application
func NewApplication(
	contextBuilder contexts.Builder,
	registry events.Registry,
) Application {
	return &application{
		contextBuilder: contextBuilder,
		registry:       registry,
	}
}

// Application represents the event application
type Application interface {
	Execute(msg IncomingMessage) (*OutgoingMessage, error)
}
