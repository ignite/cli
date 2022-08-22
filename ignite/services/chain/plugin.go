package chain

import (
	"context"

	"github.com/ignite/cli/ignite/chainconfig"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
)

// TODO omit -cli log messages for Stargate.

type Plugin interface {
	// Name of a Cosmos version.
	Name() string

	// Gentx returns step.Exec configuration for gentx command.
	Gentx(context.Context, chaincmdrunner.Runner, Validator) (path string, err error)

	// Configure configures config defaults.
	Configure(string, chainconfig.Config) error

	// Start returns step.Exec configuration to start servers.
	Start(context.Context, chaincmdrunner.Runner, chainconfig.Config) error

	// Home returns the blockchain node's home dir.
	Home() string
}

func (c *Chain) pickPlugin() Plugin {
	return newStargatePlugin(c.app)
}
