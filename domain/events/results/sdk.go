package results

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

// NewBuilder creates a new result builder
func NewBuilder() Builder {
	return &builder{}
}

// NewOperationBuilder creates a new operation builder
func NewOperationBuilder() OperationBuilder {
	return &operationBuilder{}
}

// NewActionBuilder creates a new action builder
func NewActionBuilder() ActionBuilder {
	return &actionBuilder{}
}

// NewTargetBuilder creates a new target builder
func NewTargetBuilder() TargetBuilder {
	return &targetBuilder{}
}

// NewNavigateBuilder creates a new navigate builder
func NewNavigateBuilder() NavigateBuilder {
	return &navigateBuilder{}
}

// Builder represents a result builder.
type Builder interface {
	Create() Builder
	AddOperation(operation Operation) Builder
	WithOperations(operations []Operation) Builder
	Now() (Result, error)
}

// Result represents an event result
type Result interface {
	Operations() []Operation
}

// OperationBuilder represents an operation builder.
type OperationBuilder interface {
	Create() OperationBuilder
	WithType(operationType OperationType) OperationBuilder
	WithAction(action Action) OperationBuilder
	Now() (Operation, error)
}

// Operation represents an operation
type Operation interface {
	Type() OperationType
	Action() Action
}

// ActionBuilder represents an action builder.
type ActionBuilder interface {
	Create() ActionBuilder
	WithNavigate(navigate Navigate) ActionBuilder
	WithTarget(target Target) ActionBuilder
	Now() (Action, error)
}

// Action represents an action
type Action interface {
	IsNavigate() bool
	Navigate() Navigate
	IsTarget() bool
	Target() Target
}

// TargetBuilder represents a target action builder.
type TargetBuilder interface {
	Create() TargetBuilder
	WithElement(element string) TargetBuilder
	WithHTML(html string) TargetBuilder
	Now() (Target, error)
}

// Target represents a target
type Target interface {
	Element() string
	HTML() string
}

// NavigateBuilder represents a navigation action builder.
type NavigateBuilder interface {
	Create() NavigateBuilder
	WithURL(url string) NavigateBuilder
	Now() (Navigate, error)
}

// Navigate represents a navigation action
type Navigate interface {
	URL() string
}
