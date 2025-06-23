package cosmosgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ettle/strcase"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
)

func TestGenerateTypeScript(t *testing.T) {
	require := require.New(t)
	testdataDir := "testdata"
	appDir := filepath.Join(testdataDir, "testchain")
	tsClientDir := filepath.Join(appDir, "ts-client")

	cacheStorage, err := cache.NewStorage(filepath.Join(t.TempDir(), "cache.db"))
	require.NoError(err)

	buf, err := cosmosbuf.New(cacheStorage, t.Name())
	require.NoError(err)

	// Use module discovery to collect test module proto.
	m, err := module.Discover(t.Context(), appDir, appDir, module.WithProtoDir("proto"))
	require.NoError(err, "failed to discover module")
	require.Len(m, 1, "expected exactly one module to be discovered")

	g := newTSGenerator(&generator{
		appPath:      appDir,
		protoDir:     "proto",
		goModPath:    "go.mod",
		cacheStorage: cacheStorage,
		buf:          buf,
		appModules:   m,
		opts: &generateOptions{
			tsClientRootPath: tsClientDir,
			useCache:         false,
			jsOut: func(m module.Module) string {
				return filepath.Join(tsClientDir, fmt.Sprintf("%s.%s.%s", "ignite", "planet", strcase.ToKebab(m.Name)))
			},
		},
	})

	err = g.generateModuleTemplate(t.Context(), appDir, m[0])
	require.NoError(err, "failed to generate TypeScript files")

	err = g.generateRootTemplates(generatePayload{
		Modules:   m,
		PackageNS: strings.ReplaceAll(appDir, "/", "-"),
	})
	require.NoError(err)

	// compare all generated files to golden files
	goldenDir := filepath.Join(testdataDir, "expected_files", "ts-client")
	_ = filepath.Walk(goldenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(goldenDir, path)
		got := filepath.Join(tsClientDir, rel)
		gold, err := os.ReadFile(path)
		require.NoError(err, "failed to read golden file: %s", path)

		gotBytes, err := os.ReadFile(got)
		require.NoError(err, "failed to read generated file: %s", got)
		require.Equal(string(gold), string(gotBytes), "file %s does not match golden file", rel)

		return nil
	})
}
