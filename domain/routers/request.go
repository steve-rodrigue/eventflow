package routers

import (
	"slices"

	"github.com/steve-rodrigue/eventflow/domain/trees"
)

type request struct {
	targetKeyname string
	groupKeyname  string
	locale        string
	params        []trees.Param
}

func (r *request) TargetKeyname() string {
	return r.targetKeyname
}

func (r *request) GroupKeyname() string {
	return r.groupKeyname
}

func (r *request) Locale() string {
	return r.locale
}

func (r *request) HasParams() bool {
	return len(r.params) > 0
}

func (r *request) Params() []trees.Param {
	return slices.Clone(r.params)
}
