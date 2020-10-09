package starportserve

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
)

// TODO omit -cli log messages for Stargate.

type Plugin interface {
	// Name of a Cosmos version.
	Name() string

	// Migrate migrates apps generated with older(minor) version of Starport,
	// to the current version to make them compatible with the updated `serve` command.
	Migrate(context.Context) error

	// InstallCommands returns step.Exec configurations to install app.
	InstallCommands(ldflags string) (options []step.Option, binaries []string)

	// AddUserCommand returns step.Exec configuration to add users.
	AddUserCommand(accountName string) step.Option

	// ShowAccountCommand returns step.Exec configuration to run show account.
	ShowAccountCommand(accountName string) step.Option

	// ConfigCommands returns step.Exec configuration for config commands.
	ConfigCommands() []step.Option

	// GentxCommand returns step.Exec configuration for gentx command.
	GentxCommand(starportconf.Config) step.Option

	// PostInit hook.
	PostInit(starportconf.Config) error

	// StartCommands returns step.Exec configuration to start servers.
	StartCommands(starportconf.Config) [][]step.Option

	// StoragePaths returns a list of where persistent data kept.
	StoragePaths() []string

	// GenesisPath returns path of genesis.json.
	GenesisPath() string

	// Version of the plugin.
	Version() cosmosver.MajorVersion
}

func (s *Serve) pickPlugin() (Plugin, error) {
	version, err := cosmosver.Detect(s.app.Path)
	if err != nil {
		return nil, err
	}
	switch version {
	case cosmosver.Launchpad:
		return newLaunchpadPlugin(s.app), nil
	case cosmosver.Stargate:
		return newStargatePlugin(s.app), nil
	}
	panic("unknown cosmos version")
}
