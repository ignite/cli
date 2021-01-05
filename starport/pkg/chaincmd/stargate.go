package chaincmd

import "github.com/tendermint/starport/starport/pkg/cmdrunner/step"

// stargateAddKeyCommand returns the command to add a new key in the chain keyring
func (c ChainCmd) stargateAddKeyCommand(accountName string) step.Option {
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

// stargateImportKeyCommand returns the command to import a key into the chain keyring from a mnemonic
func (c ChainCmd) stargateImportKeyCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"add",
		accountName,
		optionRecover,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// stargateShowKeyAddressCommand returns the command to print the address of a key in the chain keyring
func (c ChainCmd) stargateShowKeyAddressCommand(accountName string) step.Option {
	command := []string{
		commandKeys,
		"show",
		accountName,
		optionAddress,
	}
	command = c.attachKeyringBackend(command)
	return step.Exec(c.appCmd, c.attachHome(command)...)
}

// stargateGentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) stargateGentxCommand(
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
