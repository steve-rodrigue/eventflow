package events

import (
	"errors"

	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/events/contexts"
	"github.com/steve-rodrigue/eventflow/domain/events/results"
)

type application struct {
	contextBuilder contexts.Builder
	registry       events.Registry
}

// Execute executes the message
func (app *application) Execute(msg IncomingMessage) (*OutgoingMessage, error) {
	if msg.Type != "event" {
		return nil, errors.New("the message type was expected to be event")
	}

	ctx, err := app.contextBuilder.Create().
		WithEventName(msg.Event).
		WithPayload(msg.Payload).
		Now()

	if err != nil {
		return nil, err
	}

	result, err := app.registry.Trigger(ctx)
	if err != nil {
		return nil, err
	}

	return &OutgoingMessage{
		Type:       "operations",
		Operations: app.toOutgoingOperations(result),
	}, nil
}

func (app *application) toOutgoingOperations(result results.Result) []OutgoingOperation {
	operations := make([]OutgoingOperation, 0, len(result.Operations()))

	for _, operation := range result.Operations() {
		out := OutgoingOperation{
			Type: string(operation.Type()),
		}

		action := operation.Action()

		if action.IsTarget() {
			target := action.Target()
			out.Target = target.Element()
			out.HTML = target.HTML()
		}

		if action.IsNavigate() {
			out.URL = action.Navigate().URL()
		}

		operations = append(operations, out)
	}

	return operations
}
