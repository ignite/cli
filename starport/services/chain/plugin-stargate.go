package chain

import (
	"context"
	"os"
	"path/filepath"

	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/pelletier/go-toml"
	starportconf "github.com/tendermint/starport/starport/chainconf"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xurl"
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

func (p *stargatePlugin) Setup(ctx context.Context) error {
	return nil
}

func (p *stargatePlugin) Configure(_ context.Context, _ chaincmdrunner.Runner, _ string) error {
	return nil
}

func (p *stargatePlugin) Gentx(ctx context.Context, runner chaincmdrunner.Runner, v Validator) (path string, err error) {
	return runner.Gentx(
		ctx,
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

func (p *stargatePlugin) PostInit(homePath string, conf starportconf.Config) error {
	if err := p.apptoml(homePath, conf); err != nil {
		return err
	}
	return p.configtoml(homePath, conf)
}

func (p *stargatePlugin) apptoml(homePath string, conf starportconf.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(homePath, "config/app.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	config.Set("api.enable", true)
	config.Set("api.enabled-unsafe-cors", true)
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("api.address", xurl.TCP(conf.Host.API))
	config.Set("grpc.address", conf.Host.GRPC)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) configtoml(homePath string, conf starportconf.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(homePath, "config/config.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("consensus.timeout_commit", "1s")
	config.Set("consensus.timeout_propose", "1s")
	config.Set("rpc.laddr", xurl.TCP(conf.Host.RPC))
	config.Set("p2p.laddr", xurl.TCP(conf.Host.P2P))
	config.Set("rpc.pprof_laddr", conf.Host.Prof)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) Start(ctx context.Context, runner chaincmdrunner.Runner, conf starportconf.Config) error {
	err := runner.Start(ctx,
		"--pruning",
		"nothing",
		"--grpc.address",
		conf.Host.GRPC,
	)
	return &CannotStartAppError{p.app.Name, err}
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
