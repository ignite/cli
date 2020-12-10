package networkbuilder

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/spn"
)

// ShowChain shows details of a chain.
func (b *Builder) ShowChain(ctx context.Context, chainID string) (spn.Chain, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return spn.Chain{}, err
	}
	return b.spnclient.ShowChain(ctx, account.Name, chainID)
}

// LaunchInformation retrieves chain's launch information.
func (b *Builder) LaunchInformation(ctx context.Context, chainID string) (spn.LaunchInformation, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return spn.LaunchInformation{}, err
	}
	return b.spnclient.LaunchInformation(ctx, account.Name, chainID)
}

// SimulatedLaunchInformation retrieves chain's simulated launch information from proposals to verify.
func (b *Builder) SimulatedLaunchInformation(ctx context.Context, chainID string, proposalIDs []int) (spn.LaunchInformation, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return spn.LaunchInformation{}, err
	}
	return b.spnclient.SimulatedLaunchInformation(ctx, account.Name, chainID, proposalIDs)
}

// ChainShow lists summaries of chains
func (b *Builder) ChainList(ctx context.Context) ([]spn.ChainSummary, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return []spn.ChainSummary{}, err
	}
	return b.spnclient.ChainList(ctx, account.Name)
}
