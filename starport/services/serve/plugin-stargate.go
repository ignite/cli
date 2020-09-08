package starportserve

import (
	"context"
	"fmt"
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

func (p *stargatePlugin) StoragePaths() []string {
	return []string{
		p.app.nd(),
	}
}

func (p *stargatePlugin) GenesisPath() string {
	return fmt.Sprintf("%s/config/genesis.json", p.app.nd())
}
