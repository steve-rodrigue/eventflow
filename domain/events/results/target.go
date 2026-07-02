package results

type target struct {
	element string
	html    string
}

func (t *target) Element() string {
	return t.element
}

func (t *target) HTML() string {
	return t.html
}
