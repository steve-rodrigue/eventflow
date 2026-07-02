package templates

import (
	"slices"
	"time"
)

type template struct {
	keyname   string
	code      string
	params    []string
	createdAt time.Time
}

func (t *template) Keyname() string {
	return t.keyname
}

func (t *template) Code() string {
	return t.code
}

func (t *template) Contains(name string) bool {
	return slices.Contains(t.params, name)
}

func (t *template) HasParams() bool {
	return len(t.params) > 0
}

func (t *template) Params() []string {
	return slices.Clone(t.params)
}
