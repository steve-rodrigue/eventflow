package trees

type request struct {
	path   string
	method string
	locale string
	target string
}

func (r *request) Path() string {
	return r.path
}

func (r *request) Method() string {
	return r.method
}

func (r *request) Locale() string {
	return r.locale
}

func (r *request) Target() string {
	return r.target
}
