package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pelletier/go-toml"

	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
	"github.com/ignite/cli/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/xurl"
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
		chaincmd.GentxWithDetails(v.Details),
		chaincmd.GentxWithIdentity(v.Identity),
		chaincmd.GentxWithWebsite(v.Website),
		chaincmd.GentxWithSecurityContact(v.SecurityContact),
	)
}

func (p *stargatePlugin) Configure(homePath string, cfg *v1.Config) error {
	if err := p.appTOML(homePath, cfg); err != nil {
		return err
	}
	if err := p.clientTOML(homePath); err != nil {
		return err
	}
	return p.configTOML(homePath, cfg)
}

func (p *stargatePlugin) appTOML(homePath string, cfg *v1.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(homePath, "config/app.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

	validator := cfg.Validators[0]
	servers, err := validator.GetServers()
	if err != nil {
		return err
	}

	apiAddr, err := xurl.TCP(servers.API.Address)
	if err != nil {
		return fmt.Errorf("invalid api address format %s: %w", servers.API.Address, err)
	}

	config.Set("api.enable", true)
	config.Set("api.enabled-unsafe-cors", true)
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("api.address", apiAddr)
	config.Set("grpc.address", servers.GRPC.Address)
	config.Set("grpc-web.address", servers.GRPCWeb.Address)

	staked, err := sdktypes.ParseCoinNormalized(validator.Bonded)
	if err != nil {
		return err
	}
	gas := sdktypes.NewInt64Coin(staked.Denom, 0)
	config.Set("minimum-gas-prices", gas.String())

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) configTOML(homePath string, cfg *v1.Config) error {
	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(homePath, "config/config.toml")
	config, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

	validator := cfg.Validators[0]
	servers, err := validator.GetServers()
	if err != nil {
		return err
	}

	rpcAddr, err := xurl.TCP(servers.RPC.Address)
	if err != nil {
		return fmt.Errorf("invalid rpc address format %s: %w", servers.RPC.Address, err)
	}

	p2pAddr, err := xurl.TCP(servers.P2P.Address)
	if err != nil {
		return fmt.Errorf("invalid p2p address format %s: %w", servers.P2P.Address, err)
	}

	config.Set("mode", "validator")
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("consensus.timeout_commit", "1s")
	config.Set("consensus.timeout_propose", "1s")
	config.Set("rpc.laddr", rpcAddr)
	config.Set("p2p.laddr", p2pAddr)
	config.Set("rpc.pprof_laddr", servers.RPC.PProfAddress)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) clientTOML(homePath string) error {
	path := filepath.Join(homePath, "config/client.toml")
	config, err := toml.LoadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	config.Set("keyring-backend", "test")
	config.Set("broadcast-mode", "block")
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	return err
}

func (p *stargatePlugin) Start(ctx context.Context, runner chaincmdrunner.Runner, cfg *v1.Config) error {
	validator := cfg.Validators[0]
	servers, err := validator.GetServers()
	if err != nil {
		return err
	}

	err = runner.Start(ctx, "--pruning", "nothing", "--grpc.address", servers.GRPC.Address)

	return &CannotStartAppError{p.app.Name, err}
}

func (p *stargatePlugin) Home() string {
	return stargateHome(p.app)
}

func stargateHome(app App) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "."+app.Name)
}

func (p *stargatePlugin) Version() cosmosver.Family { return cosmosver.Stargate }

func (p *stargatePlugin) SupportsIBC() bool { return true }
