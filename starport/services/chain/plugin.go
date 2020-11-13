package chain

import (
	"context"

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

	// InstallCommands returns step.Exec configurations to install app.
	InstallCommands(ldflags string) (options []step.Option, binaries []string)

	// AddUserCommand returns step.Exec configuration to add users.
	AddUserCommand(name string) step.Options

	// ImportUserCommand returns step.Exec configuration to import users.
	ImportUserCommand(namem, mnemonic string) step.Options

	// ShowAccountCommand returns step.Exec configuration to run show account.
	ShowAccountCommand(accountName string) step.Option

	// ConfigCommands returns step.Exec configuration for config commands.
	ConfigCommands(chainID string) []step.Option

	// GentxCommand returns step.Exec configuration for gentx command.
	GentxCommand(chainID string, c starportconf.Config) step.Option

	// PostInit hook.
	PostInit(starportconf.Config) error

	// StartCommands returns step.Exec configuration to start servers.
	StartCommands(starportconf.Config) [][]step.Option

	// StoragePaths returns a list of where persistent data kept.
	StoragePaths() []string

	// Home returns the root config dir's path of app.
	Home() (string, error)

	// Version of the plugin.
	Version() cosmosver.MajorVersion

	// SupportsIBC reports if app support IBC.
	SupportsIBC() bool
}

func (s *Chain) pickPlugin() (Plugin, error) {
	version, err := cosmosver.Detect(s.app.Path)
	if err != nil {
		return nil, err
	}
	switch version {
	case cosmosver.Launchpad:
		return newLaunchpadPlugin(s.app), nil
	case cosmosver.Stargate:
		return newStargatePlugin(s.app, s), nil
	}
	panic("unknown cosmos version")
}
