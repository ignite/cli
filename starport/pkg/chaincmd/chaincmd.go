package chaincmd

import "github.com/tendermint/starport/starport/pkg/cmdrunner/step"

const (
	commandStart             = "start"
	commandInit              = "init"
	commandKeys              = "keys"
	commandAddGenesisAccount = "add-genesis-account"
	commandGentx             = "gentx"
	commandCollectGentxs     = "collect-gentxs"
	commandValidateGenesis   = "validate-genesis"
	commandShowNodeID        = "show-node-id"

	optionHome                             = "--home"
	optionKeyringBackend                   = "--keyring-backend"
	optionChainID                          = "--chain-id"
	optionOutput                           = "--output"
	optionRecover                          = "--recover"
	optionAddress                          = "--address"
	optionAmount                           = "--amount"
	optionValidatorMoniker                 = "--moniker"
	optionValidatorCommissionRate          = "--commission-rate"
	optionValidatorCommissionMaxRate       = "--commission-max-rate"
	optionValidatorCommissionMaxChangeRate = "--commission-max-change-rate"
	optionValidatorMinSelfDelegation       = "--min-self-delegation"
	optionValidatorGasPrices               = "--gas-prices"

	constTendermint = "tendermint"
	constJSON       = "json"
)

type KeyringBackend string

const (
	KeyringBackendOS      KeyringBackend = "os"
	KeyringBackendFile    KeyringBackend = "file"
	KeyringBackendPass    KeyringBackend = "pass"
	KeyringBackendTest    KeyringBackend = "test"
	KeyringBackendKwallet KeyringBackend = "kwallet"
)

type ChainCmd struct {
	appCmd         string
	chainID        string
	homeDir        string
	keyringBackend KeyringBackend
	cliCmd         string
	cliHome        string
}

// New creates a new ChainCmd to launch command with the chain app
func New(appCmd string, options ...Option) ChainCmd {
	chainCmd := ChainCmd{
		appCmd: appCmd,
	}

	// Apply the options provided by the user
	for _, applyOption := range options {
		applyOption(&chainCmd)
	}

	return chainCmd
}

type Option func(*ChainCmd)

// WithHome replaces the default home used by the chain
func WithHome(home string) Option {
	return func(c *ChainCmd) {
		c.homeDir = home
	}
}

// WithChainID provides a specific chain ID for the commands that accept this option
func WithChainID(chainID string) Option {
	return func(c *ChainCmd) {
		c.chainID = chainID
	}
}

// WithKeyrinBackend provides a specific keyring backend for the commands that accept this option
func WithKeyrinBackend(keyringBackend KeyringBackend) Option {
	return func(c *ChainCmd) {
		c.keyringBackend = keyringBackend
	}
}

// WithLaunchpadCLI provides the name of the CLI application to call Launchpad CLI commands
func WithLaunchpadCLI(cliCmd string) Option {
	return func(c *ChainCmd) {
		c.cliCmd = cliCmd
	}
}

// WithLaunchpadCLIHome replaces the default home used by the Launchpad chain CLI
func WithLaunchpadCLIHome(cliHome string) Option {
	return func(c *ChainCmd) {
		c.cliHome = cliHome
	}
}

// StartCommand returns the command to start the daemon of the chain
func (c ChainCmd) StartCommand(options ...string) step.Option {
	command := append([]string{
		commandStart,
	}, options...)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// InitCommand returns the command to initialize the chain
func (c ChainCmd) InitCommand(moniker string) step.Option {
	command := []string{
		commandInit,
		moniker,
	}
	command = c.attachChainID(command)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// AddKeyCommand returns the command to add a new key in the chain keyring
func (c ChainCmd) AddKeyCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"add",
		accountName,
		optionOutput,
		constJSON,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// ImportKeyCommand returns the command to import a key into the chain keyring from a mnemonic
func (c ChainCmd) ImportKeyCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"add",
		accountName,
		optionRecover,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// ShowKeyAddressCommand returns the command to print the address of a key in the chain keyring
func (c ChainCmd) ShowKeyAddressCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"show",
		accountName,
		optionAddress,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// AddGenesisAccountCommand returns the command to add a new account in the genesis file of the chain
func (c ChainCmd) AddGenesisAccountCommand(address string, coins string) step.Option {
	command := []string{
		commandAddGenesisAccount,
		address,
		coins,
	}
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// Options for the GentxCommand
type GentxOption func([]string) []string

// GentxWithMoniker provides moniker option for the gentx command
func GentxWithMoniker(moniker string) GentxOption {
	return func(command []string) []string {
		if len(moniker) > 0 {
			return append(command, optionValidatorMoniker, moniker)
		}
		return command
	}
}

// GentxWithCommissionRate provides commission rate option for the gentx command
func GentxWithCommissionRate(commissionRate string) GentxOption {
	return func(command []string) []string {
		if len(commissionRate) > 0 {
			return append(command, optionValidatorCommissionRate, commissionRate)
		}
		return command
	}
}

// GentxWithCommissionMaxRate provides commission max rate option for the gentx command
func GentxWithCommissionMaxRate(commissionMaxRate string) GentxOption {
	return func(command []string) []string {
		if len(commissionMaxRate) > 0 {
			return append(command, optionValidatorCommissionMaxRate, commissionMaxRate)
		}
		return command
	}
}

// GentxWithCommissionMaxChangeRate provides commission max change rate option for the gentx command
func GentxWithCommissionMaxChangeRate(commissionMaxChangeRate string) GentxOption {
	return func(command []string) []string {
		if len(commissionMaxChangeRate) > 0 {
			return append(command, optionValidatorCommissionMaxChangeRate, commissionMaxChangeRate)
		}
		return command
	}
}

// GentxWithMinSelfDelegation provides minimum self delegation option for the gentx command
func GentxWithMinSelfDelegation(minSelfDelegation string) GentxOption {
	return func(command []string) []string {
		if len(minSelfDelegation) > 0 {
			return append(command, optionValidatorMinSelfDelegation, minSelfDelegation)
		}
		return command
	}
}

// GentxWithGasPrices provides gas price option for the gentx command
func GentxWithGasPrices(gasPrices string) GentxOption {
	return func(command []string) []string {
		if len(gasPrices) > 0 {
			return append(command, optionValidatorGasPrices, gasPrices)
		}
		return command
	}
}

// GentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) GentxCommand(
	validatorName string,
	selfDelegation string,
	options ...GentxOption,
) step.Option {
	command := []string{
		commandGentx,
		validatorName,
		optionAmount,
		selfDelegation,
	}

	// Apply the options provided by the user
	for _, applyOption := range options {
		command = applyOption(command)
	}

	// Add necessary flags
	command = c.attachChainID(command)
	command = c.attachKeyringBackend(command)

	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// CollectGentxsCommand returns the command to gather the gentxs in /gentx dir into the genesis file of the chain
func (c ChainCmd) CollectGentxsCommand() step.Option {
	command := []string{
		commandCollectGentxs,
	}
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// ValidateGenesisCommand returns the command to check the validity of the chain genesis
func (c ChainCmd) ValidateGenesisCommand() step.Option {
	command := []string{
		commandValidateGenesis,
	}
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// ShowNodeIDCommand returns the command to print the node ID of the node for the chain
func (c ChainCmd) ShowNodeIDCommand() step.Option {
	command := []string{
		constTendermint,
		commandShowNodeID,
	}
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// attachChainID appends the chain ID flag to the provided command
func (c ChainCmd) attachChainID(command []string) []string {
	if c.chainID != "" {
		command = append(command, []string{optionChainID, c.chainID}...)
	}
	return command
}

// attachKeyringBackend appends the keyring backend flag to the provided command
func (c ChainCmd) attachKeyringBackend(command []string) []string {
	if c.keyringBackend != "" {
		command = append(command, []string{optionKeyringBackend, string(c.keyringBackend)}...)
	}
	return command
}

// attachHome appends the home flag to the provided command
func (c ChainCmd) attachHome(command []string) []string {
	if c.homeDir != "" {
		command = append(command, []string{optionHome, c.homeDir}...)
	}
	return command
}
