package chain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xurl"
	starportconf "github.com/tendermint/starport/starport/services/chain/conf"
)

type stargatePlugin struct {
	app   App
	chain *Chain
}

func newStargatePlugin(app App, chain *Chain) *stargatePlugin {
	return &stargatePlugin{
		app:   app,
		chain: chain,
	}
}

func (p *stargatePlugin) Name() string {
	return "Stargate"
}

func (p *stargatePlugin) Setup(ctx context.Context) error {
	return nil
}

func (p *stargatePlugin) InstallCommands(ldflags string) (options []step.Option, binaries []string) {
	return []step.Option{
			step.Exec(
				"go",
				"install",
				"-mod", "readonly",
				"-ldflags", ldflags,
				filepath.Join(p.app.root(), "cmd", p.app.d()),
			),
		}, []string{
			p.app.d(),
		}
}

func (p *stargatePlugin) AddUserCommand(accountName string) step.Options {
	return step.NewOptions().
		Add(
			step.Exec(
				p.app.d(),
				"keys",
				"add",
				accountName,
				"--output", "json",
				"--keyring-backend", "test",
			),
		)
}

func (p *stargatePlugin) ImportUserCommand(name, mnemonic string) step.Options {
	return step.NewOptions().
		Add(
			step.Exec(
				p.app.d(),
				"keys",
				"add",
				name,
				"--recover",
				"--keyring-backend", "test",
			),
			step.Write([]byte(mnemonic+"\n")),
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

func (p *stargatePlugin) ConfigCommands(_ string) []step.Option {
	return nil
}

func (p *stargatePlugin) GentxCommand(chainID string, conf starportconf.Config) step.Option {
	return step.Exec(
		p.app.d(),
		"gentx", conf.Validator.Name,
		"--chain-id", chainID,
		"--keyring-backend", "test",
		"--amount", conf.Validator.Staked,
	)
}

func (p *stargatePlugin) PostInit(conf starportconf.Config) error {
	if err := p.apptoml(conf); err != nil {
		return err
	}
	return p.configtoml(conf)
}

func (p *stargatePlugin) apptoml(conf starportconf.Config) error {
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
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("api.address", xurl.TCP(conf.Servers.APIAddr))
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) configtoml(conf starportconf.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, p.app.nd(), "config/config.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("consensus.timeout_commit", "1s")
	config.Set("consensus.timeout_propose", "1s")
	config.Set("rpc.laddr", xurl.TCP(conf.Servers.RPCAddr))
	config.Set("p2p.laddr", xurl.TCP(conf.Servers.P2PAddr))
	config.Set("rpc.pprof_laddr", conf.Servers.ProfAddr)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) StartCommands(conf starportconf.Config) [][]step.Option {
	return [][]step.Option{
		step.NewOptions().
			Add(
				step.Exec(
					p.app.d(),
					"start",
					"--pruning", "nothing",
					"--grpc.address", conf.Servers.GRPCAddr,
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

func (p *stargatePlugin) Home() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, p.app.nd()), nil
}

func (p *stargatePlugin) Version() cosmosver.MajorVersion { return cosmosver.Stargate }

func (p *stargatePlugin) SupportsIBC() bool { return true }
