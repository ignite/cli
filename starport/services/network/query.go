package network

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/services/network/networktypes"

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

// GenesisInformation returns all the information to construct the genesis from a chain ID
func (n Network) GenesisInformation(ctx context.Context, launchID uint64) (gi networktypes.GenesisInformation, err error) {
	genAccs, err := n.GenesisAccounts(ctx, launchID)
	if err != nil {
		return gi, errors.Wrap(err, "error querying genesis accounts")
	}

	vestingAccs, err := n.VestingAccounts(ctx, launchID)
	if err != nil {
		return gi, errors.Wrap(err, "error querying vesting accounts")
	}

	genVals, err := n.GenesisValidators(ctx, launchID)
	if err != nil {
		return gi, errors.Wrap(err, "error querying genesis validators")
	}

	return networktypes.NewGenesisInformation(genAccs, vestingAccs, genVals), nil
}

// GenesisAccounts returns the list of approved genesis account for a launch from SPN
func (n Network) GenesisAccounts(ctx context.Context, launchID uint64) (genAccs []networktypes.GenesisAccount, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching genesis accounts"))
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).GenesisAccountAll(ctx, &launchtypes.QueryAllGenesisAccountRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return genAccs, err
	}

	for _, acc := range res.GenesisAccount {
		genAccs = append(genAccs, networktypes.ParseGenesisAccount(acc))
	}

	return genAccs, nil
}

// VestingAccounts returns the list of approved genesis vesting account for a launch from SPN
func (n Network) VestingAccounts(ctx context.Context, launchID uint64) (vestingAccs []networktypes.VestingAccount, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching genesis vesting accounts"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).VestingAccountAll(ctx, &launchtypes.QueryAllVestingAccountRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return vestingAccs, err
	}

	for i, acc := range res.VestingAccount {
		parsedAcc, err := networktypes.ParseVestingAccount(acc)
		if err != nil {
			return vestingAccs, errors.Wrapf(err, "error parsing vesting account %d", i)
		}

		vestingAccs = append(vestingAccs, parsedAcc)
	}

	return vestingAccs, nil
}

// GenesisValidators returns the list of approved genesis validators for a launch from SPN
func (n Network) GenesisValidators(ctx context.Context, launchID uint64) (genVals []networktypes.GenesisValidator, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching genesis validators"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).GenesisValidatorAll(ctx, &launchtypes.QueryAllGenesisValidatorRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return genVals, err
	}

	for _, acc := range res.GenesisValidator {
		genVals = append(genVals, networktypes.ParseGenesisValidator(acc))
	}

	return genVals, nil
}
