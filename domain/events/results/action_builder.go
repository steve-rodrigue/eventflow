package results

import "errors"

type actionBuilder struct {
	navigate Navigate
	target   Target
}

func (b *actionBuilder) Create() ActionBuilder {
	return &actionBuilder{}
}

func (b *actionBuilder) WithNavigate(navigate Navigate) ActionBuilder {
	b.navigate = navigate
	return b
}

func (b *actionBuilder) WithTarget(target Target) ActionBuilder {
	b.target = target
	return b
}

func (b *actionBuilder) Now() (Action, error) {
	if b.navigate == nil && b.target == nil {
		return nil, errors.New("action requires navigate or target")
	}

	if b.navigate != nil && b.target != nil {
		return nil, errors.New("action cannot have both navigate and target")
	}

	return &action{
		navigate: b.navigate,
		target:   b.target,
	}, nil
}
