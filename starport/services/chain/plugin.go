package chain

import (
	"context"

	starportconf "github.com/tendermint/starport/starport/chainconf"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
)

// TODO omit -cli log messages for Stargate.

type Plugin interface {
	// Name of a Cosmos version.
	Name() string

	// Setup performs the initial setup for plugin.
	Setup(context.Context) error

	// ConfigCommands returns step.Exec configuration for config commands.
	Configure(context.Context, chaincmdrunner.Runner, string) error

	// GentxCommand returns step.Exec configuration for gentx command.
	Gentx(context.Context, chaincmdrunner.Runner, Validator) (path string, err error)

	// PostInit hook.
	PostInit(string, starportconf.Config) error

	// StartCommands returns step.Exec configuration to start servers.
	Start(context.Context, chaincmdrunner.Runner, starportconf.Config) error

	// Home returns the blockchain node's home dir.
	Home() string
}

func (c *Chain) pickPlugin() Plugin {
	return newStargatePlugin(c.app)
}
