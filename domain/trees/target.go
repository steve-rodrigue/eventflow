package trees

import "slices"

type target struct {
	keyname string
	groups  []Group
}

func (t *target) Keyname() string {
	return t.keyname
}

func (t *target) HasGroups() bool {
	return len(t.groups) > 0
}

func (t *target) Groups() []Group {
	return slices.Clone(t.groups)
}

func (t *target) Group(keyname string) (Group, bool) {
	for _, group := range t.groups {
		if group.Keyname() == keyname {
			return group, true
		}
	}
	return nil, false
}
