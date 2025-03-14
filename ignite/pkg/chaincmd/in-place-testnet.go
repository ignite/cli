package chaincmd

import (
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
)

type InPlaceOption func([]string) []string

func InPlaceWithPrvKey(prvKey string) InPlaceOption {
	return func(s []string) []string {
		if len(prvKey) > 0 {
			return append(s, optionValidatorPrivateKey, prvKey)
		}
		return s
	}
}

func InPlaceWithAccountToFund(accounts string) InPlaceOption {
	return func(s []string) []string {
		if len(accounts) > 0 {
			return append(s, optionAccountToFund, accounts)
		}
		return s
	}
}

func InPlaceWithSkipConfirmation() InPlaceOption {
	return func(s []string) []string {
		return append(s, optionSkipConfirmation)
	}
}

// TestnetInPlaceCommand return command to start testnet in-place.
func (c ChainCmd) TestnetInPlaceCommand(newChainID, newOperatorAddress string, options ...InPlaceOption) step.Option {
	command := []string{
		commandTestnetInPlace,
		newChainID,
		newOperatorAddress,
	}

	// Apply the options provided by the user
	for _, apply := range options {
		command = apply(command)
	}

	return c.daemonCommand(command)
}

// Options for testnet multi node.
type MultiNodeOption func([]string) []string

// MultiNodeWithChainID returns a MultiNodeOption that appends the chainID option
// to the provided slice of strings.
func MultiNodeWithChainID(chainID string) MultiNodeOption {
	return func(s []string) []string {
		if len(chainID) > 0 {
			return append(s, optionChainID, chainID)
		}
		return s
	}
}

// MultiNodeWithDirOutput returns a MultiNodeOption that appends the output directory option
// to the provided slice of strings.
func MultiNodeWithDirOutput(dirOutput string) MultiNodeOption {
	return func(s []string) []string {
		if len(dirOutput) > 0 {
			return append(s, optionOutPutDir, dirOutput)
		}
		return s
	}
}

// MultiNodeWithNumValidator returns a MultiNodeOption that appends the number of validators option
// to the provided slice of strings.
func MultiNodeWithNumValidator(numVal string) MultiNodeOption {
	return func(s []string) []string {
		if len(numVal) > 0 {
			return append(s, optionNumValidator, numVal)
		}
		return s
	}
}

// MultiNodeWithValidatorsStakeAmount returns a MultiNodeOption that appends the stake amounts option
// to the provided slice of strings.
func MultiNodeWithValidatorsStakeAmount(satkeAmounts string) MultiNodeOption {
	return func(s []string) []string {
		if len(satkeAmounts) > 0 {
			return append(s, optionAmountStakes, satkeAmounts)
		}
		return s
	}
}

// MultiNodeDirPrefix returns a MultiNodeOption that appends the node directory prefix option
// to the provided slice of strings.
func MultiNodeDirPrefix(nodeDirPrefix string) MultiNodeOption {
	return func(s []string) []string {
		if len(nodeDirPrefix) > 0 {
			return append(s, optionNodeDirPrefix, nodeDirPrefix)
		}
		return s
	}
}

func MultiNodePorts(ports string) MultiNodeOption {
	return func(s []string) []string {
		if len(ports) > 0 {
			return append(s, optionPorts, ports)
		}
		return s
	}
}

// TestnetMultiNodeCommand return command to start testnet multinode.
func (c ChainCmd) TestnetMultiNodeCommand(options ...MultiNodeOption) step.Option {
	command := []string{
		commandTestnetMultiNode,
	}

	// Apply the options provided by the user
	for _, apply := range options {
		command = apply(command)
	}

	return c.daemonCommand(command)
}
