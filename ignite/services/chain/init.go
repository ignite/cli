package chain

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/ignite/pkg/cliui/view/accountview"
	"github.com/ignite/cli/ignite/pkg/confile"
	"github.com/ignite/cli/ignite/pkg/events"
)

type (
	// InitArgs represents argument to add additional initialization for the chain.
	// InitAccounts initializes chain accounts from the Ignite config.
	// InitConfiguration initializes node configuration from the Ignite config.
	// InitGenesis initializes genesis state for the chain from Ignite config.
	InitArgs struct {
		InitAccounts      bool
		InitConfiguration bool
		InitGenesis       bool
	}
)

const (
	moniker = "mynode"
)

var (
	// InitArgsAll performs all initialization for the chain.
	InitArgsAll = InitArgs{
		InitAccounts:      true,
		InitConfiguration: true,
		InitGenesis:       true,
	}

	// InitArgsNone performs minimal initialization for the chain by only initializing a node.
	InitArgsNone = InitArgs{
		InitAccounts:      false,
		InitConfiguration: false,
		InitGenesis:       false,
	}
)

// Init initializes the chain and accounts.
func (c *Chain) Init(ctx context.Context, args InitArgs) error {
	if err := c.InitChain(ctx, args.InitConfiguration, args.InitGenesis); err != nil {
		return err
	}

	if args.InitAccounts {
		conf, err := c.Config()
		if err != nil {
			return &CannotBuildAppError{err}
		}

		return c.InitAccounts(ctx, conf)
	}
	return nil
}

// InitChain initializes the chain.
func (c *Chain) InitChain(ctx context.Context, initConfiguration, initGenesis bool) error {
	chainID, err := c.ID()
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

	var conf *chainconfig.Config
	if initConfiguration || initGenesis {
		conf, err = c.Config()
		if err != nil {
			return err
		}
	}

	// ovewrite app config files with the values defined in Ignite's config file
	if initConfiguration {
		if err := c.Configure(home, conf); err != nil {
			return err
		}
	}

	if initGenesis {
		// make sure that chain id given during chain.New() has the most priority.
		if conf.Genesis != nil {
			conf.Genesis["chain_id"] = chainID
		}

		// update genesis file with the genesis values defined in the config
		if err := c.UpdateGenesisFile(conf.Genesis); err != nil {
			return err
		}
	}

	return nil
}

// InitAccounts initializes the chain accounts and creates validator gentxs.
func (c *Chain) InitAccounts(ctx context.Context, cfg *chainconfig.Config) error {
	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	c.ev.Send("Initializing accounts...", events.ProgressUpdate())

	var accounts accountview.Accounts

	// add accounts from config into genesis
	for _, account := range cfg.Accounts {
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
			accounts = accounts.Append(accountview.NewAccount(
				generatedAccount.Name,
				accountAddress,
				accountview.WithMnemonic(generatedAccount.Mnemonic),
			))
		} else {
			accounts = accounts.Append(accountview.NewAccount(account.Name, accountAddress))
		}
	}

	c.ev.SendView(accounts, events.ProgressFinish())

	// 0 length validator set when using network config
	if len(cfg.Validators) != 0 {
		_, err = c.IssueGentx(ctx, createValidatorFromConfig(cfg))
	}

	return err
}

// IssueGentx generates a gentx from the validator information in chain config and imports it in the chain genesis.
func (c Chain) IssueGentx(ctx context.Context, v Validator) (string, error) {
	commands, err := c.Commands(ctx)
	if err != nil {
		return "", err
	}

	// create the gentx from the validator from the config
	gentxPath, err := c.Gentx(ctx, commands, v)
	if err != nil {
		return "", err
	}

	// import the gentx into the genesis
	return gentxPath, commands.CollectGentxs(ctx)
}

// IsInitialized checks if the chain is initialized.
// The check is performed by checking if the gentx dir exists in the config.
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

// UpdateGenesisFile updates the chain genesis with a generic map of data.
// Updates are made using an override merge strategy.
func (c Chain) UpdateGenesisFile(data map[string]interface{}) error {
	path, err := c.GenesisPath()
	if err != nil {
		return err
	}

	genesis := make(map[string]interface{})
	cf := confile.New(confile.DefaultJSONEncodingCreator, path)
	if err := cf.Load(&genesis); err != nil {
		return err
	}

	if err := mergo.Merge(&genesis, data, mergo.WithOverride); err != nil {
		return err
	}

	return cf.Save(genesis)
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

func createValidatorFromConfig(conf *chainconfig.Config) (validator Validator) {
	// Currently, we support the config file with one valid validator.
	validatorFromConfig := conf.Validators[0]
	validator.Name = validatorFromConfig.Name
	validator.StakingAmount = validatorFromConfig.Bonded

	if validatorFromConfig.Gentx != nil {
		if validatorFromConfig.Gentx.Amount != "" {
			validator.StakingAmount = validatorFromConfig.Gentx.Amount
		}
		if validatorFromConfig.Gentx.Moniker != "" {
			validator.Moniker = validatorFromConfig.Gentx.Moniker
		}
		if validatorFromConfig.Gentx.CommissionRate != "" {
			validator.CommissionRate = validatorFromConfig.Gentx.CommissionRate
		}
		if validatorFromConfig.Gentx.CommissionMaxRate != "" {
			validator.CommissionMaxRate = validatorFromConfig.Gentx.CommissionMaxRate
		}
		if validatorFromConfig.Gentx.CommissionMaxChangeRate != "" {
			validator.CommissionMaxChangeRate = validatorFromConfig.Gentx.CommissionMaxChangeRate
		}
		if validatorFromConfig.Gentx.GasPrices != "" {
			validator.GasPrices = validatorFromConfig.Gentx.GasPrices
		}
		if validatorFromConfig.Gentx.Details != "" {
			validator.Details = validatorFromConfig.Gentx.Details
		}
		if validatorFromConfig.Gentx.Identity != "" {
			validator.Identity = validatorFromConfig.Gentx.Identity
		}
		if validatorFromConfig.Gentx.Website != "" {
			validator.Website = validatorFromConfig.Gentx.Website
		}
		if validatorFromConfig.Gentx.SecurityContact != "" {
			validator.SecurityContact = validatorFromConfig.Gentx.SecurityContact
		}
		if validatorFromConfig.Gentx.MinSelfDelegation != "" {
			validator.MinSelfDelegation = validatorFromConfig.Gentx.MinSelfDelegation
		}
	}
	return validator
}
