// +build !relayer

package integration_test

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestCreateModuleWithIBC(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("ibcblog", Stargate)
	)

	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "foo", "--ibc"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a type in an IBC module",
		step.NewSteps(step.New(
			step.Exec("starport", "type", "user", "email", "--module", "foo"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an IBC module with an ordered channel",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "--ibc", "orderedfoo", "--ordering", "ordered"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an IBC module with an unordered channel",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "--ibc", "unorderedfoo", "--ordering", "unordered"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an non IBC module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "foobar"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}

func TestCreateIBCPacket(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("ibcblog2", Stargate)
	)

	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "foo", "--ibc"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a packet",
		step.NewSteps(step.New(
			step.Exec("starport", "packet", "bar", "text", "--module", "foo"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a packet with no module specified",
		step.NewSteps(step.New(
			step.Exec("starport", "packet", "bar", "text"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating a packet in a non existent module",
		step.NewSteps(step.New(
			step.Exec("starport", "packet", "bar", "text", "--module", "nomodule"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("should prevent creating an existing packet",
		step.NewSteps(step.New(
			step.Exec("starport", "packet", "bar", "post", "--module", "foo"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.Must(env.Exec("create a packet with custom type fields",
		step.NewSteps(step.New(
			step.Exec("starport", "packet", "ticket", "num:int", "victory:bool", "--module", "foo"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a packet with no field",
		step.NewSteps(step.New(
			step.Exec("starport", "packet", "empty", "--module", "foo"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a non-IBC module",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "bar"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("should prevent creating a packet in a non IBC module",
		step.NewSteps(step.New(
			step.Exec("starport", "packet", "foo", "text", "--module", "bar"),
			step.Workdir(path),
		)),
		ExecShouldError(),
	))

	env.EnsureAppIsSteady(path)
}
