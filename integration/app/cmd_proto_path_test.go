//go:build !relayer

package app_test

import (
	"os"
	"path/filepath"
	"strings"
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

// TestGenerateAppCheckBufPulsarPath tests scaffolding a new chain and checks if the buf.gen.pulsar.yaml file is correct.
func TestGenerateAppCheckBufPulsarPath(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blog")
	)

	bufGenPulsarPath := filepath.Join(app.SourcePath(), "proto", "buf.gen.pulsar.yaml")
	_, statErr := os.Stat(bufGenPulsarPath)
	require.False(t, os.IsNotExist(statErr), "buf.gen.pulsar.yaml should be scaffolded")

	result, err := os.ReadFile(bufGenPulsarPath)
	require.NoError(t, err)

	require.True(t, strings.Contains(string(result), "default: github.com/test/blog/api"), "buf.gen.pulsar.yaml should contain the correct api override")

	app.EnsureSteady()
}

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

	env.Must(env.Exec("create a list with a custom proto path from config",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteExtension, "s", "list", "--yes", "listUser", "email"),
			step.Workdir(app.SourcePath()),
		)),
	))

	app.EnsureSteady()
}
