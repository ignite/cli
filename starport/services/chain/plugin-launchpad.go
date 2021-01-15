package chain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"golang.org/x/sync/errgroup"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xurl"
	starportconf "github.com/tendermint/starport/starport/services/chain/conf"
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

func (p *launchpadPlugin) Setup(ctx context.Context) error {
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

func (p *launchpadPlugin) Binaries() []string {
	return []string{
		p.app.D(),
		p.app.CLI(),
	}
}

func (p *launchpadPlugin) Configure(ctx context.Context, runner chaincmdrunner.Runner, chainID string) error {
	return runner.LaunchpadSetConfigs(ctx,
		chaincmdrunner.NewKV("keyring-backend", "test"),
		chaincmdrunner.NewKV("chain-id", chainID),
		chaincmdrunner.NewKV("output", "json"),
		chaincmdrunner.NewKV("indent", "true"),
		chaincmdrunner.NewKV("trust-node", "true"),
	)
}

func (p *launchpadPlugin) Gentx(ctx context.Context, runner chaincmdrunner.Runner, v Validator) (path string, err error) {
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

func (p *launchpadPlugin) PostInit(homePath string, conf starportconf.Config) error {
	return p.configtoml(homePath, conf)
}

func (p *launchpadPlugin) configtoml(homePath string, conf starportconf.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(homePath, "config/config.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}
	config.Set("rpc.cors_allowed_origins", []string{"*"})
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

func (p *launchpadPlugin) Start(ctx context.Context, runner chaincmdrunner.Runner, conf starportconf.Config) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := runner.Start(ctx)
		return errors.Wrapf(err, "cannot run %[1]vd start", p.app.Name)
	})

	g.Go(func() error {
		err := runner.LaunchpadStartRestServer(ctx, xurl.TCP(conf.Servers.APIAddr), xurl.TCP(conf.Servers.RPCAddr))
		return errors.Wrapf(err, "cannot run %[1]vcli rest-server", p.app.Name)
	})

	return g.Wait()
}

func (p *launchpadPlugin) StoragePaths() []string {
	return []string{
		launchpadHome(p.app),
		launchpadCLIHome(p.app),
	}
}

func (p *launchpadPlugin) Home() string {
	return launchpadHome(p.app)
}

func (p *launchpadPlugin) CLIHome() string {
	return launchpadCLIHome(p.app)
}

func launchpadHome(app App) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "."+app.ND())
}

func launchpadCLIHome(app App) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "."+app.NCLI())
}

func (p *launchpadPlugin) Version() cosmosver.MajorVersion { return cosmosver.Launchpad }

func (p *launchpadPlugin) SupportsIBC() bool { return false }
