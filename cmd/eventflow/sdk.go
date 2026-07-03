package eventflow

import (
	"log"

	"github.com/steve-rodrigue/eventflow/applications"
	infrahttps "github.com/steve-rodrigue/eventflow/infrastructure/https"
)

// New creates a new eventflow command
func New(
	address string,
	assetsBasePath string,
	tree applications.Tree,
	logger *log.Logger,
) (*Command, error) {
	if logger == nil {
		logger = log.Default()
	}

	app := applications.NewDefaultApplication(assetsBasePath)

	if err := app.Initialize(tree); err != nil {
		return nil, err
	}

	command := &Command{
		app:            app,
		address:        address,
		assetsBasePath: assetsBasePath,
		logger:         logger,
	}

	handlers, err := command.handlers()
	if err != nil {
		return nil, err
	}

	command.server = infrahttps.NewServer(address, handlers)

	return command, nil
}
