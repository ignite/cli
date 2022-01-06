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
	defaultFaucetCoins = []string{"5token", "100000stake"}
	sampleCoins = []string{"5token", "5stake"}
)

type QueryAllBalancesResponse struct {
	Balances sdk.Coins `json:"balances"`
}

func TestRequestCoinsFromFaucet(t *testing.T) {
	var (
		env               = envtest.New(t)
		apath             = env.Scaffold("faucet")
		servers           = env.RandomizeServerPorts(apath, "")
		faucetURL            = env.ConfigureFaucet(apath, "")
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
	resp := faucetRequest(t, faucetURL, addr, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	time.Sleep(time.Second*1)
	checkAccountBalance(t, servers, addr, defaultFaucetCoins)

	// the faucet can send a specified amount of coins
	resp = faucetRequest(t, faucetURL, addr, sampleCoins)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	time.Sleep(time.Second*1)
	checkAccountBalance(t, servers, addr, []string{"10token", "100005stake"})

	// faucet request fails on malformed coins
	resp = faucetRequest(t, faucetURL, addr, []string{"no-token"})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()


}

func faucetRequest(t *testing.T, faucetURL string, accAddr string, coins []string) *http.Response {
	req := cosmosfaucet.TransferRequest{
		AccountAddress: accAddr,
		Coins: coins,
	}
	mReq, err := json.Marshal(req)
	require.NoError(t, err)
	resp, err := http.Post(faucetURL, "application/json", bytes.NewBuffer(mReq))
	require.NoError(t, err)
	return resp
}

func checkAccountBalance(t *testing.T, servers chainconfig.Host, accAddr string, coins []string) {
	balanceResp, err := http.Get(xurl.HTTP(servers.API) + fmt.Sprintf("/cosmos/bank/v1beta1/balances/%s", accAddr))
	require.NoError(t, err)
	defer balanceResp.Body.Close()
	var balanceRes QueryAllBalancesResponse
	res, err := io.ReadAll(balanceResp.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(res, &balanceRes))
	expectedCoins, err := sdk.ParseCoinsNormalized(strings.Join(coins, ","))
	require.NoError(t, err)
	require.True(t, balanceRes.Balances.IsEqual(expectedCoins))
}
