package trees

import (
	"errors"
	"strings"
)

type groupBuilder struct {
	keyname   string
	resources []Resource
}

func (b *groupBuilder) Create() GroupBuilder {
	return &groupBuilder{}
}

func (b *groupBuilder) WithKeyname(keyname string) GroupBuilder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}

func (b *groupBuilder) AddResource(resource Resource) GroupBuilder {
	if resource != nil {
		b.resources = append(b.resources, resource)
	}
	return b
}

func (b *groupBuilder) WithResources(resources []Resource) GroupBuilder {
	b.resources = resources
	return b
}

func (b *groupBuilder) Now() (Group, error) {
	if b.keyname == "" {
		return nil, errors.New("group keyname is required")
	}

	seenLocales := map[string]bool{}
	for _, resource := range b.resources {
		if resource == nil {
			continue
		}
		if seenLocales[resource.Locale()] {
			return nil, errors.New("group resource locale already exists")
		}
		seenLocales[resource.Locale()] = true
	}

	return &group{keyname: b.keyname, resources: b.resources}, nil
}
