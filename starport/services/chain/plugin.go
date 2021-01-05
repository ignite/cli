package chain

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	starportconf "github.com/tendermint/starport/starport/services/chain/conf"
)

// TODO omit -cli log messages for Stargate.

type Plugin interface {
	// Name of a Cosmos version.
	Name() string

	// Setup performs the initial setup for plugin.
	Setup(context.Context) error

	// Binaries returns a list of binaries that will be compiled for the app.
	Binaries() []string

	// AddUserCommand returns step.Exec configuration to add users.
	AddUserCommand(cmd chaincmd.ChainCmd, name string) step.Options

	// ImportUserCommand returns step.Exec configuration to import users.
	ImportUserCommand(cmd chaincmd.ChainCmd, name, mnemonic string) step.Options

	// ShowAccountCommand returns step.Exec configuration to run show account.
	ShowAccountCommand(cmd chaincmd.ChainCmd, accountName string) step.Option

	// ConfigCommands returns step.Exec configuration for config commands.
	ConfigCommands(cmd chaincmd.ChainCmd, chainID string) []step.Option

	// GentxCommand returns step.Exec configuration for gentx command.
	GentxCommand(cmd chaincmd.ChainCmd, v Validator) step.Option

	// StartCommands returns step.Exec configuration to start servers.
	StartCommands(cmd chaincmd.ChainCmd, config starportconf.Config) [][]step.Option

	// PostInit hook.
	PostInit(starportconf.Config) error

	// StoragePaths returns a list of where persistent data kept.
	StoragePaths() []string

	// Home returns the blockchain node's home dir.
	Home() string

	// CLIHome returns the cli blockchain node's home dir.
	CLIHome() string

	// Version of the plugin.
	Version() cosmosver.MajorVersion

	// SupportsIBC reports if app support IBC.
	SupportsIBC() bool
}

func (c *Chain) pickPlugin() (Plugin, error) {
	version, err := c.CosmosVersion()
	if err != nil {
		return nil, err
	}
	switch version {
	case cosmosver.Launchpad:
		return newLaunchpadPlugin(c.app), nil
	case cosmosver.Stargate:
		return newStargatePlugin(c.app, c), nil
	}
	panic("unknown cosmos version")
}
