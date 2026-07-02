package contexts

import (
	"errors"
	"strings"
)

type builder struct {
	eventName string
	payload   map[string]any
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) WithEventName(eventName string) Builder {
	b.eventName = strings.TrimSpace(eventName)
	return b
}

func (b *builder) WithPayload(payload map[string]any) Builder {
	b.payload = payload
	return b

}

func (b *builder) Now() (Context, error) {

	if b.eventName == "" {
		return nil, errors.New("context event name is required")
	}
	if b.payload == nil {
		b.payload = map[string]any{}
	}
	return &context{
		eventName: b.eventName,
		payload:   b.payload,
	}, nil

}
