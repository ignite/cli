//go:build !relayer

package other_components_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	envtest "github.com/ignite/cli/v29/integration"
)

func TestGenerateAnAppWithModuleMigrations(t *testing.T) {
	var (
		env = envtest.New(t)
		app = env.ScaffoldApp("github.com/test/blog")
	)

	app.Scaffold(
		"create the first module migration",
		false,
		"migration",
		"blog",
	)

	requireMigrationFile(t, app.SourcePath(), "blog", "v2")

	moduleFilePath := filepath.Join(app.SourcePath(), "x", "blog", "module", "module.go")
	moduleContent, err := os.ReadFile(moduleFilePath)
	require.NoError(t, err)

	normalized := normalizeWhitespace(string(moduleContent))
	require.Contains(t, normalized, `migrationv2"github.com/test/blog/x/blog/migrations/v2"`)
	require.Contains(t, normalized, `cfg,ok:=registrar.(module.Configurator)`)
	require.Contains(t, normalized, `if!ok{returnnil}`)
	require.Contains(t, normalized, `cfg.RegisterMigration(types.ModuleName,1,migrationv2.Migrate)`)
	require.Contains(t, normalized, `func(AppModule)ConsensusVersion()uint64{return2}`)

	app.Scaffold(
		"create the second module migration",
		false,
		"migration",
		"blog",
	)

	requireMigrationFile(t, app.SourcePath(), "blog", "v3")

	moduleContent, err = os.ReadFile(moduleFilePath)
	require.NoError(t, err)

	normalized = normalizeWhitespace(string(moduleContent))
	require.Equal(t, 1, strings.Count(normalized, `registrar.(module.Configurator)`))
	require.Contains(t, normalized, `migrationv3"github.com/test/blog/x/blog/migrations/v3"`)
	require.Contains(t, normalized, `cfg.RegisterMigration(types.ModuleName,1,migrationv2.Migrate)`)
	require.Contains(t, normalized, `cfg.RegisterMigration(types.ModuleName,2,migrationv3.Migrate)`)
	require.Contains(t, normalized, `func(AppModule)ConsensusVersion()uint64{return3}`)

	app.EnsureSteady()
}

func requireMigrationFile(t *testing.T, appPath, moduleName, version string) {
	t.Helper()

	migrationPath := filepath.Join(appPath, "x", moduleName, "migrations", version, "migrate.go")
	content, err := os.ReadFile(migrationPath)
	require.NoError(t, err)

	normalized := normalizeWhitespace(string(content))
	require.Contains(t, normalized, "package"+version)
	require.Contains(t, normalized, `funcMigrate(_sdk.Context)error{returnnil}`)
}

func normalizeWhitespace(content string) string {
	return strings.Join(strings.Fields(content), "")
}
