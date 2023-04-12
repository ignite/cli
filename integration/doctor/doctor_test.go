package doctor_test

import (
	_ "embed"
	"testing"

	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"

	"github.com/ignite/cli/ignite/config"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
	envtest "github.com/ignite/cli/integration"
)

func TestDoctor(t *testing.T) {
	// Ensure ignite binary is compiled
	envtest.New(t)
	// Prepare params
	params := testscript.Params{
		Setup: func(env *testscript.Env) error {
			env.Vars = append(env.Vars,
				// Pass ignite binary path
				"IGNITE="+envtest.IgniteApp,
				// Pass ignite config dir
				// (testscript resets envs so even if envtest.New has properly set
				// IGNT_CONFIG_DIR, we need to set it again)
				"IGNT_CONFIG_DIR="+xfilepath.MustInvoke(config.DirPath),
			)
			return nil
		},
		Dir: "testdata",
	}
	// Add other setup for go environment
	if err := gotooltest.Setup(&params); err != nil {
		t.Fatal(err)
	}
	// Run all scripts from testdata
	testscript.Run(t, params)
}
