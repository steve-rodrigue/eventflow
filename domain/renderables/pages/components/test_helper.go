package components

import (
	"errors"
	"testing"

	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/events/contexts"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

func mustTemplate(t *testing.T, keyname string, code string) templates.Template {
	t.Helper()

	template, err := templates.NewMustacheBuilder().
		Create().
		WithKeyname(keyname).
		WithCode(code).
		Now()

	if err != nil {
		t.Fatalf("expected template, got error: %v", err)
	}

	return template
}

func mustComponent(t *testing.T, keyname string, template templates.Template, children ...Component) Component {
	t.Helper()

	component, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname(keyname).
		WithTemplate(template).
		WithChildren(children).
		Now()

	if err != nil {
		t.Fatalf("expected component, got error: %v", err)
	}

	return component
}

func mustEvent(t *testing.T, keyname string) events.Event {
	t.Helper()

	event, err := events.NewBuilder().
		Create().
		WithKeyname(keyname).
		WithEventName(keyname).
		WithAction(func(ctx contexts.Context) (results.Result, error) {
			return nil, errors.New("not implemented")
		}).
		Now()

	if err != nil {
		t.Fatalf("expected event, got error: %v", err)
	}

	return event
}
