package plugin_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestAddRemovePlugin(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/test/blog")

		assertPlugins = func(expectedLocalPlugins, expectedGlobalPlugins []pluginsconfig.Plugin) {
			localCfg, err := pluginsconfig.ParseDir(app.SourcePath())
			require.NoError(err)
			assert.ElementsMatch(expectedLocalPlugins, localCfg.Apps, "unexpected local plugins")

			globalCfgPath, err := plugin.PluginsPath()
			require.NoError(err)
			globalCfg, err := pluginsconfig.ParseDir(globalCfgPath)
			require.NoError(err)
			assert.ElementsMatch(expectedGlobalPlugins, globalCfg.Apps, "unexpected global plugins")
		}
	)

	// no plugins expected
	assertPlugins(nil, nil)

	// Note: Originally plugin repo was "github.com/ignite/example-plugin" instead of a local one
	pluginRepo, err := filepath.Abs("testdata/example-plugin")
	require.NoError(err)

	env.Must(env.Exec("add plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginRepo, "k1=v1", "k2=v2"),
			step.Workdir(app.SourcePath()),
		)),
	))

	// one local plugin expected
	assertPlugins(
		[]pluginsconfig.Plugin{
			{
				Path: pluginRepo,
				With: map[string]string{
					"k1": "v1",
					"k2": "v2",
				},
			},
		},
		nil,
	)

	env.Must(env.Exec("uninstall plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "uninstall", pluginRepo),
			step.Workdir(app.SourcePath()),
		)),
	))

	// no plugins expected
	assertPlugins(nil, nil)

	env.Must(env.Exec("install plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginRepo, "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	// one global plugins expected
	assertPlugins(
		nil,
		[]pluginsconfig.Plugin{
			{
				Path: pluginRepo,
			},
		},
	)

	env.Must(env.Exec("uninstall plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "uninstall", pluginRepo, "-g"),
			step.Workdir(app.SourcePath()),
		)),
	))

	// no plugins expected
	assertPlugins(nil, nil)
}

// TODO install network plugin test

func TestPluginScaffold(t *testing.T) {
	env := envtest.New(t)

	env.Must(env.Exec("install a plugin",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "scaffold", "test"),
			step.Workdir(env.TmpDir()),
		)),
	))
}
