package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/app"
)

var (
	source = []byte(`
package foo

type Foo struct {
	FooKeeper foo.keeper
}

func (f Foo) RegisterAPIRoutes() {}
func (f Foo) RegisterGRPCServer() {}
func (f Foo) RegisterTxService() {}
func (f Foo) RegisterTendermintService() {}
`)

	sourceNoApp = []byte(`
package foo

type Bar struct {
	FooKeeper foo.keeper
}
`)
)

func TestCheckKeeper(t *testing.T) {
	tmpDir := os.TempDir()

	// Test with a source file containing an app
	tmpFile := filepath.Join(tmpDir, "source")
	err := os.WriteFile(tmpFile, source, 0644)
	require.NoError(t, err)
	t.Cleanup(func() { os.Remove(tmpFile) })

	err = app.CheckKeeper(tmpFile, "FooKeeper")
	require.NoError(t, err)
	err = app.CheckKeeper(tmpFile, "BarKeeper")
	require.Error(t, err)

	// No app in source must return an error
	tmpFileNoApp := filepath.Join(tmpDir, "source")
	err = os.WriteFile(tmpFileNoApp, sourceNoApp, 0644)
	require.NoError(t, err)
	t.Cleanup(func() { os.Remove(tmpFileNoApp) })

	err = app.CheckKeeper(tmpFileNoApp, "FooKeeper")
	require.Error(t, err)
}

