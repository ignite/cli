//go:build !relayer

package tx_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/randstr"
	"github.com/ignite/cli/v29/ignite/pkg/xurl"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestSignTxWithDashedAppName(t *testing.T) {
	var (
		env         = envtest.New(t)
		appname     = "da-shed-a-p-p"
		app         = env.ScaffoldApp(appname)
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)

	nodeAddr, err := xurl.TCP(servers.RPC)
	require.NoErrorf(t, err, "cant read nodeAddr from host.RPC %v", servers.RPC)

	app.Scaffold(
		"scaffold a simple list",
		false,
		"list",
		"item",
		"str",
	)

	var (
		output            = &bytes.Buffer{}
		isTxBodyRetrieved bool
		txResponse        struct {
			Code   int
			RawLog string `json:"raw_log"`
		}
	)
	// sign tx to add an item to the list.
	steps := step.NewSteps(
		step.New(
			step.Stdout(output),
			step.Exec(
				app.Binary(),
				"tx",
				"dashedapp",
				"create-item",
				"helloworld",
				"--chain-id", "dashedapp",
				"--from", "alice",
				"--node", nodeAddr,
				"--output", "json",
				"--log_format", "json",
				"--yes",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				err := json.Unmarshal(output.Bytes(), &txResponse)
				if err != nil {
					return errors.Errorf("unmarshling tx response: %w", err)
				}
				return nil
			}),
		),
	)

	go func() {
		defer cancel()
		app.WaitChainUp(ctx, servers.API)
		isTxBodyRetrieved = env.Exec("sign a tx", steps, envtest.ExecRetry())
	}()

	app.MustServe(ctx)

	if !isTxBodyRetrieved {
		t.FailNow()
	}
	require.Equal(t, 0, txResponse.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)
}

func TestGetTxViaGRPCGateway(t *testing.T) {
	var (
		env         = envtest.New(t)
		appname     = randstr.Runes(10)
		app         = env.ScaffoldApp(fmt.Sprintf("github.com/test/%s", appname))
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)

	var (
		output            = &bytes.Buffer{}
		isTxBodyRetrieved bool
		txBody            = struct {
			Tx struct {
				Body struct {
					Messages []struct {
						Amount []struct {
							Denom  string `json:"denom"`
							Amount string `json:"amount"`
						} `json:"amount"`
					} `json:"messages"`
				} `json:"body"`
			} `json:"tx"`
		}{}
	)

	// 1- list accounts
	// 2- send tokens from one to other.
	// 3- verify tx by using gRPC Gateway API.
	steps := step.NewSteps(
		step.New(
			step.Exec(
				app.Binary(),
				"keys",
				"list",
				"--keyring-backend", "test",
				"--output", "json",
				"--log_format", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}

				// collect addresses of alice and bob.
				var (
					accounts []struct {
						Name    string `json:"name"`
						Address string `json:"address"`
					}
					addresses []string
				)
				if err := json.NewDecoder(output).Decode(&accounts); err != nil {
					return err
				}
				for _, account := range accounts {
					if account.Name == "alice" || account.Name == "bob" {
						addresses = append(addresses, account.Address)
					}
				}
				if len(addresses) != 2 {
					return errors.New("expected alice and bob accounts to be created")
				}

				nodeAddr, err := xurl.TCP(servers.RPC)
				require.NoErrorf(t, err, "cant read nodeAddr from host.RPC %v", servers.RPC)

				// send some tokens from alice to bob and confirm the corresponding tx via gRPC gateway
				// endpoint by asserting denom and amount.
				return cmdrunner.New().Run(ctx, step.New(
					step.Exec(
						app.Binary(),
						"tx",
						"bank",
						"send",
						addresses[0],
						addresses[1],
						"10token",
						"--keyring-backend", "test",
						"--chain-id", appname,
						"--node", nodeAddr,
						"--output", "json",
						"--log_format", "json",
						"--yes",
					),
					step.PreExec(func() error {
						output.Reset()
						return nil
					}),
					step.PostExec(func(execErr error) error {
						if execErr != nil {
							return execErr
						}

						tx := struct {
							Hash string `json:"txHash"`
						}{}
						if err := json.NewDecoder(output).Decode(&tx); err != nil {
							return err
						}

						apiAddr, err := xurl.HTTP(servers.API)
						if err != nil {
							return err
						}

						addr := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", apiAddr, tx.Hash)
						req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
						if err != nil {
							return errors.Wrap(err, "call to get tx via gRPC gateway")
						}

						time.Sleep(5 * time.Second)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return err
						}
						defer resp.Body.Close()

						// Send error if the request failed
						if resp.StatusCode != http.StatusOK {
							return errors.New(resp.Status)
						}

						if err := json.NewDecoder(resp.Body).Decode(&txBody); err != nil {
							return err
						}
						return nil
					}),
					step.Stdout(output),
				))
			}),
			step.Stdout(output),
		))

	go func() {
		defer cancel()
		app.WaitChainUp(ctx, servers.API)
		isTxBodyRetrieved = env.Exec("retrieve account addresses", steps, envtest.ExecRetry())
	}()

	app.MustServe(ctx)

	if !isTxBodyRetrieved {
		t.FailNow()
	}

	require.Len(t, txBody.Tx.Body.Messages, 1)
	require.Len(t, txBody.Tx.Body.Messages[0].Amount, 1)
	require.Equal(t, "token", txBody.Tx.Body.Messages[0].Amount[0].Denom)
	require.Equal(t, "10", txBody.Tx.Body.Messages[0].Amount[0].Amount)
}
