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

type ChainCmd struct {
	appCmd         string
	chainID        string
	homeDir        string
	keyringBackend string
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
func WithKeyrinBackend(keyringBackend string) Option {
	return func(c *ChainCmd) {
		c.keyringBackend = keyringBackend
	}
}

// SetChainID sets the chain ID attached to commands
func (c *ChainCmd) SetChainID(chainID string) {
	c.chainID = chainID
}

// SetKeyringBackend sets the keyring backend attached to commands
func (c *ChainCmd) SetKeyringBackend(keyringBackend string) {
	c.keyringBackend = keyringBackend
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

// GentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) GentxCommand(
	validatorName string,
	selfDelegation string,
	moniker string,
	commissionRate string,
	commissionMaxRate string,
	commissionMaxChangeRate string,
	minSelfDelegation string,
	gasPrices string,
) step.Option {
	command := []string{
		commandGentx,
		validatorName,
		optionAmount,
		selfDelegation,
	}

	// Append optional validator information
	if moniker != "" {
		command = append(command, optionValidatorMoniker, moniker)
	}
	if commissionRate != "" {
		command = append(command, optionValidatorCommissionRate, commissionRate)
	}
	if commissionMaxRate != "" {
		command = append(command, optionValidatorCommissionMaxRate, commissionMaxRate)
	}
	if commissionMaxChangeRate != "" {
		command = append(command, optionValidatorCommissionMaxChangeRate, commissionMaxChangeRate)
	}
	if minSelfDelegation != "" {
		command = append(command, optionValidatorMinSelfDelegation, minSelfDelegation)
	}
	if gasPrices != "" {
		command = append(command, optionValidatorGasPrices, gasPrices)
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
		command = append(command, []string{optionKeyringBackend, c.keyringBackend}...)
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
