package results

import "errors"

type operationBuilder struct {
	operationType OperationType
	action        Action
}

func (b *operationBuilder) Create() OperationBuilder {
	return &operationBuilder{}
}

func (b *operationBuilder) WithType(operationType OperationType) OperationBuilder {
	b.operationType = operationType
	return b
}

func (b *operationBuilder) WithAction(action Action) OperationBuilder {
	b.action = action
	return b
}

func (b *operationBuilder) Now() (Operation, error) {
	if b.operationType == "" {
		return nil, errors.New("operation type is required")
	}

	if b.action == nil {
		return nil, errors.New("operation action is required")
	}

	return &operation{
		operationType: b.operationType,
		action:        b.action,
	}, nil
}
