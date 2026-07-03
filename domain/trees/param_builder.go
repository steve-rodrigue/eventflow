package trees

import (
	"errors"
	"strings"
)

type paramBuilder struct {
	keyname string
	value   string
}

func (b *paramBuilder) Create() ParamBuilder {
	return &paramBuilder{}
}

func (b *paramBuilder) WithKeyname(keyname string) ParamBuilder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}
func (b *paramBuilder) WithValue(value string) ParamBuilder {
	b.value = strings.TrimSpace(value)
	return b
}

func (b *paramBuilder) Now() (Param, error) {
	if b.keyname == "" {
		return nil, errors.New("param keyname is required")
	}
	return &param{keyname: b.keyname, value: b.value}, nil
}
