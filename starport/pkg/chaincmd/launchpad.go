package chaincmd

import "github.com/tendermint/starport/starport/pkg/cmdrunner/step"

const (
	commandConfig     = "config"
	commandRestServer = "rest-server"

	optionUnsafeCors = "--unsafe-cors"
	optionAPIAddress = "--laddr"
	optionRPCAddress = "--node"
	optionName 		= "--name"
)

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
	return step.Exec(c.cliCmd, c.attachHome(command)...)
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
	return step.Exec(c.cliCmd, c.attachHome(command)...)
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
	return step.Exec(c.cliCmd, c.attachHome(command)...)
}

// LaunchpadSetConfigCommand
func (c ChainCmd) LaunchpadSetConfigCommand(name string, value string) step.Option {
	command := []string{
		commandConfig,
		name,
		value,
	}
	return step.Exec(c.cliCmd, c.attachHome(command)...)
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
	command = c.attachKeyringBackend(command)
	return step.Exec(c.cliCmd, c.attachHome(command)...)
}

// LaunchpadGentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) LaunchpadGentxCommand(
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

	return step.Exec(c.appCmd, c.attachHome(command)...)
}
