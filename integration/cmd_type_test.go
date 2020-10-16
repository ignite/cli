package integration_test

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func TestGenerateAnAppWithTypeAndVerify(t *testing.T) {
	t.Parallel()

	var (
		env  = newEnv(t)
		path = env.Scaffold("blog", Stargate)
	)

	env.Must(env.Exec("add CosmWasm module",
		step.New(
			step.Exec("starport", "type", "user", "email"),
			step.Workdir(path),
		),
	))

	env.EnsureAppIsSteady(path)
}
