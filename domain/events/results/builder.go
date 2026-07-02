package results

type builder struct {
	operations []Operation
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) AddOperation(operation Operation) Builder {
	if operation != nil {
		b.operations = append(b.operations, operation)
	}

	return b
}

func (b *builder) WithOperations(operations []Operation) Builder {
	b.operations = operations
	return b
}

func (b *builder) Now() (Result, error) {
	return &result{
		operations: b.operations,
	}, nil
}
