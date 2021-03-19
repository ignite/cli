package networkbuilder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
)

type GentxInfo struct {
	ValidatorAddress string
	SelfDelegation   sdk.Coin
}

// VerifyProposals if proposals are correct and simulate them with the current launch information
// Correctness means checks that have to be performed off-chain
func (b *Builder) VerifyProposals(ctx context.Context, chainID string, homeDir string, proposals []int, commandOut io.Writer) (bool, string, error) {

	// Check all proposal
	for _, id := range proposals {
		proposal, err := b.ProposalGet(ctx, chainID, id)
		if err != nil {
			return false, "", err
		}

		// If this is a add validator proposal
		if proposal.Validator != nil {
			valAddress := proposal.Validator.ValidatorAddress
			selfDelegation := proposal.Validator.SelfDelegation

			// Check values inside the gentx are correct
			gentxInfo, err := parseGentx(proposal.Validator.Gentx)
			if err != nil {
				return false, "", err
			}

			// Check validator address
			if valAddress != gentxInfo.ValidatorAddress {
				return false, fmt.Sprintf(
					"proposal %v contains a validator address %v that doesn't match the one inside the gentx %v",
					id,
					valAddress,
					gentxInfo.ValidatorAddress,
				), nil
			}

			// Check self delagation
			if !selfDelegation.IsEqual(gentxInfo.SelfDelegation) {
				return false, fmt.Sprintf(
					"proposal %v contains a self delegation %v that doesn't match the one inside the gentx %v",
					id,
					selfDelegation.String(),
					gentxInfo.SelfDelegation.String(),
				), nil
			}
		}
	}

	// If all proposals are correct, simulate them
	return b.SimulateProposals(ctx, chainID, homeDir, proposals, commandOut)
}

type LaunchpadGentx struct {
	Value struct {
		Msg []struct {
			Value struct {
				ValidatorAddress string `json:"validator_address"`
				Value            struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"value"`
			} `json:"valmue"`
		} `json:"msg"`
	} `json:"value"`
}

type StargateGentx struct {
	Body struct {
		Messages []struct {
			ValidatorAddress string `json:"validator_address"`
			Value            struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"value"`
		} `json:"messages"`
	} `json:"body"`
}

func parseGentx(gentx jsondoc.Doc) (info GentxInfo, err error) {
	// Try parsing Launchpad gentx
	var launchpadGentx LaunchpadGentx
	if err := json.Unmarshal(gentx, &launchpadGentx); err != nil {
		return info, err
	}

	// Try parsing Stargate gentx
	var stargateGentx StargateGentx
	if err := json.Unmarshal(gentx, &stargateGentx); err != nil {
		return info, err
	}
	info.ValidatorAddress = stargateGentx.Body.Messages[0].ValidatorAddress
	amount, ok := sdk.NewIntFromString(stargateGentx.Body.Messages[0].Value.Amount)
	if !ok {
		return info, errors.New("the self-delegation inside the gentx is invalid")
	}
	info.SelfDelegation = sdk.NewCoin(
		stargateGentx.Body.Messages[0].Value.Denom,
		amount,
	)

	// Unrecognized gentx
	return info, errors.New("the gentx cannot be parsed")
}
