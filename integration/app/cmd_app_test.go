//go:build !relayer

package app_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/services/chain"
	envtest "github.com/ignite/cli/v29/integration"
)

// TestGenerateAnApp tests scaffolding a new chain.
func TestGenerateAnApp(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	app.EnsureSteady()
}

// TestGenerateAnAppMinimal tests scaffolding a new minimal chain.
func TestGenerateAnAppMinimal(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("blog", "--minimal")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	app.EnsureSteady()
}

// TestGenerateAnAppWithName tests scaffolding a new chain using a local name instead of a GitHub URI.
func TestGenerateAnAppWithName(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("blog")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	app.EnsureSteady()
}

// TestGenerateAnAppWithInvalidName tests scaffolding a new chain using an invalid name.
func TestGenerateAnAppWithInvalidName(t *testing.T) {
	buf := new(bytes.Buffer)

	env := envtest.New(t)
	env.Must(env.Exec("should prevent creating an app with an invalid name",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "chain", "blog2"),
			step.Stdout(buf),
			step.Stderr(buf),
		)),
		envtest.ExecShouldError(),
	))

	require.Contains(t, buf.String(), "Invalid app name blog2: cannot contain numbers")
}

func TestGenerateAnAppWithNoDefaultModule(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog", "--no-module")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.True(t, os.IsNotExist(statErr), "the default module should not be scaffolded")

	app.EnsureSteady()
}

func TestGenerateAnAppWithNoDefaultModuleAndCreateAModule(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog", "--no-module")
	)

	defer app.EnsureSteady()

	app.Scaffold(
		"should scaffold a new module into a chain that never had modules before",
		false,
		"module", "first_module",
	)
}

func TestGenerateAppWithEmptyModule(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	app.Scaffold(
		"create a module",
		false,
		"module", "example", "--require-registration",
	)

	app.Scaffold(
		"should prevent creating an existing module",
		true,
		"module", "example", "--require-registration",
	)

	app.Scaffold(
		"should prevent creating a module with an invalid name",
		true,
		"module", "example1", "--require-registration",
	)

	app.Scaffold(
		"should prevent creating a module with a reserved name",
		true,
		"module", "tx", "--require-registration",
	)

	app.Scaffold(
		"should prevent creating a module with a forbidden prefix",
		true,
		"module", "ibcfoo", "--require-registration",
	)

	app.Scaffold(
		"should prevent creating a module prefixed with an existing module",
		true,
		"module", "examplefoo", "--require-registration",
	)

	app.Scaffold(
		"create a module with dependencies",
		false,
		"module",
		"with_dep",
		"--dep",
		"auth,bank,staking,slashing,example",
		"--require-registration",
	)

	app.Scaffold(
		"should prevent creating a module with invalid dependencies",
		true,
		"module",
		"with_wrong_dep",
		"--dep",
		"dup,dup",
		"--require-registration",
	)

	app.Scaffold(
		"should prevent creating a module with a non registered dependency",
		true,
		"module",
		"with_no_dep",
		"--dep",
		"inexistent",
		"--require-registration",
	)

	app.EnsureSteady()
}

func TestGenerateAnAppWithAddressPrefix(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog", "--address-prefix=gm", "--coin-type=60")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	c, err := chain.New(app.SourcePath())
	require.NoError(t, err, "failed to get new chain")

	bech32Prefix, err := c.Bech32Prefix()
	require.NoError(t, err)

	require.Equal(t, bech32Prefix, "gm")

	coinType, err := c.CoinType()
	require.NoError(t, err, "failed to get coin type")
	require.Equal(t, coinType, uint32(60))

	app.EnsureSteady()
}

func TestGenerateAnAppWithDefaultDenom(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog", "--default-denom=good")
	)

	_, statErr := os.Stat(filepath.Join(app.SourcePath(), "x", "blog"))
	require.False(t, os.IsNotExist(statErr), "the default module should be scaffolded")

	c, err := chain.New(app.SourcePath())
	require.NoError(t, err, "failed to get new chain")

	cfg, err := c.Config()
	require.NoError(t, err)

	require.Equal(t, cfg.DefaultDenom, "good")

	app.EnsureSteady()
}
