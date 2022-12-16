package account_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/randstr"
	envtest "github.com/ignite/cli/integration"
)

const testAccountMnemonic = "develop mansion drum glow husband trophy labor jelly fault run pause inside jazz foil page injury foam oppose fruit chunk segment morning series nation"

func TestAccount(t *testing.T) {
	var (
		env         = envtest.New(t)
		tmpDir      = t.TempDir()
		accountName = randstr.Runes(10)
	)

	env.Must(env.Exec("create account",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "account", "create", accountName, "--keyring-dir", tmpDir),
		)),
	))

	listOutputBuffer := &bytes.Buffer{}
	env.Must(env.Exec("list accounts",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "account", "list", "--keyring-dir", tmpDir),
		)),
		envtest.ExecStdout(listOutputBuffer),
	))
	require.True(t, strings.Contains(listOutputBuffer.String(), accountName))

	env.Must(env.Exec("delete account",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "account", "delete", accountName, "--keyring-dir", tmpDir),
		)),
	))

	listOutputAfterDeleteBuffer := &bytes.Buffer{}
	env.Must(env.Exec("list accounts after delete",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "account", "list", "--keyring-dir", tmpDir),
		)),
		envtest.ExecStdout(listOutputAfterDeleteBuffer),
	))
	require.Equal(t, "Name \tAddress Public Key \t\n\n", listOutputAfterDeleteBuffer.String())

	env.Must(env.Exec("import account with mnemonic",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp, "account", "import", "testaccount42",
				"--keyring-dir", tmpDir,
				"--secret", testAccountMnemonic,
			),
		)),
	))

	env.Must(env.Exec("import account with private key",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp, "account", "import", "testaccount43",
				"--non-interactive",
				"--keyring-dir", tmpDir,
				"--secret", "testdata/key",
				"--passphrase", "passpass",
			),
		)),
	))

	listOutputAfterImportBuffer := &bytes.Buffer{}
	env.Must(env.Exec("list accounts after import",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "account", "list", "--keyring-dir", tmpDir),
		)),
		envtest.ExecStdout(listOutputAfterImportBuffer),
	))
	require.Equal(t, `Name 		Address 					Public Key 										
testaccount42 	cosmos1ytnkpns7mfd6jjkvq9ztdvjdrt2xvmft2qxzqd 	PubKeySecp256k1{02FDF6D6F63B6B8E3CC71D03669BE0808F9990EE2A7FDBBF47E6BBEC4176E7763C} 	
testaccount43 	cosmos18p4xchk2aqp39nsjwr69ql44upzsfnh8r9lzql 	PubKeySecp256k1{0228B8FC609973D91BFF7A9933424F31E15A54B97F8FDF7CE5A83B4DC20988068F} 	

`, listOutputAfterImportBuffer.String())

	showOutputBuffer := &bytes.Buffer{}
	env.Must(env.Exec("show account",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "account", "show", "testaccount42", "--keyring-dir", tmpDir),
		)),
		envtest.ExecStdout(showOutputBuffer),
	))
	require.Equal(t, `Name 		Address 					Public Key 										
testaccount42 	cosmos1ytnkpns7mfd6jjkvq9ztdvjdrt2xvmft2qxzqd 	PubKeySecp256k1{02FDF6D6F63B6B8E3CC71D03669BE0808F9990EE2A7FDBBF47E6BBEC4176E7763C} 	

`, showOutputBuffer.String())

	showOutputWithDifferentPrefixBuffer := &bytes.Buffer{}
	env.Must(env.Exec("show account with address prefix",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "account", "show", "testaccount42", "--keyring-dir", tmpDir, "--address-prefix", "test"),
		)),
		envtest.ExecStdout(showOutputWithDifferentPrefixBuffer),
	))
	require.Equal(t, `Name 		Address 					Public Key 										
testaccount42 	test1ytnkpns7mfd6jjkvq9ztdvjdrt2xvmftxemuve 	PubKeySecp256k1{02FDF6D6F63B6B8E3CC71D03669BE0808F9990EE2A7FDBBF47E6BBEC4176E7763C} 	

`, showOutputWithDifferentPrefixBuffer.String())
}
