package templates

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var mustacheParamPattern = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

type mustacheBuilder struct {
	keyname string
	code    string
}

func (b *mustacheBuilder) Create() Builder {
	return &mustacheBuilder{}
}

func (b *mustacheBuilder) WithKeyname(keyname string) Builder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}

func (b *mustacheBuilder) WithCode(code string) Builder {
	b.code = code
	return b
}

func (b *mustacheBuilder) Now() (Template, error) {
	if b.keyname == "" {
		return nil, errors.New("template keyname is required")
	}

	if b.code == "" {
		return nil, errors.New("template code is required")
	}

	return &template{
		keyname:   b.keyname,
		code:      b.code,
		params:    extractMustacheParams(b.code),
		createdAt: time.Now(),
	}, nil
}

func extractMustacheParams(code string) []string {
	matches := mustacheParamPattern.FindAllStringSubmatch(code, -1)

	params := make([]string, 0, len(matches))
	seen := map[string]bool{}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		name := match[1]

		if seen[name] {
			continue
		}

		seen[name] = true
		params = append(params, name)
	}

	return params
}
