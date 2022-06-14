package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/pkg/events"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// CoordinatorIDByAddress returns the CoordinatorByAddress from SPN
func (n Network) CoordinatorIDByAddress(ctx context.Context, address string) (uint64, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching coordinator by address"))
	resCoordByAddr, err := n.profileQuery.
		CoordinatorByAddress(ctx,
			&profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: address,
			},
		)

	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return 0, ErrObjectNotFound
	} else if err != nil {
		return 0, err
	}
	return resCoordByAddr.CoordinatorByAddress.CoordinatorID, nil
}

// Coordinator returns the Coordinator by address from SPN
func (n Network) Coordinator(ctx context.Context, address string) (networktypes.Coordinator, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching coordinator details"))
	coordinatorID, err := n.CoordinatorIDByAddress(ctx, address)
	if err != nil {
		return networktypes.Coordinator{}, err
	}
	resCoord, err := n.profileQuery.
		Coordinator(ctx,
			&profiletypes.QueryGetCoordinatorRequest{
				CoordinatorID: coordinatorID,
			},
		)
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return networktypes.Coordinator{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.Coordinator{}, err
	}
	return networktypes.ToCoordinator(resCoord.Coordinator), nil
}

// Validator returns the Validator by address from SPN
func (n Network) Validator(ctx context.Context, address string) (networktypes.Validator, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching validator details"))
	res, err := n.profileQuery.
		Validator(ctx,
			&profiletypes.QueryGetValidatorRequest{
				Address: address,
			},
		)
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return networktypes.Validator{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.Validator{}, err
	}
	return networktypes.ToValidator(res.Validator), nil
}

// Balances returns the all balances by address from SPN
func (n Network) Balances(ctx context.Context, address string) (sdk.Coins, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching address balances"))
	res, err := banktypes.NewQueryClient(n.cosmos.Context()).AllBalances(ctx,
		&banktypes.QueryAllBalancesRequest{
			Address: address,
		},
	)
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return sdk.Coins{}, ErrObjectNotFound
	} else if err != nil {
		return sdk.Coins{}, err
	}
	return res.Balances, nil
}

// MainnetAccountBalance returns the spn account or vesting account balance by address from SPN
func (n Network) MainnetAccountBalance(
	ctx context.Context,
	launchID uint64,
	address string,
) (balance sdk.Coins, vesting sdk.Coins, err error) {
	acc, err := n.GenesisAccount(ctx, launchID, address)
	switch {
	case err == ErrObjectNotFound:
		accVest, err := n.VestingAccount(ctx, launchID, address)
		if err != nil && err != ErrObjectNotFound {
			return nil, nil, err
		}
		balance = accVest.TotalBalance
		vesting = accVest.Vesting
	case err != nil:
		return nil, nil, err
	default:
		balance = acc.Coins
	}
	return
}

// Profile returns the address profile info
func (n Network) Profile(ctx context.Context, campaign bool, campaignID uint64) (networktypes.Profile, error) {
	address := n.account.Address(networktypes.SPN)

	vouchers, err := n.Balances(ctx, address)
	if err != nil {
		return networktypes.Profile{}, err
	}

	var (
		shares             campaigntypes.Shares
		chainShares        []networktypes.ChainShare
		chainVestingShares []networktypes.ChainShare
	)

	if campaign {
		campaignChains, err := n.CampaignChains(ctx, campaignID)
		if err == ErrObjectNotFound {
			return networktypes.Profile{}, fmt.Errorf("invalid campaign id %d", campaignID)
		} else if err != nil {
			return networktypes.Profile{}, err
		}

		acc, err := n.MainnetAccount(ctx, campaignID, address)
		if err != nil && err != ErrObjectNotFound {
			return networktypes.Profile{}, err
		}
		shares = acc.Shares

		for _, chain := range campaignChains.Chains {
			balance, vesting, err := n.MainnetAccountBalance(ctx, chain, address)
			if err != nil && err != ErrObjectNotFound {
				return networktypes.Profile{}, err
			}
			chainShares = append(chainShares, networktypes.ChainShare{
				LaunchID: chain,
				Shares:   balance,
			})
			chainVestingShares = append(chainVestingShares, networktypes.ChainShare{
				LaunchID: chain,
				Shares:   vesting,
			})
		}
	}

	var p networktypes.ProfileAcc
	p, err = n.Validator(ctx, address)
	if err == ErrObjectNotFound {
		p, err = n.Coordinator(ctx, address)
		if err == ErrObjectNotFound {
			p = networktypes.Coordinator{Address: address}
		} else if err != nil {
			return networktypes.Profile{}, err
		}
	} else if err != nil {
		return networktypes.Profile{}, err
	}
	return p.ToProfile(campaignID, vouchers, shares, chainShares, chainVestingShares), nil
}
