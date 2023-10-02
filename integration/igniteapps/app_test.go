package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appsconfig "github.com/ignite/cli/ignite/config/apps"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/services/plugin"
	envtest "github.com/ignite/cli/integration"
)

func TestAddRemoveApp(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/blog")
		appRepo = "github.com/ignite/example-app"

		assertApps = func(expectedLocalApps, expectedGlobalApps []appsconfig.App) {
			localCfg, err := appsconfig.ParseDir(app.SourcePath())
			require.NoError(err)
			assert.ElementsMatch(expectedLocalApps, localCfg.Apps, "unexpected local apps")

			globalCfgPath, err := plugin.PluginsPath()
			require.NoError(err)
			globalCfg, err := appsconfig.ParseDir(globalCfgPath)
			require.NoError(err)
			assert.ElementsMatch(expectedGlobalApps, globalCfg.Apps, "unexpected global apps")
		}
	)

	// no apps expected
	assertApps(nil, nil)

	env.Must(env.Exec("install app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", appRepo, "k1=v1", "k2=v2"),
			step.Workdir(app.SourcePath()),
		)),
	))

	// one local app expected
	assertApps(
		[]appsconfig.App{
			{
				Path: appRepo,
				With: map[string]string{
					"k1": "v1",
					"k2": "v2",
				},
			},
		},
		nil,
	)

	env.Must(env.Exec("uninstall app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "uninstall", appRepo),
			step.Workdir(app.SourcePath()),
		)),
	))

	// no apps expected
	assertApps(nil, nil)

	env.Must(env.Exec("install app globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", appRepo, "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	// one global apps expected
	assertApps(
		nil,
		[]appsconfig.App{
			{
				Path: appRepo,
			},
		},
	)

	env.Must(env.Exec("uninstall app globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "uninstall", appRepo, "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	// no apps expected
	assertApps(nil, nil)
}

func TestAppScaffold(t *testing.T) {
	env := envtest.New(t)

	env.Must(env.Exec("install a app",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "scaffold", "test"),
			step.Workdir(env.TmpDir()),
		)),
	))
}
