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
