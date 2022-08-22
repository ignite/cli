package chain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/tendermint/tendermint/p2p"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/ignite/pkg/confile"
)

const (
	moniker  = "mynode"
	localnet = "0.0.0.0"
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
		if err := c.InitAccounts(ctx, conf); err != nil {
			return err
		}
	}

	return c.symlinkValidatorConfig(ctx)
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
	home := c.AppHome()
	if err := os.RemoveAll(home); err != nil {
		return err
	}

	validatorNodeIDs := make([]p2p.ID, len(conf.Validators))

	// for each validator
	for i, validator := range conf.Validators {
		commands, err := c.Commands(ctx, validator)
		if err != nil {
			return err
		}

		// init node.
		if err := commands.Init(ctx, validator.Moniker()); err != nil {
			return err
		}

		// make sure that chain id given during chain.New() has the most priority.
		if conf.Genesis != nil {
			conf.Genesis["chain_id"] = chainID
		}

		// todo: update for each validator
		// now: hard code first validator
		// validator := c.validator
		appTOMLPath, err := c.appTOMLPathForValidator(validator)
		if err != nil {
			return err
		}
		clientTOMLPath, err := c.clientTOMLPathForValidator(validator)
		if err != nil {
			return err
		}
		configTOMLPath, err := c.configTOMLPathForValidator(validator)
		if err != nil {
			return err
		}
		appconfigs := []appconfig{
			{confile.DefaultTOMLEncodingCreator, appTOMLPath, validator.App},
			{confile.DefaultTOMLEncodingCreator, clientTOMLPath, validator.Client},
			{confile.DefaultTOMLEncodingCreator, configTOMLPath, validator.Config},
		}

		for _, ac := range appconfigs {
			applyConfig(ac)
		}

		vhome := c.homeForValidator(validator)
		if err := c.plugin.Configure(vhome, validator); err != nil {
			return err
		}

		// collect Validator NodeID
		nodeKeyPath := filepath.Join(vhome, "config/node_key.json")
		nodeKey, err := p2p.LoadNodeKey(nodeKeyPath)
		if err != nil {
			return err
		}
		validatorNodeIDs[i] = nodeKey.ID()
	}

	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}
	keyringPath := filepath.Join(c.AppHome(), "keyring-test")
	genTxPath := filepath.Join(c.AppHome(), "config/gentx")

	for i, val := range conf.Validators {
		vhome := c.homeForValidator(val)
		vgenesisPath := filepath.Join(vhome, "config/genesis.json")
		vkeyringPath := filepath.Join(vhome, "keyring-test")
		vgenTxPath := filepath.Join(vhome, "config/gentx")

		// copy the initialized genesis from the first validator into the app home
		if i == 0 { // only run on first validator
			buf, err := os.ReadFile(vgenesisPath)
			if err != nil {
				return err
			}

			appConfigPath := filepath.Join(c.AppHome(), "config")
			// ensure the config folder exists
			if err := ensureDirectory(appConfigPath); err != nil {
				return err
			}
			// copy the genesis path to the app root config
			if err := os.WriteFile(genesisPath, buf, 0644); err != nil {
				return err
			}

			// ensure the keyring-test folder exists
			if err := ensureDirectory(keyringPath); err != nil {
				return err
			}

			// ensure the gentx folder exists
			if err := ensureDirectory(genTxPath); err != nil {
				return err
			}
		}

		// delete it from all
		// then symlink back
		if err := os.Remove(vgenesisPath); err != nil {
			return err
		}

		// symlink the root genesis path into each validator config path
		if err := os.Symlink(genesisPath, vgenesisPath); err != nil {
			return err
		}

		// symlink the root gentx path into each validator config gentx path
		if err := os.Symlink(genTxPath, vgenTxPath); err != nil {
			return err
		}

		// symlink the root keyring path into each validator keyring path
		if err := os.Symlink(keyringPath, vkeyringPath); err != nil {
			return err
		}

		// apply persistent peers to each validator config
		var persistentPeers string
		for j, nodeID := range validatorNodeIDs {
			// skip matching validators
			if i == j {
				continue
			}

			peer := formatPeerID(val, string(nodeID))
			persistentPeers += peer

			// if we aren't at the end, include a comma
			if j < len(validatorNodeIDs)-2 {
				persistentPeers += ","
			}
		}

		// apply persistentpeers to validator config

		configTOMLPath, err := c.configTOMLPathForValidator(val)
		if err != nil {
			return err
		}
		p2pConfig := val.Config["p2p"]
		if cfg, ok := p2pConfig.(map[string]interface{}); ok {
			cfg["persistent_peers"] = persistentPeers
			val.Config["p2p"] = cfg
		}

		ac := appconfig{confile.DefaultTOMLEncodingCreator, configTOMLPath, val.Config}
		if err := applyConfig(ac); err != nil {
			return err
		}
	}

	// overwrite configuration changes from Ignite CLI's config.yml to
	// over app's sdk configs.
	ac := appconfig{confile.DefaultJSONEncodingCreator, genesisPath, conf.Genesis}
	if err := applyConfig(ac); err != nil {
		return err
	}

	return nil
}

func applyConfig(ac appconfig) error {
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
	return nil
}

// InitAccounts initializes the chain accounts and creates validator gentxs
func (c *Chain) InitAccounts(ctx context.Context, conf chainconfig.Config) error {
	commands, err := c.Commands(ctx, c.validator, chaincmd.WithHome(c.AppHome()))
	if err != nil {
		return err
	}

	// add accounts from config into genesis
	for _, account := range conf.Accounts {
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

	for _, v := range conf.Validators {
		// todo: for each validator
		// now: hardcode validator
		_, err = c.IssueGentx(ctx, v)
		if err != nil {
			return err
		}
	}

	return c.CollectGentxs(ctx)
	// return nil
}

// IssueGentx generates a gentx from the validator information in chain config.
// *Does not* run `collect-gentxs`.
func (c Chain) IssueGentx(ctx context.Context, v chainconfig.Validator) (string, error) {
	commands, err := c.Commands(ctx, v)
	if err != nil {
		return "", err
	}

	// create the gentx from the validator from the config
	return c.plugin.Gentx(ctx, commands, Validator{
		Name:          v.Name,
		StakingAmount: v.Bonded,
		Moniker:       v.Moniker(),
	})
}

func (c Chain) CollectGentxs(ctx context.Context) error {
	commands, err := c.Commands(ctx, c.validator, chaincmd.WithHome(c.AppHome()))
	if err != nil {
		return nil
	}

	return commands.CollectGentxs(ctx)
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

// initSymlink will symlink the root validators config data
// into
func (c *Chain) symlinkValidatorConfig(ctx context.Context) error {
	configSymlink := []string{
		"config/app.toml",
		"config/client.toml",
		"config/config.toml",
		"config/node_key.json",
		"config/priv_validator_key.json",
	}

	// symlink the root validator's config into the main config folder
	vhome := c.homeForValidator(c.validator)
	for _, cfgpath := range configSymlink {
		rootPath := filepath.Join(c.AppHome(), cfgpath)
		valPath := filepath.Join(vhome, cfgpath)

		// make sure the target root path is deleted first
		if err := os.Remove(rootPath); err != nil {
			return err
		}

		if err := os.Symlink(valPath, rootPath); err != nil {
			return err
		}
	}

	return nil
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

// ToConfig converts this type to chainconfig.Validator
func (v Validator) ToConfig() chainconfig.Validator {
	return chainconfig.Validator{
		Name:   v.Name,
		Bonded: v.StakingAmount,
		GenTx: chainconfig.GenTx{
			Moniker:                 v.Moniker,
			CommisionRate:           v.CommissionRate,
			CommissionMaxRate:       v.CommissionMaxRate,
			CommissionMaxChangeRate: v.CommissionMaxChangeRate,
			MinSelfDelegation:       v.MinSelfDelegation,
			GasPrices:               v.GasPrices,
			Details:                 v.Details,
			Identity:                v.Identity,
			Website:                 v.Website,
			SecurityContact:         v.SecurityContact,
		},
	}
}

// Account represents an account in the chain.
type Account struct {
	Name     string
	Address  string
	Mnemonic string `json:"mnemonic"`
	CoinType string
	Coins    string
}

type appconfig struct {
	ec      confile.EncodingCreator
	path    string
	changes map[string]interface{}
}

func ensureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func formatPeerID(validator chainconfig.Validator, nodeID string) string {
	var laddr string
	if subcfg, ok := validator.Config["p2p"].(map[string]interface{}); ok {
		laddr = subcfg["laddr"].(string)
	}
	return fmt.Sprintf("%s@%s", nodeID, laddr)
}
