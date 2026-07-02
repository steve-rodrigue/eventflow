package events

type event struct {
	keyname string
	action  ActionFn
}

func (e *event) Keyname() string {
	return e.keyname
}

func (e *event) Action() ActionFn {
	return e.action
}
