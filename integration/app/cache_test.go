package app_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/integration"
)

func TestCliWithCaching(t *testing.T) {
	var (
		env               = envtest.New(t)
		path              = env.Scaffold("github.com/test/cacheblog")
		vueGenerated      = filepath.Join(path, "vue/src/store/generated")
		openapiGenerated  = filepath.Join(path, "docs/static/openapi.yml")
		typesDir          = filepath.Join(path, "x/cacheblog/types")
		servers           = env.RandomizeServerPorts(path, "")
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)

	env.Must(env.Exec("create a message",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"message",
				"mymessage",
				"myfield1",
				"myfield2:bool",
				"--yes",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("create a query",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"s",
				"query",
				"myQuery",
				"mytypefield",
				"--yes",
			),
			step.Workdir(path),
		)),
	))

	env.Must(env.Exec("build",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"c",
				"build",
				"--proto-all-modules",
			),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)

	deleteCachedFiles(t, vueGenerated, openapiGenerated, typesDir)

	env.Must(env.Exec("build",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"c",
				"build",
				"--proto-all-modules",
			),
			step.Workdir(path),
		)),
	))

	env.EnsureAppIsSteady(path)

	deleteCachedFiles(t, vueGenerated, openapiGenerated, typesDir)

	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers)
	}()
	env.Must(env.Serve("should serve with Stargate version", path, "", "", envtest.ExecCtx(ctx)))

	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}

func deleteCachedFiles(t *testing.T, vueGenerated, openapiGenerated, typesDir string) {
	require.NoError(t, os.RemoveAll(vueGenerated))
	require.NoError(t, os.Remove(openapiGenerated))

	typesDirEntries, err := os.ReadDir(typesDir)
	require.NoError(t, err)

	for _, v := range typesDirEntries {
		if v.IsDir() {
			continue
		}

		if strings.Contains(v.Name(), ".pb") {
			require.NoError(t, os.Remove(filepath.Join(typesDir, v.Name())))
		}
	}
}
