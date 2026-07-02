package pages

type header struct {
	name  string
	value string
}

func (h *header) Name() string {
	return h.name
}

func (h *header) Value() string {
	return h.value
}
