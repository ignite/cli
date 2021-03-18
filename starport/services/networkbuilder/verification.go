package networkbuilder

import (
	"context"
	"io"
)

// VerifyProposals if proposals are correct and simulate them with the current launch information
// Correctness means checks that have to be performed off-chain
func (b *Builder) VerifyProposals(ctx context.Context, chainID string, homeDir string, proposals []int, commandOut io.Writer) (bool, error) {

	// Check all proposal
	//for _, id := range proposals {
	//	proposal, err := b.ProposalGet(ctx, chainID, id)
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	// If this is a add validator proposal
	//	if proposal.Validator != nil {
	//		var gentx interface{}
	//		valAddress := proposal.Validator.ValidatorAddress
	//		selfDelegation := proposal.Validator.SelfDelegation
	//
	//		if err := json.Unmarshal(proposal.Validator.Gentx, &gentx); err != nil {
	//			// gentx cannot be json unmarshaled
	//			return false, nil
	//		}
	//
	//
	//	}
	//}

	// If all proposals are correct, simulate them
	return b.SimulateProposals(ctx, chainID, homeDir, proposals, commandOut)
}