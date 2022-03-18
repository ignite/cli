package network

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// Coordinator returns the Coordinator by address from SPN
func (n Network) Coordinator(ctx context.Context, address string) (networktypes.Coordinator, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching coordinator details"))
	resCoordByAddr, err := profiletypes.NewQueryClient(n.cosmos.Context).
		CoordinatorByAddress(ctx,
			&profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: address,
			},
		)
	if err != nil {
		return networktypes.Coordinator{}, err
	}
	resCoord, err := profiletypes.NewQueryClient(n.cosmos.Context).
		Coordinator(ctx,
			&profiletypes.QueryGetCoordinatorRequest{
				CoordinatorID: resCoordByAddr.CoordinatorByAddress.CoordinatorID,
			},
		)
	if err != nil {
		return networktypes.Coordinator{}, err
	}
	return networktypes.ToCoordinator(resCoord.Coordinator), nil
}

// Validator returns the Validator by address from SPN
func (n Network) Validator(ctx context.Context, address string) (networktypes.Validator, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching validator details"))
	res, err := profiletypes.NewQueryClient(n.cosmos.Context).
		Validator(ctx,
			&profiletypes.QueryGetValidatorRequest{
				Address: address,
			},
		)
	if err != nil {
		return networktypes.Validator{}, err
	}
	return networktypes.ToValidator(res.Validator), nil
}

// Balances returns the all balances by address from SPN
func (n Network) Balances(ctx context.Context, address string) (sdk.Coins, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching address balances"))
	res, err := banktypes.NewQueryClient(n.cosmos.Context).AllBalances(ctx,
		&banktypes.QueryAllBalancesRequest{
			Address: address,
		},
	)
	if err != nil {
		return sdk.Coins{}, err
	}
	return res.Balances, nil
}
