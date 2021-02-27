package chaincmd

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xurl"
)

const (
	commandStart             = "start"
	commandInit              = "init"
	commandKeys              = "keys"
	commandAddGenesisAccount = "add-genesis-account"
	commandGentx             = "gentx"
	commandCollectGentxs     = "collect-gentxs"
	commandValidateGenesis   = "validate-genesis"
	commandShowNodeID        = "show-node-id"
	commandStatus            = "status"
	commandTx                = "tx"
	commandQuery             = "query"
	commandUnsafeReset       = "unsafe-reset-all"
	commandExport            = "export"

	optionHome                             = "--home"
	optionNode                             = "--node"
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
	optionYes                              = "--yes"
	optionHomeClient                       = "--home-client"

	constTendermint = "tendermint"
	constJSON       = "json"
)

type KeyringBackend string

const (
	KeyringBackendUnspecified KeyringBackend = ""
	KeyringBackendOS          KeyringBackend = "os"
	KeyringBackendFile        KeyringBackend = "file"
	KeyringBackendPass        KeyringBackend = "pass"
	KeyringBackendTest        KeyringBackend = "test"
	KeyringBackendKwallet     KeyringBackend = "kwallet"
)

type ChainCmd struct {
	appCmd          string
	chainID         string
	homeDir         string
	keyringBackend  KeyringBackend
	KeyringPassword string
	cliCmd          string
	cliHome         string
	nodeAddress     string

	isAutoChainIDDetectionEnabled bool

	sdkVersion cosmosver.Version
}

// New creates a new ChainCmd to launch command with the chain app
func New(appCmd string, options ...Option) ChainCmd {
	c := ChainCmd{
		appCmd:     appCmd,
		sdkVersion: cosmosver.Versions.Latest(),
	}

	applyOptions(&c, options)

	return c
}

// Copy makes a copy of ChainCmd by overwriting its options with given options.
func (c ChainCmd) Copy(options ...Option) ChainCmd {
	applyOptions(&c, options)

	return c
}

// Option configures ChainCmd.
type Option func(*ChainCmd)

func applyOptions(c *ChainCmd, options []Option) {
	for _, applyOption := range options {
		applyOption(c)
	}
}

// WithVersion sets the version of the blockchain.
// when this is not provided, latest version of SDK is assumed.
func WithVersion(v cosmosver.Version) Option {
	return func(c *ChainCmd) {
		c.sdkVersion = v
	}
}

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

// WithAutoChainIDDetection finds out the chain id by communicating with the node running.
func WithAutoChainIDDetection() Option {
	return func(c *ChainCmd) {
		c.isAutoChainIDDetectionEnabled = true
	}
}

// WithKeyringBackend provides a specific keyring backend for the commands that accept this option
func WithKeyringBackend(keyringBackend KeyringBackend) Option {
	return func(c *ChainCmd) {
		c.keyringBackend = keyringBackend
	}
}

// WithKeyringPassword provides a password to unlock keyring
func WithKeyringPassword(password string) Option {
	return func(c *ChainCmd) {
		c.KeyringPassword = password
	}
}

// WithNodeAddress sets the node address for the commands that needs to make an
// API request to the node that has a different node address other than the default one.
func WithNodeAddress(addr string) Option {
	return func(c *ChainCmd) {
		c.nodeAddress = xurl.TCP(addr)
	}
}

// WithLaunchpadCLI provides the CLI application name for the blockchain
// this is necessary for Launchpad applications since it has two different binaries but
// not needed by Stargate applications
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
	return c.daemonCommand(command)
}

// InitCommand returns the command to initialize the chain
func (c ChainCmd) InitCommand(moniker string) step.Option {
	command := []string{
		commandInit,
		moniker,
	}
	command = c.attachChainID(command)
	return c.daemonCommand(command)
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

	return c.cliCommand(command)
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

	return c.cliCommand(command)
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

	return c.cliCommand(command)
}

// ListKeysCommand returns the command to print the list of a keys in the chain keyring
func (c ChainCmd) ListKeysCommand() step.Option {
	command := []string{
		commandKeys,
		"list",
		optionOutput,
		constJSON,
	}
	command = c.attachKeyringBackend(command)

	return c.cliCommand(command)
}

// AddGenesisAccountCommand returns the command to add a new account in the genesis file of the chain
func (c ChainCmd) AddGenesisAccountCommand(address string, coins string) step.Option {
	command := []string{
		commandAddGenesisAccount,
		address,
		coins,
	}

	return c.daemonCommand(command)
}

// GentxOption for the GentxCommand
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

func (c ChainCmd) IsAutoChainIDDetectionEnabled() bool {
	return c.isAutoChainIDDetectionEnabled
}

func (c ChainCmd) SDKVersion() cosmosver.Version {
	return c.sdkVersion
}

// GentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) GentxCommand(
	validatorName string,
	selfDelegation string,
	options ...GentxOption,
) step.Option {
	command := []string{
		commandGentx,
	}

	if c.sdkVersion.Is(cosmosver.StargateZeroFourtyAndAbove) {
		command = append(command,
			validatorName,
			selfDelegation,
		)
	}

	if c.sdkVersion.Is(cosmosver.StargateBelowZeroFourty) {
		command = append(command,
			validatorName,
			optionAmount,
			selfDelegation,
		)
	}

	if c.sdkVersion.Is(cosmosver.LaunchpadAny) {
		command = append(command,
			optionName,
			validatorName,
			optionAmount,
			selfDelegation,
		)

		// Attach home client option
		if c.cliHome != "" {
			command = append(command, []string{optionHomeClient, c.cliHome}...)
		}
	}

	// Apply the options provided by the user
	for _, applyOption := range options {
		command = applyOption(command)
	}

	// Add necessary flags
	if c.sdkVersion.Major().Is(cosmosver.Stargate) {
		command = c.attachChainID(command)
	}

	command = c.attachKeyringBackend(command)

	return c.daemonCommand(command)
}

// CollectGentxsCommand returns the command to gather the gentxs in /gentx dir into the genesis file of the chain
func (c ChainCmd) CollectGentxsCommand() step.Option {
	command := []string{
		commandCollectGentxs,
	}
	return c.daemonCommand(command)
}

// ValidateGenesisCommand returns the command to check the validity of the chain genesis
func (c ChainCmd) ValidateGenesisCommand() step.Option {
	command := []string{
		commandValidateGenesis,
	}
	return c.daemonCommand(command)
}

// ShowNodeIDCommand returns the command to print the node ID of the node for the chain
func (c ChainCmd) ShowNodeIDCommand() step.Option {
	command := []string{
		constTendermint,
		commandShowNodeID,
	}
	return c.daemonCommand(command)
}

// UnsafeResetCommand returns the command to reset the blockchain database
func (c ChainCmd) UnsafeResetCommand() step.Option {
	command := []string{
		commandUnsafeReset,
	}
	return c.daemonCommand(command)
}

// ExportCommand returns the command to export the state of the blockchain into a genesis file
func (c ChainCmd) ExportCommand() step.Option {
	command := []string{
		commandExport,
	}
	return c.daemonCommand(command)
}

// BankSendCommand returns the command for transferring tokens.
func (c ChainCmd) BankSendCommand(fromAddress, toAddress, amount string) step.Option {
	command := []string{
		commandTx,
	}

	if c.sdkVersion.Major().Is(cosmosver.Stargate) {
		command = append(command,
			"bank",
		)
	}

	command = append(command,
		"send",
		fromAddress,
		toAddress,
		amount,
		optionYes,
	)

	command = c.attachChainID(command)
	command = c.attachKeyringBackend(command)
	command = c.attachNode(command)

	if c.sdkVersion.Major().Is(cosmosver.Launchpad) {
		command = append(command, optionOutput, constJSON)
	}

	return c.cliCommand(command)
}

// QueryTxEventsCommand returns the command to query events.
func (c ChainCmd) QueryTxEventsCommand(query string) step.Option {
	command := []string{
		commandQuery,
		"txs",
		"--events",
		query,
		"--page", "1",
		"--limit", "1000",
	}

	if c.sdkVersion.Major().Is(cosmosver.Launchpad) {
		command = append(command,
			"--trust-node",
		)
	}

	return c.cliCommand(command)
}

// LaunchpadSetConfigCommand returns the command to set config value
func (c ChainCmd) LaunchpadSetConfigCommand(name string, value string) step.Option {
	// Check version
	if c.isStargate() {
		panic("config command doesn't exist for Stargate")
	}
	return c.launchpadSetConfigCommand(name, value)
}

// LaunchpadRestServerCommand returns the command to start the CLI REST server
func (c ChainCmd) LaunchpadRestServerCommand(apiAddress string, rpcAddress string) step.Option {
	// Check version
	if c.isStargate() {
		panic("rest-server command doesn't exist for Stargate")
	}
	return c.launchpadRestServerCommand(apiAddress, rpcAddress)
}

// StatusCommand returns the command that fetches node's status.
func (c ChainCmd) StatusCommand() step.Option {
	command := []string{
		commandStatus,
	}

	return c.cliCommand(command)
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

// attacNode appends the node flag to the provided command
func (c ChainCmd) attachNode(command []string) []string {
	if c.nodeAddress != "" {
		command = append(command, []string{optionNode, c.nodeAddress}...)
	}
	return command
}

// isStargate checks if the version for commands is Stargate
func (c ChainCmd) isStargate() bool {
	return c.sdkVersion.Major() == cosmosver.Stargate
}

// daemonCommand returns the daemon command from the provided command
func (c ChainCmd) daemonCommand(command []string) step.Option {
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// cliCommand returns the cli command from the provided command
// cli is the daemon for Stargate
func (c ChainCmd) cliCommand(command []string) step.Option {
	// Check version
	if c.isStargate() {
		return step.Exec(c.appCmd, c.attachHome(command)...)
	}
	return step.Exec(c.cliCmd, c.attachCLIHome(command)...)
}

// KeyringBackendFromString returns the keyring backend from its string
func KeyringBackendFromString(kb string) (KeyringBackend, error) {
	existingKeyringBackend := map[KeyringBackend]bool{
		KeyringBackendUnspecified: true,
		KeyringBackendOS:          true,
		KeyringBackendFile:        true,
		KeyringBackendPass:        true,
		KeyringBackendTest:        true,
		KeyringBackendKwallet:     true,
	}

	if _, ok := existingKeyringBackend[KeyringBackend(kb)]; ok {
		return KeyringBackend(kb), nil
	}
	return KeyringBackendUnspecified, fmt.Errorf("unrecognized keyring backend: %s", kb)
}
