package envtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/xurl"
)

const defaultRequestTimeout = 90 * time.Second

type TxResponse struct {
	Code      int    `json:"code"`
	Codespace string `json:"codespace"`
	RawLog    string `json:"raw_log"`
	TxHash    string `json:"txhash"`
	Height    string `json:"height"`
	Data      string `json:"data"`
	Info      string `json:"info"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
	Timestamp string `json:"timestamp"`
}

func (a App) CLITx(chainRPC, module, method string, args ...string) TxResponse {
	nodeAddr, err := xurl.TCP(chainRPC)
	require.NoErrorf(a.env.T(), err, "cant read nodeAddr from host.RPC %v", chainRPC)

	args = append(args,
		"--node", nodeAddr,
		"--home", a.homePath,
		"--from", "alice",
		"--output", "json",
		"--log_format", "json",
		"--keyring-backend", "test",
		"--yes",
	)
	var (
		output     = &bytes.Buffer{}
		outErr     = &bytes.Buffer{}
		txResponse = TxResponse{}
	)
	stepsTx := step.NewSteps(
		step.New(
			step.Stdout(output),
			step.Stderr(outErr),
			step.PreExec(func() error {
				output.Reset()
				outErr.Reset()
				return nil
			}),
			step.Exec(
				a.Binary(),
				append([]string{"tx", module, method}, args...)...,
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if outErr.Len() > 0 {
					return errors.Errorf("error executing request: %s", outErr.String())
				}
				output := output.Bytes()
				if err := json.Unmarshal(output, &txResponse); err != nil {
					return errors.Errorf("unmarshalling tx response error: %w, response: %s", err, string(output))
				}
				return nil
			}),
		))

	ctx, cancel := context.WithTimeout(a.env.Ctx(), defaultRequestTimeout)
	defer cancel()

	if !a.env.Exec("sending chain request "+args[0], stepsTx, ExecRetry(), ExecCtx(ctx)) {
		cancel()
		a.env.t.FailNow()
	}

	return txResponse
}

func (a App) CLIQueryTx(chainRPC, txHash string) (txResponse TxResponse) {
	output := a.query(chainRPC, "tx", txHash)
	err := json.Unmarshal(output, &txResponse)
	require.NoError(a.env.T(), err, "unmarshalling tx response: %s", string(output))
	return txResponse
}

func (a App) CLIQuery(chainRPC, module, method string, args ...string) []byte {
	return a.query(chainRPC, module, method, args...)
}

func (a App) query(chainRPC, module, method string, args ...string) []byte {
	nodeAddr, err := xurl.TCP(chainRPC)
	require.NoErrorf(a.env.T(), err, "cant read nodeAddr from host.RPC %v", chainRPC)

	var (
		output = &bytes.Buffer{}
		outErr = &bytes.Buffer{}
	)

	cmd := append([]string{"query", module, method}, args...)
	cmd = append(cmd,
		"--node", nodeAddr,
		"--home", a.homePath,
		"--output", "json",
		"--log_format", "json",
	)
	steps := step.NewSteps(
		step.New(
			step.Stdout(output),
			step.Stderr(outErr),
			step.PreExec(func() error {
				output.Reset()
				outErr.Reset()
				return nil
			}),
			step.Exec(a.Binary(), cmd...),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if outErr.Len() > 0 {
					return errors.Errorf("error executing request: %s", outErr.String())
				}
				return nil
			}),
		))

	if !a.env.Exec(fmt.Sprintf("fetching query data %s => %s", module, method), steps, ExecRetry()) {
		a.env.t.FailNow()
	}

	return output.Bytes()
}

func (a App) APIQuery(ctx context.Context, chainAPI, namespace, module, method string, args ...string) []byte {
	ctx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	chainAPI, err := xurl.HTTP(chainAPI)
	require.NoErrorf(a.env.T(), err, "failed to convert chain API %s to HTTP", chainAPI)

	modulePath := gomodulepath.ExtractAppPath(namespace)
	apiURL, err := url.JoinPath(chainAPI, modulePath, module, "v1", method, strings.Join(args, "/"))
	require.NoErrorf(a.env.T(), err, "failed to create API URL")

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	require.NoErrorf(a.env.T(), err, "failed to create HTTP request")

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoErrorf(a.env.T(), err, "failed to execute HTTP request")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		require.Failf(a.env.T(), "unexpected status code", "expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	require.NoErrorf(a.env.T(), err, "failed to read response body")

	return body
}
