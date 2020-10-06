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

var (
	rpcPort  = os.Getenv("RPC_PORT")
	p2pPort  = os.Getenv("P2P_PORT")
	profPort = os.Getenv("PROF_PORT")
	apiPort  = os.Getenv("API_PORT")
	grpcPort = os.Getenv("GRPC_PORT")
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
		"--chain-id", p.app.n(),
		"--keyring-backend", "test",
		"--amount", conf.Validator.Staked,
	)
}

func (p *stargatePlugin) PostInit() error {
	if err := p.apptoml(); err != nil {
		return err
	}
	return p.configtoml()
}

func (p *stargatePlugin) apptoml() error {
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
	if apiPort != "" {
		config.Set("api.address", "tcp://0.0.0.0:"+apiPort)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) configtoml() error {
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
	if rpcPort != "" {
		config.Set("rpc.laddr", "tcp://0.0.0.0:"+rpcPort)
	}
	if p2pPort != "" {
		config.Set("p2p.laddr", "tcp://0.0.0.0:"+p2pPort)
	}
	if profPort != "" {
		config.Set("prof_laddr", "localhost:"+profPort)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) StartCommands() [][]step.Option {
	c := []string{
		p.app.d(),
		"start",
		"--pruning=nothing",
	}
	if grpcPort != "" {
		c = append(c, "--grpc.address", "0.0.0.0:"+grpcPort)
	}
	return [][]step.Option{
		step.NewOptions().
			Add(
				step.Exec(c[0], c[1:]...),
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
