package networkbuilder

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/spn"
)

// ProposalList lists proposals on a chain by status.
func (b *Builder) ProposalList(ctx context.Context, chainID string, status spn.ProposalStatus) ([]spn.Proposal, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return nil, err
	}
	return b.spnclient.ProposalList(ctx, account.Name, chainID, status)
}

// ProposalGet retrieves a proposal on a chain by id.
func (b *Builder) ProposalGet(ctx context.Context, chainID string, id int) (spn.Proposal, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return spn.Proposal{}, err
	}
	return b.spnclient.ProposalGet(ctx, account.Name, chainID, id)
}

// Propose proposes given proposals in batch for chainID by using SPN accountName.
func (b *Builder) Propose(ctx context.Context, chainID string, proposals ...spn.ProposalOption) error {
	acc, err := b.AccountInUse()
	if err != nil {
		return err
	}
	return b.spnclient.Propose(ctx, acc.Name, chainID, proposals...)
}

// SubmitReviewals submits reviewals for proposals in batch for chainID by using SPN accountName.
func (b *Builder) SubmitReviewals(ctx context.Context, chainID string, reviewals ...spn.Reviewal) error {
	acc, err := b.AccountInUse()
	if err != nil {
		return err
	}
	return b.spnclient.SubmitReviewals(ctx, acc.Name, chainID, reviewals...)
}
