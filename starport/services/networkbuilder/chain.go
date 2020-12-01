package networkbuilder

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/spn"
)

// ChainShow shows details of a chain.
func (b *Builder) ChainShow(ctx context.Context, chainID string) (spn.ChainInformation, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return spn.ChainInformation{}, err
	}
	return b.spnclient.GetChainInformation(ctx, account.Name, chainID)
}
