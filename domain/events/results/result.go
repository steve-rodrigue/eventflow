package results

type result struct {
	operations []Operation
}

func (r *result) Operations() []Operation {
	return append([]Operation(nil), r.operations...)
}
