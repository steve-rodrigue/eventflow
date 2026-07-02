package contexts

// NewBuilder creates a new context builder.
func NewBuilder() Builder {
	return &builder{}
}

// Builder represents a context builder
type Builder interface {
	Create() Builder
	WithEventName(eventName string) Builder
	WithPayload(payload map[string]any) Builder
	Now() (Context, error)
}

// Context represents a context
type Context interface {
	EventName() string
	Payload() map[string]any
	HasValue(key string) bool
	Value(key string) any
}
