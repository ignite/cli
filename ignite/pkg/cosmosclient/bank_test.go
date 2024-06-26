package cosmosclient_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientBankBalances(t *testing.T) {
	var (
		ctx              = context.Background()
		address          = "address"
		pagination       = &query.PageRequest{Offset: 1}
		expectedBalances = sdk.NewCoins(
			sdk.NewCoin("token", math.NewInt(1000)),
			sdk.NewCoin("stake", math.NewInt(2000)),
		)
	)
	c := newClient(t, func(s suite) {
		req := &banktypes.QueryAllBalancesRequest{
			Address:    address,
			Pagination: pagination,
		}

		s.bankQueryClient.EXPECT().AllBalances(ctx, req).
			Return(&banktypes.QueryAllBalancesResponse{
				Balances: expectedBalances,
			}, nil)
	})

	balances, err := c.BankBalances(ctx, address, pagination)

	require.NoError(t, err)
	assert.Equal(t, expectedBalances, balances)
}
