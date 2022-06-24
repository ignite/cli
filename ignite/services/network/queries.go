package network

import (
	"context"
	"sort"
	"sync"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

var (
	// ErrObjectNotFound is returned when the query returns a not found error.
	ErrObjectNotFound = errors.New("query object not found")
)

// ChainLaunch fetches the chain launch from Network by launch id.
func (n Network) ChainLaunch(ctx context.Context, id uint64) (networktypes.ChainLaunch, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching chain information"))

	res, err := n.launchQuery.
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

// ChainLaunchesWithReward fetches the chain launches with rewards from Network
func (n Network) ChainLaunchesWithReward(ctx context.Context) ([]networktypes.ChainLaunch, error) {
	g, ctx := errgroup.WithContext(ctx)

	n.ev.Send(events.New(events.StatusOngoing, "Fetching chains information"))
	res, err := n.launchQuery.
		ChainAll(ctx, &launchtypes.QueryAllChainRequest{})
	if err != nil {
		return nil, err
	}

	n.ev.Send(events.New(events.StatusOngoing, "Fetching reward information"))
	var chainLaunches []networktypes.ChainLaunch
	var mu sync.Mutex

	// Parse fetched chains and fetch rewards
	for _, chain := range res.Chain {
		chain := chain
		g.Go(func() error {
			chainLaunch := networktypes.ToChainLaunch(chain)
			reward, err := n.ChainReward(ctx, chain.LaunchID)
			if err != nil && err != ErrObjectNotFound {
				return err
			}
			chainLaunch.Reward = reward.RemainingCoins.String()
			mu.Lock()
			chainLaunches = append(chainLaunches, chainLaunch)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	// sort filenames by launch id
	sort.Slice(chainLaunches, func(i, j int) bool {
		return chainLaunches[i].ID > chainLaunches[j].ID
	})
	return chainLaunches, nil
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
	res, err := n.launchQuery.
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
	res, err := n.launchQuery.
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
	res, err := n.launchQuery.
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

// MainnetAccounts returns the list of campaign mainnet accounts for a launch from SPN
func (n Network) MainnetAccounts(ctx context.Context, campaignID uint64) (genAccs []networktypes.MainnetAccount, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching campaign mainnet accounts"))
	res, err := n.campaignQuery.
		MainnetAccountAll(ctx,
			&campaigntypes.QueryAllMainnetAccountRequest{
				CampaignID: campaignID,
			},
		)
	if err != nil {
		return genAccs, err
	}

	for _, acc := range res.MainnetAccount {
		genAccs = append(genAccs, networktypes.ToMainnetAccount(acc))
	}

	return genAccs, nil
}

// ChainReward fetches the chain reward from SPN by launch id
func (n Network) ChainReward(ctx context.Context, launchID uint64) (rewardtypes.RewardPool, error) {
	res, err := n.rewardQuery.
		RewardPool(ctx,
			&rewardtypes.QueryGetRewardPoolRequest{
				LaunchID: launchID,
			},
		)

	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return rewardtypes.RewardPool{}, ErrObjectNotFound
	} else if err != nil {
		return rewardtypes.RewardPool{}, err
	}
	return res.RewardPool, nil
}

// stakingParams fetches the staking module params
func (n Network) stakingParams(ctx context.Context) (stakingtypes.Params, error) {
	res, err := n.stakingQuery.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return stakingtypes.Params{}, err
	}
	return res.Params, nil
}

// RewardsInfo Fetches the consensus state with the validator set,
// the unbounding time, and the last block height from chain rewards.
func (n Network) RewardsInfo(
	ctx context.Context,
	launchID uint64,
	height int64,
) (
	rewardsInfo networktypes.Reward,
	lastRewardHeight int64,
	unboundingTime int64,
	err error,
) {
	rewardsInfo, err = RewardsInfo(ctx, n.cosmos, height)
	if err != nil {
		return rewardsInfo, 0, 0, err
	}

	stakingParams, err := n.stakingParams(ctx)
	if err != nil {
		return rewardsInfo, 0, 0, err
	}
	unboundingTime = int64(stakingParams.UnbondingTime.Seconds())

	chainReward, err := n.ChainReward(ctx, launchID)
	if err == ErrObjectNotFound {
		return rewardsInfo, 1, unboundingTime, nil
	} else if err != nil {
		return rewardsInfo, 0, 0, err
	}
	lastRewardHeight = chainReward.LastRewardHeight

	return
}

// ChainID fetches the network chain id
func (n Network) ChainID(ctx context.Context) (string, error) {
	status, err := n.cosmos.Status(ctx)
	if err != nil {
		return "", err
	}
	return status.NodeInfo.Network, nil
}
