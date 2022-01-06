package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	envtest "github.com/tendermint/starport/integration"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	addr = "cosmos1zqr2gd7hwkyw55knad0l6ml6ngutd70878evqj"
)

var (
	defaultCoins = []string{"10token", "1stake"}
	maxCoins = []string{"102token", "100000000stake"}
)

type QueryAllBalancesResponse struct {
	Balances sdk.Coins `json:"balances"`
}

func TestRequestCoinsFromFaucet(t *testing.T) {
	var (
		env               = envtest.New(t)
		apath             = env.Scaffold("faucet")
		servers           = env.RandomizeServerPorts(apath, "")
		faucetURL            = env.ConfigureFaucet(apath, "", defaultCoins, maxCoins)
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	)
	// serve the app
	go func() {
		env.Must(env.Serve("should serve app", apath, "", "", envtest.ExecCtx(ctx)))
	}()

	// wait servers to be online
	defer cancel()
	err := env.IsAppServed(ctx, servers)
	require.NoError(t, err)

	err = env.IsFaucetServed(ctx, faucetURL)
	require.NoError(t, err)

	// error "account doesn't have any balances" occurs if a sleep is not included
	time.Sleep(time.Second*1)

	// the faucet sends the default faucet coins value when not specified
	resp, err := faucetRequest(faucetURL, addr, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, resp.Body.Close())
	checkAccountBalance(t, servers, addr, defaultCoins)

	// the faucet can send a specified amount of coins
	resp, err = faucetRequest(faucetURL, addr, []string{"20token", "2stake"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, resp.Body.Close())
	checkAccountBalance(t, servers, addr, []string{"30token", "3stake"})

	// faucet request fails on malformed coins
	resp, err = faucetRequest(faucetURL, addr, []string{"no-token"})
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.NoError(t, resp.Body.Close())

	// faucet request fails when requesting more than max coins
	resp, err = faucetRequest(faucetURL, addr, []string{"500token"})
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.NoError(t, resp.Body.Close())

	// send several request in parallel and check max coins is not overflown
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i<10; i++ {
		i := i
		g.Go(func() error {
			resp, err := faucetRequest(faucetURL, addr, nil)
			if err != nil {
				return err
			}

			res, err := io.ReadAll(resp.Body)
			var r cosmosfaucet.TransferResponse
			json.Unmarshal(res, &r)
			fmt.Printf("%d: %v\n", i, r.Error)

			return resp.Body.Close()
		})
	}
	require.NoError(t, g.Wait())
	checkAccountBalance(t, servers, addr, []string{"100token", "10stake"})
}

func faucetRequest(faucetURL string, accAddr string, coins []string) (*http.Response, error) {
	req := cosmosfaucet.TransferRequest{
		AccountAddress: accAddr,
		Coins: coins,
	}
	mReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(faucetURL,  "application/json", bytes.NewBuffer(mReq))
	return resp, err
}

func checkAccountBalance(t *testing.T, servers chainconfig.Host, accAddr string, coins []string) {
	// delay for the balance to be updated after faucet request
	time.Sleep(time.Second*1)

	balanceResp, err := http.Get(xurl.HTTP(servers.API) + fmt.Sprintf("/cosmos/bank/v1beta1/balances/%s", accAddr))
	require.NoError(t, err)
	defer balanceResp.Body.Close()

	var balanceRes QueryAllBalancesResponse
	res, err := io.ReadAll(balanceResp.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(res, &balanceRes))
	require.Len(t, balanceRes.Balances, len(coins))
	expectedCoins, err := sdk.ParseCoinsNormalized(strings.Join(coins, ","))
	require.NoError(t, err)
	require.True(t, balanceRes.Balances.IsEqual(expectedCoins), fmt.Sprintf("%s should be equals to %s", balanceRes.Balances.String(), expectedCoins.String()))
}
