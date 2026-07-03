package trees

import (
	"errors"
)

type nodeBuilder struct {
	targets   []Target
	fallbacks []Fallback
}

func (b *nodeBuilder) Create() NodeBuilder {
	return &nodeBuilder{}
}

func (b *nodeBuilder) AddTarget(target Target) NodeBuilder {
	if target != nil {
		b.targets = append(b.targets, target)
	}
	return b
}

func (b *nodeBuilder) WithTargets(targets []Target) NodeBuilder {
	b.targets = targets
	return b
}

func (b *nodeBuilder) AddFallback(fallback Fallback) NodeBuilder {
	if fallback != nil {
		b.fallbacks = append(b.fallbacks, fallback)
	}
	return b
}

func (b *nodeBuilder) WithFallbacks(fallbacks []Fallback) NodeBuilder {
	b.fallbacks = fallbacks
	return b
}

func (b *nodeBuilder) Now() (Node, error) {
	seenTargets := map[string]bool{}
	for _, target := range b.targets {
		if target == nil {
			continue
		}
		if seenTargets[target.Keyname()] {
			return nil, errors.New("node target already exists")
		}
		seenTargets[target.Keyname()] = true
	}

	return &node{targets: b.targets, fallbacks: b.fallbacks}, nil
}
