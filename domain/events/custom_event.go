package events

type customEvent struct {
	event

	eventName string
}

func (e *customEvent) EventName() string {
	return e.eventName
}
