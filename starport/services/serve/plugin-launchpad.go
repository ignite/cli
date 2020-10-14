package starportserve

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xurl"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
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

func (p *launchpadPlugin) InstallCommands(ldflags string) (options []step.Option, binaries []string) {
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
		}, []string{
			p.app.d(),
			p.app.cli(),
		}
}

func (p *launchpadPlugin) AddUserCommand(accountName string) step.Option {
	return step.Exec(
		p.app.cli(),
		"keys",
		"add",
		accountName,
		"--output", "json",
		"--keyring-backend", "test",
	)
}

func (p *launchpadPlugin) ShowAccountCommand(accountName string) step.Option {
	return step.Exec(
		p.app.cli(),
		"keys",
		"show",
		accountName,
		"-a",
		"--keyring-backend", "test",
	)
}

func (p *launchpadPlugin) ConfigCommands() []step.Option {
	return []step.Option{
		step.Exec(
			p.app.cli(),
			"config",
			"keyring-backend",
			"test",
		),
		step.Exec(
			p.app.cli(),
			"config",
			"chain-id",
			p.app.n(),
		),
		step.Exec(
			p.app.cli(),
			"config",
			"output",
			"json",
		),
		step.Exec(
			p.app.cli(),
			"config",
			"indent",
			"true",
		),
		step.Exec(
			p.app.cli(),
			"config",
			"trust-node",
			"true",
		),
	}
}

func (p *launchpadPlugin) GentxCommand(conf starportconf.Config) step.Option {
	return step.Exec(
		p.app.d(),
		"gentx",
		"--name", conf.Validator.Name,
		"--keyring-backend", "test",
		"--amount", conf.Validator.Staked,
	)
}

func (p *launchpadPlugin) PostInit(conf starportconf.Config) error {
	return p.configtoml(conf)
}

func (p *launchpadPlugin) configtoml(conf starportconf.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, "."+p.app.nd(), "config/config.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("rpc.laddr", xurl.TCP(conf.Servers.RPCAddr))
	config.Set("p2p.laddr", xurl.TCP(conf.Servers.P2PAddr))
	config.Set("rpc.prof_laddr", conf.Servers.ProfAddr)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *launchpadPlugin) StartCommands(conf starportconf.Config) [][]step.Option {
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
		step.NewOptions().
			Add(
				step.Exec(
					p.app.cli(),
					"rest-server",
					"--unsafe-cors",
					"--laddr", xurl.TCP(conf.Servers.APIAddr),
					"--node", xurl.TCP(conf.Servers.RPCAddr),
				),
				step.PostExec(func(exitErr error) error {
					return errors.Wrapf(exitErr, "cannot run %[1]vcli rest-server", p.app.Name)
				}),
			),
	}
}

func (p *launchpadPlugin) StoragePaths() []string {
	return []string{
		fmt.Sprintf(".%s", p.app.nd()),
		fmt.Sprintf(".%s", p.app.ncli()),
	}
}

func (p *launchpadPlugin) GenesisPath() string {
	return fmt.Sprintf(".%s/config/genesis.json", p.app.nd())
}

func (p *launchpadPlugin) Version() cosmosver.MajorVersion { return cosmosver.Launchpad }
