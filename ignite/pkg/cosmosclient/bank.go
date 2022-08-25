package cosmosclient

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

func (c Client) BankBalances(ctx context.Context, address string, pagination *query.PageRequest) (sdk.Coins, error) {
	defer c.lockBech32Prefix()()

	req := &banktypes.QueryAllBalancesRequest{
		Address:    address,
		Pagination: pagination,
	}

	resp, err := banktypes.NewQueryClient(c.context).AllBalances(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Balances, nil
}

func (c Client) BankSendTx(fromAccount cosmosaccount.Account, toAddress string, amount sdk.Coins) (TxService, error) {
	addr, err := fromAccount.Address(c.addressPrefix)
	if err != nil {
		return TxService{}, err
	}

	msg := &banktypes.MsgSend{
		FromAddress: addr,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	return c.CreateTx(fromAccount, msg)
}
