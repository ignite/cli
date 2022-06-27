package cosmosanalysis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis"
)

var (
	expectedinterface = []string{"foo", "bar", "foobar"}

	file1 = []byte(`
package foo

type Foo struct {}
func (f Foo) foo() {}
func (f Foo) bar() {}
func (f Foo) foobar() {}

type Bar struct {}
func (b Bar) foo() {}
func (b Bar) bar() {}
func (b Bar) barfoo() {}
`)

	file2 = []byte(`
package foo

type Foobar struct {}
func (f Foobar) foo() {}
func (f Foobar) bar() {}
func (f Foobar) foobar() {}
func (f Foobar) barfoo() {}
`)
	noImplementation = []byte(`
package foo
type Foo struct {}
func (f Foo) nofoo() {}
func (f Foo) nobar() {}
func (f Foo) nofoobar() {}
`)

	partialImplementation = []byte(`
package foo
type Foo struct {}
func (f Foo) foo() {}
func (f Foo) bar() {}
`)
	restOfImplementation = []byte(`
package foo
func (f Foo) foobar() {}
`)

	appFile = []byte(`
package app
type App struct {}
func (app *App) Name() string { return app.BaseApp.Name() }
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
`)
	appTestFile = []byte(`
package app_test
type App struct {}
func (app *App) Name() string { return app.BaseApp.Name() }
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
`)
)

func TestFindImplementation(t *testing.T) {
	tmpDir := t.TempDir()

	f1 := filepath.Join(tmpDir, "1.go")
	err := os.WriteFile(f1, file1, 0644)
	require.NoError(t, err)

	f2 := filepath.Join(tmpDir, "2.go")
	err = os.WriteFile(f2, file2, 0644)
	require.NoError(t, err)

	// find in dir
	found, err := cosmosanalysis.FindImplementation(tmpDir, expectedinterface)
	require.NoError(t, err)
	require.Len(t, found, 2)
	require.Contains(t, found, "Foo")
	require.Contains(t, found, "Foobar")

	// empty directory
	emptyDir := t.TempDir()
	found, err = cosmosanalysis.FindImplementation(emptyDir, expectedinterface)
	require.NoError(t, err)
	require.Empty(t, found)

	// can't provide file
	_, err = cosmosanalysis.FindImplementation(filepath.Join(tmpDir, "1.go"), expectedinterface)
	require.Error(t, err)
}

func TestFindImplementationInSpreadInMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	f1 := filepath.Join(tmpDir, "1.go")
	err := os.WriteFile(f1, partialImplementation, 0644)
	require.NoError(t, err)
	f2 := filepath.Join(tmpDir, "2.go")
	err = os.WriteFile(f2, restOfImplementation, 0644)
	require.NoError(t, err)

	found, err := cosmosanalysis.FindImplementation(tmpDir, expectedinterface)
	require.NoError(t, err)
	require.Len(t, found, 1)
	require.Contains(t, found, "Foo")
}

func TestFindImplementationNotFound(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	noImplFile := filepath.Join(tmpDir1, "1.go")
	err := os.WriteFile(noImplFile, noImplementation, 0644)
	require.NoError(t, err)
	partialImplFile := filepath.Join(tmpDir2, "2.go")
	err = os.WriteFile(partialImplFile, partialImplementation, 0644)
	require.NoError(t, err)

	// No implementation
	found, err := cosmosanalysis.FindImplementation(tmpDir1, expectedinterface)
	require.Len(t, found, 0)

	// Partial implementation
	found, err = cosmosanalysis.FindImplementation(tmpDir2, expectedinterface)
	require.Len(t, found, 0)
}

func TestFindAppFilePath(t *testing.T) {
	tmpDir := t.TempDir()

	appFolder := filepath.Join(tmpDir, "app")
	secondaryAppFolder := filepath.Join(tmpDir, "myOwnAppDir")
	err := os.Mkdir(appFolder, 0700)
	require.NoError(t, err)
	err = os.Mkdir(secondaryAppFolder, 0700)
	require.NoError(t, err)

	// No file
	_, err = cosmosanalysis.FindAppFilePath(tmpDir)
	require.Equal(t, "app.go file cannot be found", err.Error())

	// Only one file with app implementation
	myOwnAppFilePath := filepath.Join(secondaryAppFolder, "my_own_app.go")
	err = os.WriteFile(myOwnAppFilePath, appFile, 0644)
	require.NoError(t, err)
	pathFound, err := cosmosanalysis.FindAppFilePath(tmpDir)
	require.NoError(t, err)
	require.Equal(t, myOwnAppFilePath, pathFound)

	// With a test file added
	appTestFilePath := filepath.Join(secondaryAppFolder, "my_own_app_test.go")
	err = os.WriteFile(appTestFilePath, appTestFile, 0644)
	require.NoError(t, err)
	pathFound, err = cosmosanalysis.FindAppFilePath(tmpDir)
	require.Contains(t, err.Error(), "cannot locate your app.go")

	// With an additional app file (that is app.go)
	appFilePath := filepath.Join(appFolder, "app.go")
	err = os.WriteFile(appFilePath, appFile, 0644)
	require.NoError(t, err)
	pathFound, err = cosmosanalysis.FindAppFilePath(tmpDir)
	require.NoError(t, err)
	require.Equal(t, appFilePath, pathFound)

	// With two app.go files
	extraAppFilePath := filepath.Join(secondaryAppFolder, "app.go")
	err = os.WriteFile(extraAppFilePath, appFile, 0644)
	require.NoError(t, err)
	pathFound, err = cosmosanalysis.FindAppFilePath(tmpDir)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(appFolder, "app.go"), pathFound)
}
