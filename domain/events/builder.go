package events

import (
	"errors"
	"strings"
)

type builder struct {
	keyname     string
	browserType BrowserType
	eventName   string
	action      ActionFn
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) WithKeyname(keyname string) Builder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}

func (b *builder) WithBrowserType(browserType BrowserType) Builder {
	b.browserType = browserType
	return b
}

func (b *builder) WithEventName(eventName string) Builder {
	b.eventName = strings.TrimSpace(eventName)
	return b
}

func (b *builder) WithAction(action ActionFn) Builder {
	b.action = action
	return b
}

func (b *builder) Now() (Event, error) {
	if b.keyname == "" {
		return nil, errors.New("event keyname is required")
	}

	if b.action == nil {
		return nil, errors.New("event action is required")
	}

	hasBrowserType := b.browserType != ""
	hasEventName := b.eventName != ""

	if hasBrowserType && hasEventName {
		return nil, errors.New("event cannot have both browser type and event name")
	}

	if !hasBrowserType && !hasEventName {
		return nil, errors.New("event requires browser type or event name")
	}

	if hasBrowserType {
		return &browserEvent{
			event: event{
				keyname: b.keyname,
				action:  b.action,
			},
			browserType: b.browserType,
		}, nil
	}

	return &customEvent{
		event: event{
			keyname: b.keyname,
			action:  b.action,
		},
		eventName: b.eventName,
	}, nil
}
