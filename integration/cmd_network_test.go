package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/chaintest"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/randstr"
	"github.com/tendermint/starport/starport/pkg/xurl"
)

func TestNetworkChainCreate(t *testing.T) {
	var (
		env       = chaintest.New(t)
		chainName = randstr.Runes(10)
		spnPath   = env.Pull(spn())
		servers   = env.RandomizeServerPorts(spnPath, "")
	)

	ctx, cancel := context.WithTimeout(env.Ctx(), chaintest.ServeTimeout)
	defer cancel()

	go func() { env.Serve("run spn", spnPath, chaintest.ServeWithExecOption(chaintest.ExecCtx(ctx))) }()

	require.NoError(t, env.IsAppServed(ctx, servers.Host), "app cannot get online in time")

	subcommand := []string{
		"network",
		"--spn-api-address",
		xurl.HTTP(servers.Host.API),
		"--spn-faucet-address",
		xurl.HTTP(servers.Faucet.Host),
		"--spn-node-address",
		xurl.HTTP(servers.Host.RPC),
	}

	env.Must(env.Exec("create a chain",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				append(subcommand, []string{
					"chain",
					"create",
					chainName,
					"https://github.com/cosmos/gaia",
					"--tag",
					"v4.2.1",
				}...)...,
			),
			step.Write([]byte("y\n")),
		)),
	))

	res := &bytes.Buffer{}

	env.Must(env.Exec("get created chain",
		step.NewSteps(step.New(
			step.Exec(
				"starport",
				append(subcommand, []string{
					"chain",
					"show",
					chainName,
				}...)...,
			),
			step.Stdout(res),
		)),
	))

	require.Contains(t, res.String(), fmt.Sprintf("chainid: %s", chainName))
}
