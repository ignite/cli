package chaincmd

import "github.com/tendermint/starport/starport/pkg/cmdrunner/step"

// stargateGentxCommand returns the command to generate a gentx for the chain
func (c ChainCmd) stargateGentxCommand(
	validatorName string,
	selfDelegation string,
	options ...GentxOption,
) step.Option {
	command := []string{
		commandGentx,
		validatorName,
		selfDelegation,
	}

	// Apply the options provided by the user
	for _, applyOption := range options {
		command = applyOption(command)
	}

	// Add necessary flags
	command = c.attachChainID(command)
	command = c.attachKeyringBackend(command)

	return c.daemonCommand(command)
}
