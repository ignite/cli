package chaincmd

import (
	"fmt"

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

type MultiNodeOption func([]string) []string

func MultiNodeWithChainID(ChainId string) MultiNodeOption {
	return func(s []string) []string {
		if len(ChainId) > 0 {
			return append(s, optionChainID, ChainId)
		}
		return s
	}
}

func MultiNodeWithDirOutput(dirOutput string) MultiNodeOption {
	return func(s []string) []string {
		if len(dirOutput) > 0 {
			return append(s, optionOutPutDir, dirOutput)
		}
		return s
	}
}

func MultiNodeWithNumValidator(numVal string) MultiNodeOption {
	return func(s []string) []string {
		if len(numVal) > 0 {
			return append(s, optionNumValidator, numVal)
		}
		return s
	}
}
func MultiNodeWithValidatorsStakeAmount(satkeAmounts string) MultiNodeOption {
	return func(s []string) []string {
		if len(satkeAmounts) > 0 {
			return append(s, optionAmountStakes, satkeAmounts)
		}
		return s
	}
}

func MultiNodeWithHome(home string) MultiNodeOption {
	return func(s []string) []string {
		if len(home) > 0 {
			return append(s, optionHome, home)
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
	fmt.Println(command)

	return c.daemonCommand(command)
}

// TestnetMultiNodeCommand return command to start testnet multinode.
func (c ChainCmd) TestnetStartMultiNodeCommand(options ...MultiNodeOption) step.Option {
	command := []string{
		commandStart,
	}

	// Apply the options provided by the user
	for _, apply := range options {
		command = apply(command)
	}

	return c.daemonCommandIncludedHomeFlag(command)
}
