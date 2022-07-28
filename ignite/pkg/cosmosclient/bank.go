package cosmosclient

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (c Client) BankBalances(ctx context.Context, address string, pagination *query.PageRequest) (sdk.Coins, error) {
	req := &banktypes.QueryAllBalancesRequest{
		Address:    address,
		Pagination: pagination,
	}

	resp, err := performQuery(c, func() (*banktypes.QueryAllBalancesResponse, error) {
		return banktypes.NewQueryClient(c.context).AllBalances(ctx, req)
	})
	if err != nil {
		return nil, err
	}

	return resp.Balances, nil
}

func (c Client) BankSendTx(fromAddress string, toAddress string, amount sdk.Coins, fromAccountName string) (TxService, error) {
	msg := &banktypes.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	return c.CreateTx(fromAccountName, msg)
}
