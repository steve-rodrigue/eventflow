package events

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/events/contexts"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
)

type registry struct {
	events map[string]Event
}

func (r *registry) Listen(event Event) error {
	if event == nil {
		return errors.New("event is required")
	}

	if event.Keyname() == "" {
		return errors.New("event keyname is required")
	}

	if event.Action() == nil {
		return errors.New("event action is required")
	}

	if _, exists := r.events[event.Keyname()]; exists {
		return errors.New("event already exists")
	}

	r.events[event.Keyname()] = event

	return nil
}

func (r *registry) ListenAll(events []Event) error {
	for _, event := range events {
		if err := r.Listen(event); err != nil {
			return err
		}
	}

	return nil
}

func (r *registry) Unlisten(keyname string) error {
	keyname = strings.TrimSpace(keyname)

	if keyname == "" {
		return errors.New("event keyname is required")
	}

	if _, exists := r.events[keyname]; !exists {
		return errors.New("event does not exist")
	}

	delete(r.events, keyname)

	return nil
}

func (r *registry) UnlistenAll(keynames []string) error {
	for _, keyname := range keynames {
		if err := r.Unlisten(keyname); err != nil {
			return err
		}
	}

	return nil
}

func (r *registry) Clear() {
	r.events = map[string]Event{}
}

func (r *registry) Trigger(ctx contexts.Context) (results.Result, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}

	keyname := strings.TrimSpace(ctx.EventName())

	if keyname == "" {
		return nil, errors.New("event name is required")
	}

	event, exists := r.events[keyname]

	if !exists {
		return nil, errors.New("event does not exist")
	}

	return event.Action()(ctx)
}
