package starportserve

import (
	"context"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

type launchpadPlugin struct {
	app App
}

func newLaunchpadPlugin(app App) *launchpadPlugin {
	return &launchpadPlugin{
		app: app,
	}
}

func (p *launchpadPlugin) Name() string {
	return "Launchpad"
}

func (p *launchpadPlugin) Migrate(ctx context.Context) error {
	// migrate:
	//	appcli rest-server with --unsafe-cors (available only since v0.39.1).
	return cmdrunner.
		New(
			cmdrunner.DefaultWorkdir(p.app.Path),
		).
		Run(ctx,
			step.New(
				step.Exec(
					"go",
					"mod",
					"edit",
					"-require=github.com/cosmos/cosmos-sdk@v0.39.1",
				),
			),
		)
}

func (p *launchpadPlugin) Install(ctx context.Context, ldflags string) []step.Option {
	return []step.Option{
		step.Exec(
			"go",
			"install",
			"-mod", "readonly",
			"-ldflags", ldflags,
			filepath.Join(p.app.root(), "cmd", p.app.d()),
		),
		step.Exec(
			"go",
			"install",
			"-mod", "readonly",
			"-ldflags", ldflags,
			filepath.Join(p.app.root(), "cmd", p.app.cli()),
		),
	}
}
