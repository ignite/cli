package network

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	statusErr, ok := status.FromError(err)
	if ok && statusErr.Code() == codes.NotFound {
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
	statusErr, ok := status.FromError(err)
	if ok && statusErr.Code() == codes.NotFound {
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
	statusErr, ok := status.FromError(err)
	if ok && statusErr.Code() == codes.NotFound {
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
	statusErr, ok := status.FromError(err)
	if ok && statusErr.Code() == codes.NotFound {
		return sdk.Coins{}, ErrObjectNotFound
	} else if err != nil {
		return sdk.Coins{}, err
	}
	return res.Balances, nil
}

// Profile returns the address profile info
func (n Network) Profile(ctx context.Context, campaignID uint64) (networktypes.Profile, error) {
	address := n.account.Address(networktypes.SPN)
	vouchers, err := n.Balances(ctx, address)
	if err != nil {
		return networktypes.Profile{}, err
	}

	var shares, vestingShares campaigntypes.Shares
	if campaignID > 0 {
		acc, err := n.MainnetAccount(ctx, campaignID, address)
		switch {
		case err == ErrObjectNotFound:
			accVest, err := n.MainnetVestingAccount(ctx, campaignID, address)
			if err != nil && err != ErrObjectNotFound {
				return networktypes.Profile{}, err
			}
			shares = accVest.TotalShares
			vestingShares = accVest.Vesting
		case err != nil:
			return networktypes.Profile{}, err
		default:
			shares = acc.Shares
		}
	}

	var p interface{}
	p, err = n.Validator(ctx, address)
	if err == ErrObjectNotFound {
		p, err = n.Coordinator(ctx, address)
		if err != nil {
			return networktypes.Profile{}, err
		}
	} else if err != nil {
		return networktypes.Profile{}, err
	}
	return networktypes.ToProfile(p, campaignID, vouchers, shares, vestingShares), err
}
