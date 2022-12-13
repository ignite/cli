package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pelletier/go-toml"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/ignite/pkg/xurl"
)

// Gentx wraps the "testd gentx"  command for generating a gentx for a validator.
// Returns path of generated gentx.
func (c Chain) Gentx(ctx context.Context, runner chaincmdrunner.Runner, v Validator) (path string, err error) {
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

// Start wraps the "appd start" command to begin running a chain from the daemon.
func (c Chain) Start(ctx context.Context, runner chaincmdrunner.Runner, cfg *chainconfig.Config) error {
	validator, err := chainconfig.FirstValidator(cfg)
	if err != nil {
		return err
	}

	servers, err := validator.GetServers()
	if err != nil {
		return err
	}

	err = runner.Start(ctx, "--pruning", "nothing", "--grpc.address", servers.GRPC.Address)

	return &CannotStartAppError{runner.Cmd().Name(), err}
}

// Configure sets the runtime configurations files for a chain (app.toml, client.toml, config.toml).
func (c Chain) Configure(homePath string, cfg *chainconfig.Config) error {
	if err := c.appTOML(homePath, cfg); err != nil {
		return err
	}
	if err := c.clientTOML(homePath, cfg); err != nil {
		return err
	}
	return c.configTOML(homePath, cfg)
}

func (c Chain) appTOML(homePath string, cfg *chainconfig.Config) error {
	validator, err := chainconfig.FirstValidator(cfg)
	if err != nil {
		return err
	}

	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(homePath, "config/app.toml")
	appConfig, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

	servers, err := validator.GetServers()
	if err != nil {
		return err
	}

	apiAddr, err := xurl.TCP(servers.API.Address)
	if err != nil {
		return fmt.Errorf("invalid api address format %s: %w", servers.API.Address, err)
	}

	// Set default config values
	appConfig.Set("api.enable", true)
	appConfig.Set("api.enabled-unsafe-cors", true)
	appConfig.Set("rpc.cors_allowed_origins", []string{"*"})

	// Update config values with the validator's Cosmos SDK app config
	updateTomlTreeValues(appConfig, validator.App)

	// Make sure the API address have the protocol prefix
	appConfig.Set("api.address", apiAddr)

	staked, err := sdktypes.ParseCoinNormalized(validator.Bonded)
	if err != nil {
		return err
	}
	gas := sdktypes.NewInt64Coin(staked.Denom, 0)
	appConfig.Set("minimum-gas-prices", gas.String())

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = appConfig.WriteTo(file)
	return err
}

func (c Chain) configTOML(homePath string, cfg *chainconfig.Config) error {
	validator, err := chainconfig.FirstValidator(cfg)
	if err != nil {
		return err
	}

	// TODO find a better way in order to not delete comments in the toml.yml
	path := filepath.Join(homePath, "config/config.toml")
	tmConfig, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

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

	// Set default config values
	tmConfig.Set("mode", "validator")
	tmConfig.Set("rpc.cors_allowed_origins", []string{"*"})
	tmConfig.Set("consensus.timeout_commit", "1s")
	tmConfig.Set("consensus.timeout_propose", "1s")

	// Update config values with the validator's Tendermint config
	updateTomlTreeValues(tmConfig, validator.Config)

	// Make sure the addresses have the protocol prefix
	tmConfig.Set("rpc.laddr", rpcAddr)
	tmConfig.Set("p2p.laddr", p2pAddr)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = tmConfig.WriteTo(file)
	return err
}

func (c Chain) clientTOML(homePath string, cfg *chainconfig.Config) error {
	validator, err := chainconfig.FirstValidator(cfg)
	if err != nil {
		return err
	}

	path := filepath.Join(homePath, "config/client.toml")
	tmConfig, err := toml.LoadFile(path)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	// Set default config values
	tmConfig.Set("keyring-backend", "test")
	tmConfig.Set("broadcast-mode", "block")

	// Update config values with the validator's client config
	updateTomlTreeValues(tmConfig, validator.Client)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = tmConfig.WriteTo(file)
	return err
}

func (c Chain) appHome() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "."+c.app.Name)
}

func updateTomlTreeValues(t *toml.Tree, values map[string]interface{}) {
	for name, v := range values {
		// Map are treated as TOML sections where the section names are the key values
		if m, ok := v.(map[string]interface{}); ok {
			section := name

			for name, v := range m {
				path := fmt.Sprintf("%s.%s", section, name)

				t.Set(path, v)
			}
		} else {
			// By default set top a level key/value
			t.Set(name, v)
		}
	}
}
