package trees

import "slices"

type node struct {
	targets   []Target
	fallbacks []Fallback
}

func (n *node) HasTargets() bool {
	return len(n.targets) > 0
}

func (n *node) Targets() []Target {
	return slices.Clone(n.targets)
}

func (n *node) HasFallbacks() bool {
	return len(n.fallbacks) > 0
}

func (n *node) Fallbacks() []Fallback {
	return slices.Clone(n.fallbacks)
}

func (n *node) Resource(locale string, targetKeyname string, groupKeyname string) (Resource, bool) {
	for _, target := range n.targets {
		if target.Keyname() != targetKeyname {
			continue
		}

		group, ok := target.Group(groupKeyname)
		if !ok {
			return nil, false
		}

		return group.Resource(locale)
	}

	return nil, false
}
