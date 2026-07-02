package results

type operation struct {
	operationType OperationType
	action        Action
}

func (o *operation) Type() OperationType {
	return o.operationType
}

func (o *operation) Action() Action {
	return o.action
}
