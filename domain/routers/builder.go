package routers

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/trees"
)

type builder struct {
	targetKeyname string
	groupKeyname  string
	locale        string
	params        []trees.Param
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) WithTargetKeyname(targetKeyname string) Builder {
	b.targetKeyname = strings.TrimSpace(targetKeyname)
	return b
}

func (b *builder) WithGroupKeyname(groupKeyname string) Builder {
	b.groupKeyname = strings.TrimSpace(groupKeyname)
	return b
}

func (b *builder) WithLocale(locale string) Builder {
	b.locale = strings.TrimSpace(locale)
	return b
}

func (b *builder) AddParam(param trees.Param) Builder {
	if param != nil {
		b.params = append(b.params, param)
	}

	return b
}

func (b *builder) WithParams(params []trees.Param) Builder {
	b.params = params
	return b
}

func (b *builder) Now() (Request, error) {
	if b.targetKeyname == "" {
		return nil, errors.New("route request target keyname is required")
	}

	if b.groupKeyname == "" {
		return nil, errors.New("route request group keyname is required")
	}

	if b.locale == "" {
		return nil, errors.New("route request locale is required")
	}

	return &request{
		targetKeyname: b.targetKeyname,
		groupKeyname:  b.groupKeyname,
		locale:        b.locale,
		params:        b.params,
	}, nil
}
