package faucet_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
	"github.com/ignite/cli/ignite/pkg/xurl"
	envtest "github.com/ignite/cli/integration"
)

const (
	addr = "cosmos1zqr2gd7hwkyw55knad0l6ml6ngutd70878evqj"
)

var (
	defaultCoins = []string{"10token", "1stake"}
	maxCoins     = []string{"102token", "100000000stake"}
)

func TestRequestCoinsFromFaucet(t *testing.T) {
	var (
		env          = envtest.New(t)
		app          = env.Scaffold("github.com/test/faucet")
		servers      = app.RandomizeServerPorts()
		faucetURL    = app.EnableFaucet(defaultCoins, maxCoins)
		ctx, cancel  = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		faucetClient = cosmosfaucet.NewClient(faucetURL)
	)
	isErrTransferRequest := func(err error, expectedCode int) {
		require.ErrorAs(t, err, &cosmosfaucet.ErrTransferRequest{})
		errTransfer := err.(cosmosfaucet.ErrTransferRequest)
		require.EqualValues(t, expectedCode, errTransfer.StatusCode)
	}

	// serve the app
	go func() {
		app.Serve("should serve app", envtest.ExecCtx(ctx))
	}()

	// wait servers to be online
	defer cancel()
	err := env.IsAppServed(ctx, servers.API)
	require.NoError(t, err)

	err = env.IsFaucetServed(ctx, faucetClient)
	require.NoError(t, err)

	// error "account doesn't have any balances" occurs if a sleep is not included
	time.Sleep(time.Second * 1)

	nodeAddr, err := xurl.HTTP(servers.RPC)
	require.NoError(t, err)

	cosmosClient, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(nodeAddr))
	require.NoError(t, err)

	// the faucet sends the default faucet coins value when not specified
	_, err = faucetClient.Transfer(ctx, cosmosfaucet.NewTransferRequest(addr, nil))
	require.NoError(t, err)
	checkAccountBalance(t, ctx, cosmosClient, addr, defaultCoins)

	// the faucet can send a specified amount of coins
	_, err = faucetClient.Transfer(ctx, cosmosfaucet.NewTransferRequest(addr, []string{"20token", "2stake"}))
	require.NoError(t, err)
	checkAccountBalance(t, ctx, cosmosClient, addr, []string{"30token", "3stake"})

	// faucet request fails on malformed coins
	_, err = faucetClient.Transfer(ctx, cosmosfaucet.NewTransferRequest(addr, []string{"no-token"}))
	isErrTransferRequest(err, http.StatusBadRequest)

	// faucet request fails when requesting more than max coins
	_, err = faucetClient.Transfer(ctx, cosmosfaucet.NewTransferRequest(addr, []string{"500token"}))
	isErrTransferRequest(err, http.StatusInternalServerError)

	// faucet request fails when transfer should fail
	_, err = faucetClient.Transfer(ctx, cosmosfaucet.NewTransferRequest(addr, []string{"500nonexistent"}))
	isErrTransferRequest(err, http.StatusInternalServerError)

	// send several request in parallel and check max coins is not overflown
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			c := faucetClient
			_, err := c.Transfer(ctx, cosmosfaucet.NewTransferRequest(addr, nil))
			if err != nil && !errors.As(err, &cosmosfaucet.ErrTransferRequest{}) {
				return err
			}
			return nil
		})
	}
	require.NoError(t, g.Wait())
	checkAccountBalance(t, ctx, cosmosClient, addr, []string{"130token", "13stake"})
}

func checkAccountBalance(t *testing.T, ctx context.Context, c cosmosclient.Client, accAddr string, coins []string) {
	resp, err := banktypes.NewQueryClient(c.Context()).AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
		Address: accAddr,
	})
	require.NoError(t, err)

	require.Len(t, resp.Balances, len(coins))
	expectedCoins, err := sdk.ParseCoinsNormalized(strings.Join(coins, ","))
	require.NoError(t, err)
	require.True(t, resp.Balances.IsEqual(expectedCoins),
		fmt.Sprintf("%s should be equals to %s", resp.Balances.String(), expectedCoins.String()),
	)
}
