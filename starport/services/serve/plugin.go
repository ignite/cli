package starportserve

import (
	"context"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/gomodule"
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
	InstallCommands(ldflags string) []step.Option

	// AddUserCommand returns step.Exec configuration to add users.
	AddUserCommand(accountName string) step.Option

	// ShowAccountCommand returns step.Exec configuration to run show account.
	ShowAccountCommand(accountName string) step.Option

	// ConfigCommands returns step.Exec configuration for config commands.
	ConfigCommands() []step.Option

	// GentxCommand returns step.Exec configuration for gentx command.
	GentxCommand(starportconf.Config) step.Option

	// PostInit hook.
	PostInit() error

	// StartCommands returns step.Exec configuration to start servers.
	StartCommands() [][]step.Option

	// StoragePaths returns a list of where persistent data kept.
	StoragePaths() []string

	// GenesisPath returns path of genesis.json.
	GenesisPath() string
}

type CosmosMajorVersion int

const (
	Launchpad CosmosMajorVersion = iota
	Stargate
)

const (
	tendermintPath                = "github.com/tendermint/tendermint"
	cosmosStargateTendermintMajor = "v0.34.0"
)

// detectCosmos dedects major version of Cosmos.
func (s *starportServe) detectCosmos() (CosmosMajorVersion, error) {
	parsed, err := gomodule.ParseAt(s.app.Path)
	if err != nil {
		return 0, err
	}
	for _, r := range parsed.Require {
		v := r.Mod
		if v.Path == tendermintPath {
			if strings.HasPrefix(v.Version, cosmosStargateTendermintMajor) {
				return Stargate, nil
			}
			break
		}
	}
	return Launchpad, nil
}

func (s *starportServe) pickPlugin() (Plugin, error) {
	version, err := s.detectCosmos()
	if err != nil {
		return nil, err
	}
	switch version {
	case Launchpad:
		return newLaunchpadPlugin(s.app), nil
	case Stargate:
		return newStargatePlugin(s.app), nil
	}
	panic("unknown cosmos version")
}
