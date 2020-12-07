package networkbuilder

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/spn"
)

// ChainShow shows details of a chain.
func (b *Builder) ChainShow(ctx context.Context, chainID string) (spn.Chain, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return spn.Chain{}, err
	}
	return b.spnclient.ChainGet(ctx, account.Name, chainID)
}

// ChainShow lists summaries of chains
func (b *Builder) ChainList(ctx context.Context) ([]spn.ChainSummary, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return []spn.ChainSummary{}, err
	}
	return b.spnclient.ChainList(ctx, account.Name)
}
