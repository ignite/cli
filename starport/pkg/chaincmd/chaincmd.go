package chaincmd

const (
	commandStart = "start"
	commandInit = "init"
	commandKeys = "keys"
	commandAddGenesisAccount = "add-genesis-account"
	commandGentx = "gentx"
	commandCollectGentxs = "collect-gentxs"
	commandValidateGenesis = "validate-genesis"
	commandShowNodeID = "show-node-id"

	optionHome = "--home"
	optionKeyringBackend = "--keyring-backend"
	optionChainID = "--chain-id"
	optionOutput = "--output"
	optionRecover = "--recover"
	optionAddress = "--address"
	optionName = "--name"
	optionAmount = "--amount"
	optionValidatorMoniker = "--moniker"
	optionValidatorCommissionRate = "--commission-rate"
	optionValidatorCommissionMaxRate = "commission-max-rate"
	optionValidatorCommissionMaxChangeRate = "--commission-max-change-rate"
	optionValidatorMinSelfDelegation = "--min-self-delegation"
	optionValidatorGasPrices = "--gas-prices"

	constTendermint = "tendermint"
	constJSON = "json"
)

type ChainCmd struct {
	appCmd string
	chainID string
	homeDir   string
	keyringBackend string
}

// NewChainCmd creates a new chaincmd to launch comand with the chain app
func NewChainCmd(appName string, chainID string, homeDir string, keyringBackend string) ChainCmd {
	return ChainCmd{
		appCmd: appName + "d",
		chainID: chainID,
		homeDir: homeDir,
		keyringBackend: keyringBackend,
	}
}

// StartCommand returns the command to start the daemon of the chain
func (c ChainCmd) StartCommand(options... string) []string {
	command := append([]string{
		c.appCmd,
		commandStart,
	}, options...)
	return c.withFlags(command)
}

// InitCommand returns the command to initialize the chain
func (c ChainCmd) InitCommand(moniker string, chainID string) []string {
	command := []string{
		c.appCmd,
		commandInit,
		moniker,
	}
	command = c.withChainID(command)
	return c.withFlags(command)
}

// AddKeyCommand returns the command to add a new key in the chain keyring
func (c ChainCmd) AddKeyCommand(accountName string) []string {
	command := []string{
		c.appCmd,
		commandKeys,
		"add",
		accountName,
		optionOutput,
		constJSON,
	}
	return c.withFlags(command)
}

// ImportKeyCommand returns the command to import a key into the chain keyring from a mnemonic
func (c ChainCmd) ImportKeyCommand(accountName string) []string {
	command := []string{
		c.appCmd,
		commandKeys,
		"add",
		accountName,
		optionRecover,
	}
	return c.withFlags(command)
}

// ShowKeyAddressCommand returns the command to print the address of a key in the chain keyring
func (c ChainCmd) ShowKeyAddressCommand(accountName string) []string {
	command := []string{
		c.appCmd,
		commandKeys,
		"show",
		accountName,
		optionAddress,
	}
	return c.withFlags(command)
}

// AddGenesisAccountCommand returns the command to add a new account in the genesis file of the chain
func (c ChainCmd) AddGenesisAccountCommand(address string, coins string) []string {
	command := []string{
		c.appCmd,
		commandAddGenesisAccount,
		address,
		coins,
	}
	return c.withFlags(command)
}

// GentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) GentxCommand(
	validatorName string,
	chainID string,
	selfDelegation string,
	moniker string,
	commissionRate string,
	commissionMaxRate string,
	commissionMaxChangeRate string,
	minSelfDelegation string,
	gasPrices string,
	) []string {
	command := []string{
		c.appCmd,
		commandGentx,
		optionName,
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

	command = c.withChainID(command)
	return c.withFlags(command)
}

// CollectGentxsCommand returns the command to gather the gentxs in /gentx dir into the genesis file of the chain
func (c ChainCmd) CollectGentxsCommand() []string {
	command := []string{
		c.appCmd,
		commandCollectGentxs,
	}
	return c.withFlags(command)
}

// ValidateGenesisCommand returns the command to check the validity of the chain genesis
func (c ChainCmd) ValidateGenesisCommand() []string {
	command := []string{
		c.appCmd,
		commandValidateGenesis,
	}
	return c.withFlags(command)
}

// ShowNodeIDCommand returns the command to print the node ID of the node for the chain
func (c ChainCmd) ShowNodeIDCommand() []string {
	command := []string{
		c.appCmd,
		constTendermint,
		commandShowNodeID,
	}
	return c.withFlags(command)
}

// withChainID appends the chain ID flag to the provided command
func (c ChainCmd) withChainID(command []string) []string {
	return append(command, []string{optionChainID, c.chainID}...)
}

// withFlags appends the global flags defined for the chain commands to the provided command
func (c ChainCmd) withFlags(command []string) []string {
	// Attach home
	if c.homeDir != "" {
		command = append(command, []string{optionHome, c.homeDir}...)
	}
	// Attach keyring backend
	if c.keyringBackend != "" {
		command = append(command, []string{optionKeyringBackend, c.keyringBackend}...)
	}

	return command
}

