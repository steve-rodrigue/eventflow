package templates

import (
	"strings"
)

type mustacheRenderer struct{}

func (r *mustacheRenderer) Render(template Template, values map[string]string) string {
	code := template.Code()

	for _, param := range template.Params() {
		value := values[param]
		code = strings.ReplaceAll(code, "{"+param+"}", value)
	}

	return code
}
