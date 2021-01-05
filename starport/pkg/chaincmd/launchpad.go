package chaincmd

import "github.com/tendermint/starport/starport/pkg/cmdrunner/step"

const (
	commandConfig     = "config"
	commandRestServer = "rest-server"

	optionUnsafeCors = "--unsafe-cors"
	optionAPIAddress = "--laddr"
	optionRPCAddress = "--node"
	optionName       = "--name"
)

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

// SetHome sets a new home for the commands
func (c *ChainCmd) SetLaunchpadCLIHome(cliHome string) {
	c.cliHome = cliHome
}

// LaunchpadAddKeyCommand returns the command to add a new key in the chain keyring with Launchpad chains
func (c ChainCmd) LaunchpadAddKeyCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"add",
		accountName,
		optionOutput,
		constJSON,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.cliCmd, c.attachCLIHome(command)...)
}

// LaunchpadImportKeyCommand returns the command to import a key into the chain keyring from a mnemonic with Launchpad chains
func (c ChainCmd) LaunchpadImportKeyCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"add",
		accountName,
		optionRecover,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.cliCmd, c.attachCLIHome(command)...)
}

// LaunchpadShowKeyAddressCommand returns the command to print the address of a key in the chain keyring with Launchpad chains
func (c ChainCmd) LaunchpadShowKeyAddressCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"show",
		accountName,
		optionAddress,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.cliCmd, c.attachCLIHome(command)...)
}

// LaunchpadSetConfigCommand
func (c ChainCmd) LaunchpadSetConfigCommand(name string, value string) step.Option {
	command := []string{
		commandConfig,
		name,
		value,
	}
	return step.Exec(c.cliCmd, c.attachCLIHome(command)...)
}

// LaunchpadRestServerCommand
func (c ChainCmd) LaunchpadRestServerCommand(apiAddress string, rpcAddress string) step.Option {
	command := []string{
		commandRestServer,
		optionUnsafeCors,
		optionAPIAddress,
		apiAddress,
		optionRPCAddress,
		rpcAddress,
	}
	return step.Exec(c.cliCmd, c.attachCLIHome(command)...)
}

// LaunchpadGentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) LaunchpadGentxCommand(
	validatorName string,
	selfDelegation string,
	options ...GentxOption,
) step.Option {
	command := []string{
		commandGentx,
		optionName,
		validatorName,
		optionAmount,
		selfDelegation,
	}

	// Apply the options provided by the user
	for _, applyOption := range options {
		command = applyOption(command)
	}

	command = c.attachKeyringBackend(command)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// attachCLIHome appends the home flag to the provided CLI command
func (c ChainCmd) attachCLIHome(command []string) []string {
	if c.cliHome != "" {
		command = append(command, []string{optionHome, c.cliHome}...)
	}
	return command
}
