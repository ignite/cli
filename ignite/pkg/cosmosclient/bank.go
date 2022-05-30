package cosmosclient

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (c Client) BankBalances(address string, pagination *query.PageRequest) (sdk.Coins, error) {
	req := &banktypes.QueryAllBalancesRequest{
		Address:    address,
		Pagination: pagination,
	}

	resp, err := performQuery[*banktypes.QueryAllBalancesResponse](c, func() (*banktypes.QueryAllBalancesResponse, error) {
		return banktypes.NewQueryClient(c.context).AllBalances(context.Background(), req)
	})
	if err != nil {
		return nil, err
	}

	return resp.Balances, nil
}

func (c Client) BankSend(fromAddress string, toAddress string, amount sdk.Coins, fromAccountName string) error {
	msg := &banktypes.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	_, err := c.BroadcastTx(fromAccountName, msg)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) BankSendGenerateOnly(fromAddress string, toAddress string, amount sdk.Coins, fromAccountName string) (string, error) {
	msg := &banktypes.MsgSend{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	return c.GenerateTx(fromAccountName, msg)
}
