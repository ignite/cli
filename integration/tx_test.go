// +build !relayer

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/randstr"
	"github.com/tendermint/starport/starport/pkg/xurl"
)

func TestGetTxViaGRPCGateway(t *testing.T) {
	t.Parallel()

	var (
		env         = newEnv(t)
		appname     = randstr.Runes(10)
		path        = env.Scaffold(appname, Stargate)
		servers     = env.RandomizeServerPorts(path)
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
	steps := step.NewSteps(step.New(
		step.Exec(
			appname+"d",
			"keys",
			"list",
			"--keyring-backend", "test",
			"--output", "json",
		),
		step.PreExec(func() error {
			output.Reset()
			return env.IsAppServed(ctx, servers)
		}),
		step.PostExec(func(execErr error) error {
			if execErr != nil {
				return execErr
			}

			addresses := []string{}

			// collect addresses of user1 and user2.
			accounts := []struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			}{}
			if err := json.NewDecoder(output).Decode(&accounts); err != nil {
				return err
			}
			for _, account := range accounts {
				if account.Name == "user1" || account.Name == "user2" {
					addresses = append(addresses, account.Address)
				}
			}
			if len(addresses) != 2 {
				return errors.New("expected user1 and user2 accounts to be created")
			}

			// send some tokens from user1 to user2 and confirm the corresponding tx via gRPC gateway
			// endpoint by asserting denom and amount.
			return cmdrunner.New().Run(ctx, step.New(
				step.Exec(
					appname+"d",
					"tx",
					"bank",
					"send",
					addresses[0],
					addresses[1],
					"10token",
					"--keyring-backend", "test",
					"--chain-id", appname,
					"--node", xurl.TCP(servers.RPCAddr),
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

					addr := fmt.Sprintf("%s/cosmos/tx/v1beta1/tx/%s", xurl.HTTP(servers.APIAddr), tx.Hash)
					req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
					if err != nil {
						return err
					}
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return err
					}
					defer resp.Body.Close()

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

		isTxBodyRetrieved = env.Exec("retrieve account addresses", steps, ExecRetry())
	}()

	env.Must(env.Serve("should serve", path, ExecCtx(ctx)))

	if !isTxBodyRetrieved {
		t.FailNow()
	}

	require.Len(t, txBody.Tx.Body.Messages, 1)
	require.Len(t, txBody.Tx.Body.Messages[0].Amount, 1)
	require.Equal(t, "token", txBody.Tx.Body.Messages[0].Amount[0].Denom)
	require.Equal(t, "10", txBody.Tx.Body.Messages[0].Amount[0].Amount)
}
