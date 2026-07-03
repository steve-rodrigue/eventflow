package trees

import "slices"

type route struct {
	pattern string
	params  []Param
}

func (r *route) Pattern() string {
	return r.pattern
}

func (r *route) HasParams() bool {
	return len(r.params) > 0
}

func (r *route) Params() []Param {
	return slices.Clone(r.params)
}
