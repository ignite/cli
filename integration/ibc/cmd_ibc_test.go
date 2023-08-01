//go:build !relayer

package ibc_test

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCreateModuleWithIBC(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blogibc")
	)

	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "foo", "--ibc", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an IBC module with custom path",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"appPath",
				"--ibc",
				"--require-registration",
				"--path",
				"./blogibc",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a type in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "list", "--yes", "user", "email", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an IBC module with an ordered channel",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"orderedfoo",
				"--ibc",
				"--ordering",
				"ordered",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an IBC module with an unordered channel",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"unorderedfoo",
				"--ibc",
				"--ordering",
				"unordered",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a non IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "non_ibc", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an IBC module with dependencies",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"example_with_dep",
				"--ibc",
				"--dep",
				"account,bank,staking,slashing",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	app.EnsureSteady()
}

// Deprecated: Oracle functionality is no longer tested.
func TestCreateIBCOracle(t *testing.T) {
	t.Skip() // TODO remove in future
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/ibcoracle")
	)

	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "foo", "--ibc", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create an IBC module with params",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"module",
				"--yes",
				"paramsFoo",
				"--ibc",
				"--params",
				"defaultName,isLaunched:bool,minLaunch:uint,maxLaunch:int",
				"--require-registration",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create the first BandChain oracle integration",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "band", "--yes", "oracleone", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create the second BandChain oracle integration",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "band", "--yes", "oracletwo", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a BandChain oracle with no module specified",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "band", "--yes", "invalidOracle"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a BandChain oracle in a non existent module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "band", "--yes", "invalidOracle", "--module", "nomodule"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a non-IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "bar", "--params", "name,minLaunch:uint,maxLaunch:int", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a BandChain oracle in a non IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "band", "--yes", "invalidOracle", "--module", "bar"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	app.EnsureSteady()
}

func TestCreateIBCPacket(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/blogibc2")
	)

	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "foo", "--ibc", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a packet",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"packet",
				"--yes",
				"bar",
				"text",
				"texts:strings",
				"--module",
				"foo",
				"--ack",
				"foo:string,bar:int,baz:bool",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a packet with no module specified",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "packet", "--yes", "bar", "text"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a packet in a non existent module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "packet", "--yes", "bar", "text", "--module", "nomodule"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing packet",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "packet", "--yes", "bar", "post", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	env.Must(env.Exec("create a packet with custom type fields",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"packet",
				"--yes",
				"ticket",
				"numInt:int",
				"numsInt:array.int",
				"numsIntAlias:ints",
				"numUint:uint",
				"numsUint:array.uint",
				"numsUintAlias:uints",
				"textString:string",
				"textStrings:array.string",
				"textStringsAlias:strings",
				"textCoin:coin",
				"textCoins:array.coin",
				"textCoinsAlias:coins",
				"victory:bool",
				"--module",
				"foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "type", "--yes", "custom-type", "customField:uint", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a packet with a custom field type",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "packet", "--yes", "foo-baz", "customField:CustomType", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a packet with no send message",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "packet", "--yes", "nomessage", "foo", "--no-message", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a packet with no field",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "packet", "--yes", "empty", "--module", "foo"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a non-IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "module", "--yes", "bar", "--require-registration"),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("should prevent creating a packet in a non IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "s", "packet", "--yes", "foo", "text", "--module", "bar"),
			step.Workdir(app.SourcePath()),
		)),
		envtest.ExecShouldError(),
	))

	app.EnsureSteady()
}
