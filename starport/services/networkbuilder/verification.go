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

type VerificationError struct {
	Err error
}

func (e VerificationError) Error() string {
	return e.Err.Error()
}

type GentxInfo struct {
	DelegatorAddress, ValidatorAddress string
	SelfDelegation                     sdk.Coin
}

// VerifyProposals if proposals are correct and simulate them with the current launch information
// Correctness means checks that have to be performed off-chain
func (b *Builder) VerifyProposals(ctx context.Context, chainID string, proposals []int, commandOut io.Writer) error {

	// Check all proposal
	for _, id := range proposals {
		proposal, err := b.ProposalGet(ctx, chainID, id)
		if err != nil {
			return err
		}

		// If this is a add validator proposal
		if proposal.Validator != nil {
			valAddress := proposal.Validator.ValidatorAddress
			selfDelegation := proposal.Validator.SelfDelegation

			// Check values inside the gentx are correct
			gentxInfo, err := ParseGentx(proposal.Validator.Gentx)
			if err != nil {
				return VerificationError{
					fmt.Errorf("cannot parse proposal %v gentx: %v", id, err.Error()),
				}
			}

			// Check validator address
			if valAddress != gentxInfo.DelegatorAddress {
				return VerificationError{
					fmt.Errorf(
						"proposal %v contains a validator address %v that doesn't match the one inside the gentx: %v",
						id,
						valAddress,
						gentxInfo.DelegatorAddress,
					),
				}
			}

			// Check self delagation
			if !selfDelegation.IsEqual(gentxInfo.SelfDelegation) {
				return VerificationError{
					fmt.Errorf(
						"proposal %v contains a self delegation %v that doesn't match the one inside the gentx: %v",
						id,
						selfDelegation.String(),
						gentxInfo.SelfDelegation.String(),
					),
				}
			}
		}
	}

	// If all proposals are correct, simulate them
	return b.SimulateProposals(ctx, chainID, proposals, commandOut)
}

type stargateGentx struct {
	Body struct {
		Messages []struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Value            struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"value"`
		} `json:"messages"`
	} `json:"body"`
}

func ParseGentx(gentx jsondoc.Doc) (info GentxInfo, err error) {
	// Try parsing Stargate gentx
	var stargateGentx stargateGentx
	if err := json.Unmarshal(gentx, &stargateGentx); err != nil {
		return info, err
	}
	if stargateGentx.Body.Messages == nil {
		return info, errors.New("the gentx cannot be parsed")
	}

	// This is a stargate gentx
	if len(stargateGentx.Body.Messages) != 1 {
		return info, errors.New("add validator gentx must contain 1 message")
	}
	info.DelegatorAddress = stargateGentx.Body.Messages[0].DelegatorAddress
	info.ValidatorAddress = stargateGentx.Body.Messages[0].ValidatorAddress
	amount, ok := sdk.NewIntFromString(stargateGentx.Body.Messages[0].Value.Amount)
	if !ok {
		return info, errors.New("the self-delegation inside the gentx is invalid")
	}
	info.SelfDelegation = sdk.NewCoin(
		stargateGentx.Body.Messages[0].Value.Denom,
		amount,
	)

	return info, nil
}
