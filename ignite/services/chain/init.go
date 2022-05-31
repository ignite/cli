package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"

	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/ignite/pkg/confile"
)

const (
	moniker = "mynode"
)

// Init initializes the chain and applies all optional configurations.
func (c *Chain) Init(ctx context.Context, initAccounts bool) error {
	conf, err := c.Config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

	if err := c.InitChain(ctx); err != nil {
		return err
	}

	if initAccounts {
		return c.InitAccounts(ctx, conf)
	}
	return nil
}

// InitChain initializes the chain.
func (c *Chain) InitChain(ctx context.Context) error {
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

	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	// init node.
	if err := commands.Init(ctx, moniker); err != nil {
		return err
	}

	// overwrite configuration changes from Ignite CLI's config.yml to
	// over app's sdk configs.

	if err := c.plugin.Configure(home, conf); err != nil {
		return err
	}

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
	clientTOMLPath, err := c.ClientTOMLPath()
	if err != nil {
		return err
	}
	configTOMLPath, err := c.ConfigTOMLPath()
	if err != nil {
		return err
	}

	validator := conf.Validators[0]
	config := make(map[string]interface{})
	for key, value := range validator.Config {
		config[key] = value
	}
	delete(config, "rpc")
	delete(config, "p2p")
	delete(config, "pprof_laddr")
	if len(config) == 0 {
		config = nil
	}

	app := make(map[string]interface{})
	for key, value := range validator.App {
		app[key] = value
	}
	delete(app, "grpc")
	delete(app, "grpc-web")
	delete(app, "api")

	if len(app) == 0 {
		app = nil
	}
	appconfigs := []struct {
		ec      confile.EncodingCreator
		path    string
		changes map[string]interface{}
	}{
		{confile.DefaultJSONEncodingCreator, genesisPath, conf.Genesis},
		{confile.DefaultTOMLEncodingCreator, appTOMLPath, app},
		{confile.DefaultTOMLEncodingCreator, clientTOMLPath, validator.Client},
		{confile.DefaultTOMLEncodingCreator, configTOMLPath, config},
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

	return nil
}

// InitAccounts initializes the chain accounts and creates validator gentxs
func (c *Chain) InitAccounts(ctx context.Context, conf *v1.Config) error {
	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	// add accounts from config into genesis
	for _, account := range conf.ListAccounts() {
		var generatedAccount chaincmdrunner.Account
		accountAddress := account.Address

		// If the account doesn't provide an address, we create one
		if accountAddress == "" {
			generatedAccount, err = commands.AddAccount(ctx, account.Name, account.Mnemonic, account.CoinType)
			if err != nil {
				return err
			}
			accountAddress = generatedAccount.Address
		}

		coins := strings.Join(account.Coins, ",")
		if err := commands.AddGenesisAccount(ctx, accountAddress, coins); err != nil {
			return err
		}

		if account.Address == "" {
			fmt.Fprintf(
				c.stdLog().out,
				"🙂 Created account %q with address %q with mnemonic: %q\n",
				generatedAccount.Name,
				generatedAccount.Address,
				generatedAccount.Mnemonic,
			)
		} else {
			fmt.Fprintf(
				c.stdLog().out,
				"🙂 Imported an account %q with address: %q\n",
				account.Name,
				account.Address,
			)
		}
	}

	_, err = c.IssueGentx(ctx, createValidatorFromConfig(conf))

	return err
}

// IssueGentx generates a gentx from the validator information in chain config and import it in the chain genesis
func (c Chain) IssueGentx(ctx context.Context, v Validator) (string, error) {
	commands, err := c.Commands(ctx)
	if err != nil {
		return "", err
	}

	// create the gentx from the validator from the config
	gentxPath, err := c.plugin.Gentx(ctx, commands, v)
	if err != nil {
		return "", err
	}

	// import the gentx into the genesis
	return gentxPath, commands.CollectGentxs(ctx)
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

type Validator struct {
	Name                    string
	Moniker                 string
	StakingAmount           string
	CommissionRate          string
	CommissionMaxRate       string
	CommissionMaxChangeRate string
	MinSelfDelegation       string
	GasPrices               string
	Details                 string
	Identity                string
	Website                 string
	SecurityContact         string
}

// Account represents an account in the chain.
type Account struct {
	Name     string
	Address  string
	Mnemonic string `json:"mnemonic"`
	CoinType string
	Coins    string
}

func createValidatorFromConfig(conf *v1.Config) Validator {
	// Currently, we support the config file with one valid validator.
	validator := conf.ListValidators()[0]
	name := validator.Name
	monikerV := ""
	commissionRate := ""
	commissionMaxRate := ""
	commissionMaxChangeRate := ""
	gasPrices := ""
	details := ""
	identity := ""
	website := ""
	securityContact := ""
	minSelfDelegation := ""
	stakingAmount := validator.Bonded

	if validator.Gentx != nil {
		if validator.Gentx.Amount != "" {
			stakingAmount = validator.Gentx.Amount
		}
		if validator.Gentx.Moniker != "" {
			monikerV = validator.Gentx.Moniker
		}
		if validator.Gentx.CommissionRate != "" {
			commissionRate = validator.Gentx.CommissionRate
		}
		if validator.Gentx.CommissionMaxRate != "" {
			commissionMaxRate = validator.Gentx.CommissionMaxRate
		}
		if validator.Gentx.CommissionMaxChangeRate != "" {
			commissionMaxChangeRate = validator.Gentx.CommissionMaxChangeRate
		}
		if validator.Gentx.GasPrices != "" {
			gasPrices = validator.Gentx.GasPrices
		}
		if validator.Gentx.Details != "" {
			details = validator.Gentx.Details
		}
		if validator.Gentx.Identity != "" {
			identity = validator.Gentx.Identity
		}
		if validator.Gentx.Website != "" {
			website = validator.Gentx.Website
		}
		if validator.Gentx.SecurityContact != "" {
			securityContact = validator.Gentx.SecurityContact
		}
		if validator.Gentx.MinSelfDelegation != "" {
			minSelfDelegation = validator.Gentx.MinSelfDelegation
		}
	}

	return Validator{
		Name:                    name,
		StakingAmount:           stakingAmount,
		Moniker:                 monikerV,
		CommissionRate:          commissionRate,
		CommissionMaxRate:       commissionMaxRate,
		CommissionMaxChangeRate: commissionMaxChangeRate,
		MinSelfDelegation:       minSelfDelegation,
		GasPrices:               gasPrices,
		Details:                 details,
		Identity:                identity,
		Website:                 website,
		SecurityContact:         securityContact,
	}
}
