package chain

import (
	"context"

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

	// ConfigCommands returns step.Exec configuration for config commands.
	Configure(ctx context.Context, chainID string) error

	// GentxCommand returns step.Exec configuration for gentx command.
	Gentx(context.Context, Validator) (path string, err error)

	// PostInit hook.
	PostInit(starportconf.Config) error

	// StartCommands returns step.Exec configuration to start servers.
	Start(context.Context, starportconf.Config) error

	// StoragePaths returns a list of where persistent data kept.
	StoragePaths() []string

	// Home returns the blockchain node's home dir.
	Home() string

	// Version of the plugin.
	Version() cosmosver.MajorVersion

	// SupportsIBC reports if app support IBC.
	SupportsIBC() bool
}

func (c *Chain) pickPlugin() (Plugin, error) {
	version := c.app.Version
	if version == "" {
		var err error
		version, err = cosmosver.Detect(c.app.Path)
		if err != nil {
			return nil, err
		}
	}
	switch version {
	case cosmosver.Launchpad:
		return newLaunchpadPlugin(c.app, c), nil
	case cosmosver.Stargate:
		return newStargatePlugin(c.app, c), nil
	}
	panic("unknown cosmos version")
}
