package chain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/nqd/flat"
	"github.com/pelletier/go-toml"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xurl"
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

func (c Chain) InPlace(ctx context.Context, runner chaincmdrunner.Runner, args InPlaceArgs) error {
	err := runner.InPlace(ctx,
		args.NewChainID,
		args.NewOperatorAddress,
		chaincmd.InPlaceWithPrvKey(args.PrvKeyValidator),
		chaincmd.InPlaceWithAccountToFund(args.AccountsToFund),
		chaincmd.InPlaceWithSkipConfirmation(),
	)
	return err
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
func (c Chain) Configure(homePath, chainID string, cfg *chainconfig.Config) error {
	if err := appTOML(homePath, cfg); err != nil {
		return err
	}
	if err := clientTOML(homePath, chainID, cfg); err != nil {
		return err
	}
	return configTOML(homePath, cfg)
}

func appTOML(homePath string, cfg *chainconfig.Config) error {
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
		return errors.Errorf("invalid api address format %s: %w", servers.API.Address, err)
	}

	// Set default config values
	appConfig.Set("api.enable", true)
	appConfig.Set("api.enabled-unsafe-cors", true)
	appConfig.Set("rpc.cors_allowed_origins", []string{"*"})

	staked, err := sdktypes.ParseCoinNormalized(validator.Bonded)
	if err != nil {
		return err
	}
	gas := sdktypes.NewInt64Coin(staked.Denom, 0)
	appConfig.Set("minimum-gas-prices", gas.String())

	// Update config values with the validator's Cosmos SDK app config
	if err := updateTomlTreeValues(appConfig, validator.App); err != nil {
		return err
	}

	// Make sure the API address have the protocol prefix
	appConfig.Set("api.address", apiAddr)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = appConfig.WriteTo(file)
	return err
}

func configTOML(homePath string, cfg *chainconfig.Config) error {
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
		return errors.Errorf("invalid rpc address format %s: %w", servers.RPC.Address, err)
	}

	p2pAddr, err := xurl.TCP(servers.P2P.Address)
	if err != nil {
		return errors.Errorf("invalid p2p address format %s: %w", servers.P2P.Address, err)
	}

	// Set default config values
	tmConfig.Set("mode", "validator")
	tmConfig.Set("rpc.cors_allowed_origins", []string{"*"})
	tmConfig.Set("consensus.timeout_commit", "1s")
	tmConfig.Set("consensus.timeout_propose", "1s")

	// Update config values with the validator's Tendermint config
	if err := updateTomlTreeValues(tmConfig, validator.Config); err != nil {
		return err
	}

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

func clientTOML(homePath, chainID string, cfg *chainconfig.Config) error {
	validator, err := chainconfig.FirstValidator(cfg)
	if err != nil {
		return err
	}

	path := filepath.Join(homePath, "config/client.toml")
	clientConfig, err := toml.LoadFile(path)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	// Set default config values
	clientConfig.Set("chain-id", chainID)
	clientConfig.Set("keyring-backend", "test")
	clientConfig.Set("broadcast-mode", "sync")

	// Update config values with the validator's client config
	if err := updateTomlTreeValues(clientConfig, validator.Client); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = clientConfig.WriteTo(file)
	return err
}

func (c Chain) appHome() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "."+c.app.Name)
}

func updateTomlTreeValues(t *toml.Tree, values map[string]interface{}) error {
	flatValues, err := flat.Flatten(values, nil)
	if err != nil {
		return err
	}

	for name, v := range flatValues {
		t.Set(name, v)
	}
	return nil
}
