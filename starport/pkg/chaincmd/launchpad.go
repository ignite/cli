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

// WithLaunchpadCLIHome replaces the default home used by the Launchpad chain CLI
func WithLaunchpadCLIHome(cliHome string) Option {
	return func(c *ChainCmd) {
		c.cliHome = cliHome
	}
}

// launchpadGentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) launchpadGentxCommand(
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

// launchpadSetConfigCommand
func (c ChainCmd) launchpadSetConfigCommand(name string, value string) step.Option {
	command := []string{
		commandConfig,
		name,
		value,
	}
	return step.Exec(c.cliCmd, c.attachHome(command)...)
}

// launchpadRestServerCommand
func (c ChainCmd) launchpadRestServerCommand(apiAddress string, rpcAddress string) step.Option {
	command := []string{
		commandRestServer,
		optionUnsafeCors,
		optionAPIAddress,
		apiAddress,
		optionRPCAddress,
		rpcAddress,
	}
	return step.Exec(c.cliCmd, c.attachHome(command)...)
}

// attachCLIHome appends the home flag to the provided CLI command
func (c ChainCmd) attachCLIHome(command []string) []string {
	if c.cliHome != "" {
		command = append(command, []string{optionHome, c.cliHome}...)
	}
	return command
}
