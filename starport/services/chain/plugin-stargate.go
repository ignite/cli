package chain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

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

func (p *stargatePlugin) Binaries() []string {
	return []string{
		p.app.D(),
	}
}

func (p *stargatePlugin) AddUserCommand(cmd chaincmd.ChainCmd, accountName string) step.Options {
	return step.NewOptions().Add(cmd.AddKeyCommand(accountName))
}

func (p *stargatePlugin) ImportUserCommand(cmd chaincmd.ChainCmd, name, mnemonic string) step.Options {
	return step.NewOptions().
		Add(
			cmd.ImportKeyCommand(name),
			step.Write([]byte(mnemonic+"\n")),
		)
}

func (p *stargatePlugin) ShowAccountCommand(cmd chaincmd.ChainCmd, accountName string) step.Option {
	return cmd.ShowKeyAddressCommand(accountName)
}

func (p *stargatePlugin) ConfigCommands(_ chaincmd.ChainCmd, _ string) []step.Option {
	return nil
}

func (p *stargatePlugin) GentxCommand(cmd chaincmd.ChainCmd, v Validator) step.Option {
	return cmd.GentxCommand(
		v.Name,
		v.StakingAmount,
		chaincmd.GentxWithMoniker(v.Moniker),
		chaincmd.GentxWithCommissionRate(v.CommissionRate),
		chaincmd.GentxWithCommissionMaxRate(v.CommissionMaxRate),
		chaincmd.GentxWithCommissionMaxChangeRate(v.CommissionMaxChangeRate),
		chaincmd.GentxWithMinSelfDelegation(v.MinSelfDelegation),
		chaincmd.GentxWithGasPrices(v.GasPrices),
	)
}

func (p *stargatePlugin) StartCommands(cmd chaincmd.ChainCmd, conf starportconf.Config) [][]step.Option {
	return [][]step.Option{
		step.NewOptions().
			Add(
				cmd.StartCommand(
					"--pruning",
					"nothing",
					"--grpc.address",
					conf.Servers.GRPCAddr,
				),
				step.PostExec(func(exitErr error) error {
					return errors.Wrapf(exitErr, "cannot run %[1]vd start", p.app.Name)
				}),
			),
	}
}

func (p *stargatePlugin) PostInit(conf starportconf.Config) error {
	if err := p.apptoml(conf); err != nil {
		return err
	}
	return p.configtoml(conf)
}

func (p *stargatePlugin) apptoml(conf starportconf.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(p.Home(), "config/app.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	config.Set("api.enable", true)
	config.Set("api.enabled-unsafe-cors", true)
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("api.address", xurl.TCP(conf.Servers.APIAddr))
	config.Set("grpc.address", conf.Servers.GRPCAddr)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) configtoml(conf starportconf.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(p.Home(), "config/config.toml")
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
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) StoragePaths() []string {
	return []string{
		p.Home(),
	}
}

func (p *stargatePlugin) Home() string {
	return stargateHome(p.app)
}

func (p *stargatePlugin) CLIHome() string {
	return stargateHome(p.app)
}

func stargateHome(app App) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "."+app.N())
}

func (p *stargatePlugin) Version() cosmosver.MajorVersion { return cosmosver.Stargate }

func (p *stargatePlugin) SupportsIBC() bool { return true }
