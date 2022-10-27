package chaincmd

import "github.com/ignite/cli/ignite/pkg/cmdrunner/step"

const (
	commandConfig     = "config"
	commandRestServer = "rest-server"

	optionUnsafeCors = "--unsafe-cors"
	optionAPIAddress = "--laddr"
	optionRPCAddress = "--node"
	optionName       = "--name"
)

// launchpadSetConfigCommand
func (c ChainCmd) launchpadSetConfigCommand(name string, value string) step.Option {
	command := []string{
		commandConfig,
		name,
		value,
	}

	return c.cliCommand(command)
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
	return c.cliCommand(command)
}

// attachCLIHome appends the home flag to the provided CLI command
func (c ChainCmd) attachCLIHome(command []string) []string {
	if c.cliHome != "" {
		command = append(command, []string{optionHome, c.cliHome}...)
	}
	return command
}
