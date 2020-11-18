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

// ProposeAddAccount proposes to add an account to chain.
func (b *Builder) ProposeAddAccount(ctx context.Context, chainID string, account spn.ProposalAddAccount) error {
	acc, err := b.AccountInUse()
	if err != nil {
		return err
	}
	return b.spnclient.ProposeAddAccount(ctx, acc.Name, chainID, account)
}

// ProposeAddValidator proposes to add a validator to chain.
func (b *Builder) ProposeAddValidator(ctx context.Context, chainID string, validator spn.ProposalAddValidator) error {
	acc, err := b.AccountInUse()
	if err != nil {
		return err
	}
	return b.spnclient.ProposeAddValidator(ctx, acc.Name, chainID, validator)
}

// ProposalApprove approves a proposal by id.
func (b *Builder) ProposalApprove(ctx context.Context, chainID string, id int) error {
	acc, err := b.AccountInUse()
	if err != nil {
		return err
	}
	return b.spnclient.ProposalApprove(ctx, acc.Name, chainID, id)
}

// ProposalReject rejects a proposal by id.
func (b *Builder) ProposalReject(ctx context.Context, chainID string, id int) error {
	acc, err := b.AccountInUse()
	if err != nil {
		return err
	}
	return b.spnclient.ProposalReject(ctx, acc.Name, chainID, id)

}
