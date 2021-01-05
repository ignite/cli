package chain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
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

func (p *launchpadPlugin) AddUserCommand(cmd chaincmd.ChainCmd, accountName string) step.Options {
	return step.NewOptions().Add(cmd.LaunchpadAddKeyCommand(accountName))
}

func (p *launchpadPlugin) ImportUserCommand(cmd chaincmd.ChainCmd, name, mnemonic string) step.Options {
	return step.NewOptions().
		Add(
			cmd.LaunchpadImportKeyCommand(name),
			step.Write([]byte(mnemonic+"\n")),
		)
}

func (p *launchpadPlugin) ShowAccountCommand(cmd chaincmd.ChainCmd, accountName string) step.Option {
	return cmd.LaunchpadShowKeyAddressCommand(accountName)
}

func (p *launchpadPlugin) ConfigCommands(cmd chaincmd.ChainCmd, chainID string) []step.Option {
	return []step.Option{
		cmd.LaunchpadSetConfigCommand("keyring-backend", "test"),
		cmd.LaunchpadSetConfigCommand("chain-id", chainID),
		cmd.LaunchpadSetConfigCommand("output", "json"),
		cmd.LaunchpadSetConfigCommand("indent", "true"),
		cmd.LaunchpadSetConfigCommand("trust-node", "true"),
	}
}

func (p *launchpadPlugin) GentxCommand(cmd chaincmd.ChainCmd, v Validator) step.Option {
	return cmd.LaunchpadGentxCommand(
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

func (p *launchpadPlugin) StartCommands(cmd chaincmd.ChainCmd, conf starportconf.Config) [][]step.Option {
	return [][]step.Option{
		step.NewOptions().
			Add(
				cmd.StartCommand(),
				step.PostExec(func(exitErr error) error {
					return errors.Wrapf(exitErr, "cannot run %[1]vd start", p.app.Name)
				}),
			),
		step.NewOptions().
			Add(
				cmd.LaunchpadRestServerCommand(xurl.TCP(conf.Servers.APIAddr), xurl.TCP(conf.Servers.RPCAddr)),
				step.PostExec(func(exitErr error) error {
					return errors.Wrapf(exitErr, "cannot run %[1]vcli rest-server", p.app.Name)
				}),
			),
	}
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
	path := filepath.Join(home, "."+p.app.ND(), "config/config.toml")
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

func (p *launchpadPlugin) StoragePaths() []string {
	return []string{
		launchpadHome(p.app),
		launchpadCLIHome(p.app),
	}
}

func (p *launchpadPlugin) Home() string {
	return launchpadHome(p.app)
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
