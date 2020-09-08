package starportserve

import (
	"context"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

type stargatePlugin struct {
	app App
}

func newStargatePlugin(app App) *stargatePlugin {
	return &stargatePlugin{
		app: app,
	}
}

func (p *stargatePlugin) Name() string {
	return "Stargate"
}

func (p *stargatePlugin) Migrate(ctx context.Context) error {
	return nil
}

func (p *stargatePlugin) Install(ctx context.Context, ldflags string) []step.Option {
	return []step.Option{
		step.Exec(
			"go",
			"install",
			"-mod", "readonly",
			"-ldflags", ldflags,
			filepath.Join(p.app.root(), "cmd", p.app.d()),
		),
	}
}
