package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
)

// ChainLaunch returns the chain launch data from SPN
func (b *Builder) ChainLaunch(ctx context.Context, launchID uint64) (launchtypes.Chain, error) {
	b.ev.Send(events.New(events.StatusOngoing, "Fetching chain launch data"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).Chain(ctx, &launchtypes.QueryGetChainRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return launchtypes.Chain{}, err
	}
	return res.Chain, err
}

// GenesisAccounts returns the list of approved genesis account for a launch from SPN
func (b *Builder) GenesisAccounts(ctx context.Context, launchID uint64) ([]launchtypes.GenesisAccount, error) {
	b.ev.Send(events.New(events.StatusOngoing, "Fetching genesis accounts"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).GenesisAccountAll(ctx, &launchtypes.QueryAllGenesisAccountRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return []launchtypes.GenesisAccount{}, err
	}
	return res.GenesisAccount, err
}

// VestingAccounts returns the list of approved genesis vesting account for a launch from SPN
func (b *Builder) VestingAccounts(ctx context.Context, launchID uint64) ([]launchtypes.VestingAccount, error) {
	b.ev.Send(events.New(events.StatusOngoing, "Fetching genesis vesting accounts"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).VestingAccountAll(ctx, &launchtypes.QueryAllVestingAccountRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return []launchtypes.VestingAccount{}, err
	}
	return res.VestingAccount, err
}

// GenesisValidators returns the list of approved genesis validators for a launch from SPN
func (b *Builder) GenesisValidators(ctx context.Context, launchID uint64) ([]launchtypes.GenesisValidator, error) {
	b.ev.Send(events.New(events.StatusOngoing, "Fetching genesis validators"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).GenesisValidatorAll(ctx, &launchtypes.QueryAllGenesisValidatorRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return []launchtypes.GenesisValidator{}, err
	}
	return res.GenesisValidator, err
}