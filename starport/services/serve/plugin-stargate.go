package starportserve

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
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

func (p *stargatePlugin) InstallCommands(ldflags string) []step.Option {
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

func (p *stargatePlugin) AddUserCommand(accountName string) step.Option {
	return step.Exec(
		p.app.d(),
		"keys",
		"add",
		accountName,
		"--output", "json",
		"--keyring-backend", "test",
	)
}

func (p *stargatePlugin) ShowAccountCommand(accountName string) step.Option {
	return step.Exec(
		p.app.d(),
		"keys",
		"show",
		accountName,
		"-a",
		"--keyring-backend", "test",
	)
}

func (p *stargatePlugin) ConfigCommands() []step.Option {
	return nil
}

func (p *stargatePlugin) GentxCommand(conf starportconf.Config) step.Option {
	return step.Exec(
		p.app.d(),
		"gentx", conf.Validator.Name,
		"--chain-id", p.app.nd(),
		"--keyring-backend", "test",
		"--amount", conf.Validator.Staked,
	)
}

func (p *stargatePlugin) PostInit() error {
	// TODO find a better way in order to not delete comments in the toml.yml
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, p.app.nd(), "config/app.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	config.Set("api.enable", true)
	config.Set("api.enabled-unsafe-cors", true)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) StartCommands() [][]step.Option {
	return [][]step.Option{
		step.NewOptions().
			Add(
				step.Exec(
					p.app.d(),
					"start",
				),
				step.PostExec(func(exitErr error) error {
					return errors.Wrapf(exitErr, "cannot run %[1]vd start", p.app.Name)
				}),
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
