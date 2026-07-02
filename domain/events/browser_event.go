package events

type browserEvent struct {
	event

	browserType BrowserType
}

func (e *browserEvent) BrowserType() BrowserType {
	return e.browserType
}
