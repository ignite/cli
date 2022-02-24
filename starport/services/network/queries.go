package network

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// ChainLaunch fetches the chain launch from Starport Network by launch id.
func (n Network) ChainLaunch(ctx context.Context, id uint64) (networktypes.ChainLaunch, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching chain information"))

	res, err := launchtypes.NewQueryClient(n.cosmos.Context).
		Chain(ctx,
			&launchtypes.QueryGetChainRequest{
				LaunchID: id,
			},
		)
	if err != nil {
		return networktypes.ChainLaunch{}, err
	}

	return networktypes.ToChainLaunch(res.Chain), nil
}

// ChainLaunchesWithReward fetches the chain launches with rewards from Starport Network
func (n Network) ChainLaunchesWithReward(ctx context.Context) ([]networktypes.ChainLaunch, error) {
	g, ctx := errgroup.WithContext(ctx)

	n.ev.Send(events.New(events.StatusOngoing, "Fetching chains information"))
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).
		ChainAll(ctx, &launchtypes.QueryAllChainRequest{})
	if err != nil {
		return nil, err
	}

	n.ev.Send(events.New(events.StatusOngoing, "Fetching rewards"))
	var chainLaunches []networktypes.ChainLaunch
	var mu sync.Mutex

	// Parse fetched chains and fetch rewards
	for _, chain := range res.Chain {
		chain := chain
		g.Go(func() error {
			chainLaunch := networktypes.ToChainLaunch(chain)
			reward, err := n.ChainReward(ctx, chain.LaunchID)
			if err != nil {
				return err
			}
			chainLaunch.Reward = reward.Coins.String()
			mu.Lock()
			chainLaunches = append(chainLaunches, chainLaunch)
			mu.Unlock()
			return nil
		})
	}
	return chainLaunches, g.Wait()
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

// GenesisAccounts returns the list of approved genesis accounts for a launch from SPN
func (n Network) GenesisAccounts(ctx context.Context, launchID uint64) (genAccs []networktypes.GenesisAccount, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching genesis accounts"))
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).
		GenesisAccountAll(ctx,
			&launchtypes.QueryAllGenesisAccountRequest{
				LaunchID: launchID,
			},
		)
	if err != nil {
		return genAccs, err
	}

	for _, acc := range res.GenesisAccount {
		genAccs = append(genAccs, networktypes.ToGenesisAccount(acc))
	}

	return genAccs, nil
}

// VestingAccounts returns the list of approved genesis vesting accounts for a launch from SPN
func (n Network) VestingAccounts(ctx context.Context, launchID uint64) (vestingAccs []networktypes.VestingAccount, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching genesis vesting accounts"))
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).
		VestingAccountAll(ctx,
			&launchtypes.QueryAllVestingAccountRequest{
				LaunchID: launchID,
			},
		)
	if err != nil {
		return vestingAccs, err
	}

	for i, acc := range res.VestingAccount {
		parsedAcc, err := networktypes.ToVestingAccount(acc)
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
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).
		GenesisValidatorAll(ctx,
			&launchtypes.QueryAllGenesisValidatorRequest{
				LaunchID: launchID,
			},
		)
	if err != nil {
		return genVals, err
	}

	for _, acc := range res.GenesisValidator {
		genVals = append(genVals, networktypes.ToGenesisValidator(acc))
	}

	return genVals, nil
}

// ChainReward fetches the chain reward from SPN by launch id
func (n Network) ChainReward(ctx context.Context, launchID uint64) (rewardtypes.RewardPool, error) {
	res, err := rewardtypes.NewQueryClient(n.cosmos.Context).
		RewardPool(ctx,
			&rewardtypes.QueryGetRewardPoolRequest{
				LaunchID: launchID,
			},
		)
	if status.Code(err) == codes.InvalidArgument {
		return rewardtypes.RewardPool{}, nil
	} else if err != nil {
		return rewardtypes.RewardPool{}, err
	}
	return res.RewardPool, nil
}
