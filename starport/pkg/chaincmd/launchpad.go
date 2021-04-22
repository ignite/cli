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

// launchpadSetConfigCommand
func (chainCmd ChainCmd) launchpadSetConfigCommand(name string, value string) step.Option {
	command := []string{
		commandConfig,
		name,
		value,
	}

	return chainCmd.cliCommand(command)
}

// launchpadRestServerCommand
func (chainCmd ChainCmd) launchpadRestServerCommand(apiAddress string, rpcAddress string) step.Option {
	command := []string{
		commandRestServer,
		optionUnsafeCors,
		optionAPIAddress,
		apiAddress,
		optionRPCAddress,
		rpcAddress,
	}
	return chainCmd.cliCommand(command)
}

// attachCLIHome appends the home flag to the provided CLI command
func (chainCmd ChainCmd) attachCLIHome(command []string) []string {
	if chainCmd.cliHome != "" {
		command = append(command, []string{optionHome, chainCmd.cliHome}...)
	}
	return command
}
