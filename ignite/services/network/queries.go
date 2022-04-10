package network

import (
	"context"
	"sort"
	"sync"

	"github.com/pkg/errors"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
	"golang.org/x/sync/errgroup"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/pkg/events"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

var (
	// ErrObjectNotFound is returned when the query returns a not found error.
	ErrObjectNotFound = errors.New("query object not found")
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

// MainnetAccounts returns the list of campaign mainnet accounts for a launch from SPN
func (n Network) MainnetAccounts(ctx context.Context, campaignID uint64) (genAccs []networktypes.MainnetAccount, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching campaign mainnet accounts"))
	res, err := campaigntypes.NewQueryClient(n.cosmos.Context).
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

// MainnetVestingAccounts returns the list of campaign mainnet vesting accounts for a launch from SPN
func (n Network) MainnetVestingAccounts(ctx context.Context, campaignID uint64) (genAccs []networktypes.MainnetVestingAccount, err error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching campaign mainnet vesting accounts"))
	res, err := campaigntypes.NewQueryClient(n.cosmos.Context).
		MainnetVestingAccountAll(ctx,
			&campaigntypes.QueryAllMainnetVestingAccountRequest{
				CampaignID: campaignID,
			},
		)
	if err != nil {
		return genAccs, err
	}

	for _, acc := range res.MainnetVestingAccount {
		genAccs = append(genAccs, networktypes.ToMainnetVestingAccount(acc))
	}

	return genAccs, nil
}

// ChainReward fetches the chain reward from SPN by launch id
func (n Network) ChainReward(ctx context.Context, launchID uint64) (rewardtypes.RewardPool, error) {
	res, err := rewardtypes.NewQueryClient(n.cosmos.Context).
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
