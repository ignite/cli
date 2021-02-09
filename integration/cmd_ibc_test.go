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
			step.Exec("starport", "module", "create", "--ibc", "foo"),
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
			step.Exec("starport", "module", "create", "--ibc", "orderedfoo", "--ibc-ordering", "ordered"),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create an IBC module with an unordered channel",
		step.NewSteps(step.New(
			step.Exec("starport", "module", "create", "--ibc", "unorderedfoo", "--ibc-ordering", "unordered"),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)
}
