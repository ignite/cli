package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"

	conf "github.com/tendermint/starport/starport/chainconf"
	secretconf "github.com/tendermint/starport/starport/chainconf/secret"

	"github.com/imdario/mergo"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/xos"
)

const (
	moniker = "mynode"
)

// Init initializes chain.
func (c *Chain) Init(ctx context.Context) error {
	chainID, err := c.ID()
	if err != nil {
		return err
	}

	conf, err := c.Config()
	if err != nil {
		return err
	}

	// cleanup persistent data from previous `serve`.
	home, err := c.Home()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(home); err != nil {
		return err
	}
	cliHome, err := c.CLIHome()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(cliHome); err != nil {
		return err
	}

	// init node.
	if err := c.cmd.Init(ctx, moniker); err != nil {
		return err
	}

	// overwrite configuration changes from Starport's config.yml to
	// over app's sdk configs.

	// make sure that chain id given during chain.New() has the most priority.
	if conf.Genesis != nil {
		conf.Genesis["chain_id"] = chainID
	}

	// Initilize app config
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}
	appTOMLPath, err := c.AppTOMLPath()
	if err != nil {
		return err
	}
	configTOMLPath, err := c.ConfigTOMLPath()
	if err != nil {
		return err
	}

	appconfigs := []struct {
		ec      confile.EncodingCreator
		path    string
		changes map[string]interface{}
	}{
		{confile.DefaultJSONEncodingCreator, genesisPath, conf.Genesis},
		{confile.DefaultTOMLEncodingCreator, appTOMLPath, conf.Init.App},
		{confile.DefaultTOMLEncodingCreator, configTOMLPath, conf.Init.Config},
	}

	for _, ac := range appconfigs {
		cf := confile.New(ac.ec, ac.path)
		var conf map[string]interface{}
		if err := cf.Load(&conf); err != nil {
			return err
		}
		if err := mergo.Merge(&conf, ac.changes, mergo.WithOverride); err != nil {
			return err
		}
		if err := cf.Save(conf); err != nil {
			return err
		}
	}

	// run post init handler
	return c.plugin.PostInit(home, conf)
}

// InitAccounts initializes the chain accounts and creates validator gentxs
func (c *Chain) InitAccounts(ctx context.Context, conf conf.Config) error {
	sconf, err := secretconf.Open(c.app.Path)
	if err != nil {
		return err
	}

	// add accounts from config into genesis
	for _, account := range conf.Accounts {
		var generatedAccount chaincmdrunner.Account
		accountAddress := account.Address

		// If the account doesn't provide an address, we create one
		if accountAddress == "" {
			generatedAccount, err = c.Commands().AddAccount(ctx, account.Name, "")
			if err != nil {
				return err
			}
			accountAddress = generatedAccount.Address
		}

		coins := strings.Join(account.Coins, ",")
		if err := c.Commands().AddGenesisAccount(ctx, accountAddress, coins); err != nil {
			return err
		}

		if account.Address == "" {
			fmt.Fprintf(c.stdLog(logStarport).out, "🙂 Created an account. Password (mnemonic): %[1]v\n", generatedAccount.Mnemonic)
		} else {
			fmt.Fprintf(c.stdLog(logStarport).out, "🙂 Account imported: %[1]v\n", account.Address)
		}
	}

	// add accounts from secret config into genesis
	for _, account := range sconf.Accounts {
		acc, err := c.Commands().AddAccount(ctx, account.Name, account.Mnemonic)
		if err != nil {
			return err
		}

		coins := strings.Join(account.Coins, ",")
		if err := c.Commands().AddGenesisAccount(ctx, acc.Address, coins); err != nil {
			return err
		}
	}

	// perform configuration in the chain config
	if err := c.configure(ctx); err != nil {
		return err
	}

	// create the gentx from the validator from the config
	if _, err := c.plugin.Gentx(ctx, c.Commands(), Validator{
		Name:          conf.Validator.Name,
		StakingAmount: conf.Validator.Staked,
	}); err != nil {
		return err
	}

	// import the gentx into the genesis
	if err := c.Commands().CollectGentxs(ctx); err != nil {
		return err
	}

	return nil
}

// IsInitialized checks if the chain is initialized
// the check is performed by checking if the gentx dir exist in the config
func (c *Chain) IsInitialized() (bool, error) {
	home, err := c.Home()
	if err != nil {
		return false, err
	}
	gentxDir := filepath.Join(home, "config", "gentx")

	if _, err := os.Stat(gentxDir); os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		// Return error on other error
		return false, err
	}

	return true, nil
}

func (c *Chain) configure(ctx context.Context) error {
	// setup IBC Relayer.
	if err := c.checkIBCRelayerSupport(); err == nil {
		if err := xos.RemoveAllUnderHome(".relayer"); err != nil {
			return err
		}
		info, err := c.RelayerInfo()
		if err != nil {
			return err
		}
		fmt.Fprintf(c.stdLog(logStarport).out, "✨ Relayer info: %s\n", info)
		return nil
	}

	// configure blockchain.
	chainID, err := c.ID()
	if err != nil {
		return err
	}

	return c.plugin.Configure(ctx, c.Commands(), chainID)
}

type Validator struct {
	Name                    string
	Moniker                 string
	StakingAmount           string
	CommissionRate          string
	CommissionMaxRate       string
	CommissionMaxChangeRate string
	MinSelfDelegation       string
	GasPrices               string
}

// Account represents an account in the chain.
type Account struct {
	Name     string
	Address  string
	Mnemonic string `json:"mnemonic"`
	Coins    string
}
