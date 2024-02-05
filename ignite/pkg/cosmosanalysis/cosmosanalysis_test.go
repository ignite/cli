package cosmosanalysis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
)

var (
	expectedInterface = []string{"foo", "bar", "foobar"}

	file1 = []byte(`
package foo

type Foo struct {}
func (f Foo) foo() {}
func (f Foo) bar() {}
func (f Foo) foobar() {}

type Bar struct {}
func (b *Bar) foo() {}
func (b *Bar) bar() {}
func (b *Bar) barfoo() {}
`)

	file2 = []byte(`
package foo

type Foobar struct {}
func (f Foobar) foo() {}
func (f Foobar) bar() {}
func (f Foobar) foobar() {}
func (f Foobar) barfoo() {}

type Generic[T any] struct {
	i T
}
func (Generic[T]) foo(){}
func (Generic[T]) bar() {}
func (Generic[T]) foobar() {}

type GenericP[T any] struct {
	i T
}
func (*GenericP[T]) foo(){}
func (*GenericP[T]) bar() {}
func (*GenericP[T]) foobar() {}
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
func (app *App) RegisterAPIRoutes()                                   {}
func (app *App) RegisterTxService()                                   {}
func (app *App) AppCodec() codec.Codec                                { return app.appCodec }
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey        { return nil }
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey { return nil }
func (app *App) kvStoreKeys() map[string]*storetypes.KVStoreKey       { return nil }
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace   { return subspace }
func (app *App) TxConfig() client.TxConfig                       { return nil }
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
func (app *App) RegisterAPIRoutes()                                   {}
func (app *App) RegisterTxService()                                   {}
func (app *App) AppCodec() codec.Codec                                { return app.appCodec }
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey        { return nil }
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey { return nil }
func (app *App) kvStoreKeys() map[string]*storetypes.KVStoreKey       { return nil }
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace   { return subspace }
func (app *App) TxConfig() client.TxConfig                       { return nil }
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
`)
	appFileSDKv47 = []byte(`
package app
type App struct {}
func (app *App) Name() string { return app.BaseApp.Name() }
func (app *App) RegisterAPIRoutes()                                   {}
func (app *App) RegisterTxService()                                   {}
func (app *App) AppCodec() codec.Codec                                { return app.appCodec }
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey        { return nil }
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey { return nil }
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace   { return subspace }
func (app *App) TxConfig() client.TxConfig                       { return nil }
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
	err := os.WriteFile(f1, file1, 0o644)
	require.NoError(t, err)

	f2 := filepath.Join(tmpDir, "2.go")
	err = os.WriteFile(f2, file2, 0o644)
	require.NoError(t, err)

	// find in dir
	found, err := cosmosanalysis.FindImplementation(tmpDir, expectedInterface)
	require.NoError(t, err)
	require.ElementsMatch(t, found, []string{"Foo", "Foobar", "Generic", "GenericP"})

	// empty directory
	emptyDir := t.TempDir()
	found, err = cosmosanalysis.FindImplementation(emptyDir, expectedInterface)
	require.NoError(t, err)
	require.Empty(t, found)

	// can't provide file
	_, err = cosmosanalysis.FindImplementation(filepath.Join(tmpDir, "1.go"), expectedInterface)
	require.Error(t, err)
}

func TestFindImplementationInSpreadInMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	f1 := filepath.Join(tmpDir, "1.go")
	err := os.WriteFile(f1, partialImplementation, 0o644)
	require.NoError(t, err)
	f2 := filepath.Join(tmpDir, "2.go")
	err = os.WriteFile(f2, restOfImplementation, 0o644)
	require.NoError(t, err)

	found, err := cosmosanalysis.FindImplementation(tmpDir, expectedInterface)
	require.NoError(t, err)
	require.Len(t, found, 1)
	require.Contains(t, found, "Foo")
}

func TestFindImplementationNotFound(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	noImplFile := filepath.Join(tmpDir1, "1.go")
	err := os.WriteFile(noImplFile, noImplementation, 0o644)
	require.NoError(t, err)
	partialImplFile := filepath.Join(tmpDir2, "2.go")
	err = os.WriteFile(partialImplFile, partialImplementation, 0o644)
	require.NoError(t, err)

	// No implementation
	found, err := cosmosanalysis.FindImplementation(tmpDir1, expectedInterface)
	require.NoError(t, err)
	require.Len(t, found, 0)

	// Partial implementation
	found, err = cosmosanalysis.FindImplementation(tmpDir2, expectedInterface)
	require.NoError(t, err)
	require.Len(t, found, 0)
}

func TestFindAppFilePath(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	appFolder1 := filepath.Join(tmpDir1, "app")
	appFolder2 := filepath.Join(tmpDir1, "myOwnAppDir")
	appFolder3 := filepath.Join(tmpDir2, "sdk47AppDir")
	err := os.Mkdir(appFolder1, 0o700)
	require.NoError(t, err)
	err = os.Mkdir(appFolder2, 0o700)
	require.NoError(t, err)
	err = os.Mkdir(appFolder3, 0o700)
	require.NoError(t, err)

	// No file
	_, err = cosmosanalysis.FindAppFilePath(tmpDir1)
	require.Equal(t, "app.go file cannot be found", err.Error())

	// Only one file with app implementation
	myOwnAppFilePath := filepath.Join(appFolder2, "my_own_app.go")
	err = os.WriteFile(myOwnAppFilePath, appFile, 0o644)
	require.NoError(t, err)
	pathFound, err := cosmosanalysis.FindAppFilePath(tmpDir1)
	require.NoError(t, err)
	require.Equal(t, myOwnAppFilePath, pathFound)

	// With a test file added
	appTestFilePath := filepath.Join(appFolder2, "my_own_app_test.go")
	err = os.WriteFile(appTestFilePath, appTestFile, 0o644)
	require.NoError(t, err)
	_, err = cosmosanalysis.FindAppFilePath(tmpDir1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot locate your app.go")

	// With an additional app file (that is app.go)
	appFilePath := filepath.Join(appFolder1, "app.go")
	err = os.WriteFile(appFilePath, appFile, 0o644)
	require.NoError(t, err)
	pathFound, err = cosmosanalysis.FindAppFilePath(tmpDir1)
	require.NoError(t, err)
	require.Equal(t, appFilePath, pathFound)

	// With two app.go files
	extraAppFilePath := filepath.Join(appFolder2, "app.go")
	err = os.WriteFile(extraAppFilePath, appFile, 0o644)
	require.NoError(t, err)
	pathFound, err = cosmosanalysis.FindAppFilePath(tmpDir1)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(appFolder1, "app.go"), pathFound)

	// With an app.go file from a Cosmos SDK v0.47 app
	sdk47AppFilePath := filepath.Join(appFolder3, "app_sdk_47.go")
	err = os.WriteFile(sdk47AppFilePath, appFileSDKv47, 0o644)
	require.NoError(t, err)
	pathFound, err = cosmosanalysis.FindAppFilePath(tmpDir2)
	require.NoError(t, err)
	require.Equal(t, sdk47AppFilePath, pathFound)
}

func TestIsChainPath(t *testing.T) {
	err := cosmosanalysis.IsChainPath(".")
	require.ErrorAs(t, err, &cosmosanalysis.ErrPathNotChain{})

	err = cosmosanalysis.IsChainPath("testdata/chain")
	require.NoError(t, err)

	// testdata/chain-sdk-fork is a chain using a fork of the Cosmos SDK
	// so it should still be considered as a chain as ValidateGoMod
	// does not resolve the module file replacement.
	err = cosmosanalysis.IsChainPath("testdata/chain-sdk-fork")
	require.NoError(t, err)
}

func TestValidateGoMod(t *testing.T) {
	modFile, err := gomodule.ParseAt("testdata/chain")
	require.NoError(t, err)
	err = cosmosanalysis.ValidateGoMod(modFile)
	require.NoError(t, err)

	modFile, err = gomodule.ParseAt("testdata/chain-sdk-fork")
	require.NoError(t, err)
	err = cosmosanalysis.ValidateGoMod(modFile)
	require.NoError(t, err)
}
