package trees

import (
	"errors"
	"strings"
)

type targetBuilder struct {
	keyname string
	groups  []Group
}

func (b *targetBuilder) Create() TargetBuilder {
	return &targetBuilder{}
}

func (b *targetBuilder) WithKeyname(keyname string) TargetBuilder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}

func (b *targetBuilder) AddGroup(group Group) TargetBuilder {
	if group != nil {
		b.groups = append(b.groups, group)
	}
	return b
}

func (b *targetBuilder) WithGroups(groups []Group) TargetBuilder {
	b.groups = groups
	return b
}

func (b *targetBuilder) Now() (Target, error) {
	if b.keyname == "" {
		return nil, errors.New("target keyname is required")
	}

	seenGroups := map[string]bool{}
	for _, group := range b.groups {
		if group == nil {
			continue
		}
		if seenGroups[group.Keyname()] {
			return nil, errors.New("target group already exists")
		}
		seenGroups[group.Keyname()] = true
	}

	return &target{keyname: b.keyname, groups: b.groups}, nil
}
