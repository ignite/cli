package chain

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	cmtypes "github.com/cometbft/cometbft/abci/types"
	cmprivval "github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ccvconsumertypes "github.com/cosmos/interchain-security/v3/x/ccv/consumer/types"
	ccvtypes "github.com/cosmos/interchain-security/v3/x/ccv/types"
	"github.com/imdario/mergo"

	chainconfig "github.com/ignite/cli/v28/ignite/config/chain"
	chaincmdrunner "github.com/ignite/cli/v28/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/view/accountview"
	"github.com/ignite/cli/v28/ignite/pkg/confile"
	"github.com/ignite/cli/v28/ignite/pkg/events"
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
	if len(cfg.Validators) == 0 {
		return nil
	}
	if cfg.IsConsumerChain() {
		// Consumer chain writes validators in the consumer module genesis
		if err := c.writeConsumerGenesis(); err != nil {
			return err
		}
	} else {
		// Sovereign chain writes validators in gentxs.
		_, err := c.IssueGentx(ctx, createValidatorFromConfig(cfg))
		if err != nil {
			return err
		}
	}
	return nil
}

// writeConsumerGenesis updates the consumer module genesis in the genesis.json
// file of c.
func (c *Chain) writeConsumerGenesis() error {
	var (
		providerClientState = &ibctmtypes.ClientState{
			ChainId:         "provider",
			TrustLevel:      ibctmtypes.DefaultTrustLevel,
			TrustingPeriod:  time.Hour * 64,
			UnbondingPeriod: time.Hour * 128,
			MaxClockDrift:   time.Minute * 5,
		}
		providerConsState = &ibctmtypes.ConsensusState{
			Timestamp: time.Now().Add(time.Hour * 24),
			Root: commitmenttypes.NewMerkleRoot(
				[]byte("LpGpeyQVLUo9HpdsgJr12NP2eCICspcULiWa5u9udOA="),
			),
			NextValidatorsHash: []byte("E30CE736441FB9101FADDAF7E578ABBE6DFDB67207112350A9A904D554E1F5BE"),
		}
		params = ccvtypes.NewParams(
			true,
			1000, // ignore distribution
			"",   // ignore distribution
			"",   // ignore distribution
			ccvtypes.DefaultCCVTimeoutPeriod,
			ccvtypes.DefaultTransferTimeoutPeriod,
			ccvtypes.DefaultConsumerRedistributeFrac,
			ccvtypes.DefaultHistoricalEntries,
			ccvtypes.DefaultConsumerUnbondingPeriod,
			"0", // disable soft opt-out
			[]string{},
			[]string{},
		)
	)
	// Load public key from priv_validator_key.json file
	pvKeyFile, err := c.PrivValidatorKeyPath()
	if err != nil {
		return err
	}
	filePV := cmprivval.LoadFilePVEmptyState(pvKeyFile, "")
	pk, err := filePV.GetPubKey()
	if err != nil {
		return err
	}
	// Feed initial_val_set with this public key
	// Like for sovereign chain, provide only a single validator.
	valUpdates := cmtypes.ValidatorUpdates{
		cmtypes.UpdateValidator(pk.Bytes(), 1, pk.Type()),
	}
	// Build consumer genesis
	consumerGen := ccvtypes.NewInitialGenesisState(providerClientState, providerConsState, valUpdates, params)
	// Read genesis file
	genPath, err := c.GenesisPath()
	if err != nil {
		return err
	}
	genState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genPath)
	if err != nil {
		return err
	}
	// Update consumer module gen state
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	bz, err := codec.MarshalJSON(consumerGen)
	if err != nil {
		return err
	}
	genState[ccvconsumertypes.ModuleName] = bz
	// Update whole genesis
	bz, err = json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return err
	}
	genDoc.AppState = bz
	// Save genesis
	return genDoc.SaveAs(genPath)
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
// The check is performed by checking if the gentx dir exists in the config,
// unless c is a consumer chain, in that case the check relies on checking
// if the consumer genesis module is filled with validators.
func (c *Chain) IsInitialized() (bool, error) {
	home, err := c.Home()
	if err != nil {
		return false, err
	}
	cfg, err := c.Config()
	if err != nil {
		return false, err
	}
	if cfg.IsConsumerChain() {
		// Consumer chain doesn't have necessarily gentxs, so we can't rely on that
		// to determine if it's initialized. To perform that check, we need to
		// ensure the consumer genesis has InitialValSet filled.
		genPath, err := c.GenesisPath()
		if err != nil {
			return false, err
		}
		genState, _, err := genutiltypes.GenesisStateFromGenFile(genPath)
		if err != nil {
			// If the genesis isn't readable, don't propagate the error, just
			// consider the chain isn't initialized.
			return false, nil
		}
		var (
			consumerGenesis   ccvtypes.GenesisState
			interfaceRegistry = codectypes.NewInterfaceRegistry()
			codec             = codec.NewProtoCodec(interfaceRegistry)
		)
		err = codec.UnmarshalJSON(genState[ccvconsumertypes.ModuleName], &consumerGenesis)
		if err != nil {
			return false, err
		}
		return len(consumerGenesis.InitialValSet) != 0, nil
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
