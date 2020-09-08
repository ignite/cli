package starportserve

import (
	"context"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/gomodule"
)

type Plugin interface {
	// Name of a Cosmos version.
	Name() string

	// Migrate migrates apps generated with older(minor) version of Starport,
	// to the current version to make them compatible with the updated `serve` command.
	Migrate(context.Context) error

	// Install returns the installation's step.Exec configuration.
	Install(ctx context.Context, ldflags string) []step.Option
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
