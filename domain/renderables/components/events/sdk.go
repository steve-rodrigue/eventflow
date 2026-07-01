package events

// NativeType represents a DOM event supported by the browser.
type NativeType string

const (
	NativeTypeClick       NativeType = "click"
	NativeTypeDoubleClick NativeType = "dblclick"

	NativeTypeMouseDown  NativeType = "mousedown"
	NativeTypeMouseUp    NativeType = "mouseup"
	NativeTypeMouseEnter NativeType = "mouseenter"
	NativeTypeMouseLeave NativeType = "mouseleave"

	NativeTypeInput  NativeType = "input"
	NativeTypeChange NativeType = "change"
	NativeTypeSubmit NativeType = "submit"

	NativeTypeFocus NativeType = "focus"
	NativeTypeBlur  NativeType = "blur"

	NativeTypeKeyDown NativeType = "keydown"
	NativeTypeKeyUp   NativeType = "keyup"

	NativeTypeScroll NativeType = "scroll"
	NativeTypeLoad   NativeType = "load"
)

// OperationType represents the kind of UI update to apply in the browser.
type OperationType string

const (
	OperationTypeReplace  OperationType = "replace"
	OperationTypeAppend   OperationType = "append"
	OperationTypePrepend  OperationType = "prepend"
	OperationTypeRemove   OperationType = "remove"
	OperationTypeNavigate OperationType = "navigate"
	OperationTypeEmit     OperationType = "emit"
)

// Builder represents an event builder.
type Builder interface {
	Create() Builder
	WithKeyname(keyname string) Builder
	WithNativeType(nativeType NativeType) Builder
	WithEventName(eventName string) Builder
	WithAction(action ActionFn) Builder
	Now() (Event, error)
}

// Event represents an executable event.
type Event interface {
	Keyname() string
	Action() ActionFn
}

// NativeEvent represents a browser DOM event.
type NativeEvent interface {
	Event
	NativeType() NativeType
}

// CustomEvent represents an application-defined event.
type CustomEvent interface {
	Event
	EventName() string
}

// ActionFn is executed when an event is triggered.
type ActionFn func(ctx Context) (ActionResult, error)

// ContextBuilder represents a context builder.
type ContextBuilder interface {
	Create() ContextBuilder
	WithEvent(event Event) ContextBuilder
	WithPayload(payload map[string]any) ContextBuilder
	Now() (Context, error)
}

// Context represents the execution context of an event.
type Context interface {
	Event() Event
	Payload() map[string]any
	Value(key string) any
}

// ActionResultBuilder represents an action result builder.
type ActionResultBuilder interface {
	Create() ActionResultBuilder
	AddOperation(operation Operation) ActionResultBuilder
	WithOperations(operations []Operation) ActionResultBuilder
	Now() ActionResult
}

// ActionResult represents the UI changes sent back to the browser.
type ActionResult interface {
	HasOperations() bool
	Operations() []Operation
}

// OperationBuilder represents an operation builder
type OperationBuilder interface {
	Create() OperationBuilder
	WithType(operationType OperationType) OperationBuilder
	WithTarget(target string) OperationBuilder
	WithPayload(payload map[string]any) OperationBuilder
	Now() (Operation, error)
}

// Operation represents one browser-side UI operation.
type Operation interface {
	Type() OperationType
	Target() string
	Payload() map[string]any
}
