//go:build !relayer

package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/config/chain/base"
	v1 "github.com/ignite/cli/v29/ignite/config/chain/v1"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/xyaml"
	envtest "github.com/ignite/cli/v29/integration"
)

const newProtoPath = "myProto"

var (
	bobName = "bob"
	cfg     = v1.Config{
		Config: base.Config{
			Version: 1,
			Build: base.Build{
				Proto: base.Proto{
					Path: newProtoPath,
				},
			},
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "slide moment original seven milk crawl help text kick fluid boring awkward doll wonder sure fragile plate grid hard next casual expire okay body",
				},
				{
					Name:     bobName,
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "trap possible liquid elite embody host segment fantasy swim cable digital eager tiny broom burden diary earn hen grow engine pigeon fringe claim program",
				},
			},
			Faucet: base.Faucet{
				Name:  &bobName,
				Coins: []string{"500token", "100000000stake"},
				Host:  ":4501",
			},
			Genesis: xyaml.Map{"chain_id": "mars-1"},
		},
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
			},
		},
	}
)

func TestChangeProtoPath(t *testing.T) {
	var (
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/protopath", "--proto-dir", newProtoPath)
		appPath = app.SourcePath()
		cfgPath = filepath.Join(appPath, chain.ConfigFilenames[0])
	)

	// set the custom config path.
	file, err := os.Create(cfgPath)
	require.NoError(t, err)
	require.NoError(t, yaml.NewEncoder(file).Encode(cfg))
	require.NoError(t, file.Close())
	app.SetConfigPath(cfgPath)

	env.Must(env.Exec("create a list with a custom proto path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "listUser", "email"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a map with a custom proto path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "map", "--yes", "mapUser", "email"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a single with a custom proto path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "single", "--yes", "singleUser", "email"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a query with a custom proto path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "query", "--yes", "foo", "--proto-dir"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a new module with parameter",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"foo",
				"--params",
				"bla,baz:uint",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a new module parameter in the mars module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"params",
				"--yes",
				"foo",
				"bar:int",
				"--module",
				"foo",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	app.EnsureSteady()
}
