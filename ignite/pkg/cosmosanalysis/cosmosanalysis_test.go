package cosmosanalysis_test

import (
	"go/ast"
	"go/parser"
	"go/token"
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

import (
	runtime "github.com/cosmos/cosmos-sdk/runtime"
)

type App struct {
	*runtime.App
}`)
	appTestFile = []byte(`
package app_test

import (
	runtime "github.com/cosmos/cosmos-sdk/runtime"
)

type App struct {
	*runtime.App
}`)

	appFileSDKv47 = []byte(`
package app

import "github.com/cosmos/cosmos-sdk/baseapp"

type App struct {
	baseapp.BaseApp
}`)

	embeddedTypeFile = []byte(`
package foo

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	runtime "github.com/cosmos/cosmos-sdk/runtime"
)

type App1 struct {
	*runtime.App
}

type App2 struct {
	baseapp.BaseApp
}

type NotApp struct{}

// AppPointer uses a pointer to an embedded type
type AppPointer struct {
	*runtime.App
}

// AppNoEmbed has the type but doesn't embed it
type AppNoEmbed struct {
	a runtime.App
}

// OtherEmbed embeds a different type from a target package
type OtherEmbed struct {
	*runtime.Server
}
`)
	appModuleGoMod = []byte(`
module example.com/foo

go 1.19

require (
	github.com/cosmos/cosmos-sdk v0.47.0
)`)
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

func TestFindEmbed(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a dummy go.mod for the test package
	modPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(modPath, appModuleGoMod, 0o644)
	require.NoError(t, err)

	// Create the test file
	filePath := filepath.Join(tmpDir, "app.go")
	err = os.WriteFile(filePath, embeddedTypeFile, 0o644)
	require.NoError(t, err)

	targets := []string{
		"github.com/cosmos/cosmos-sdk/runtime.App",
		"github.com/cosmos/cosmos-sdk/baseapp.BaseApp",
	}

	found, err := cosmosanalysis.FindEmbed(tmpDir, targets)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"App1", "App2", "AppPointer"}, found)

	// Test with a directory that doesn't contain the target embeds
	emptyDir := t.TempDir()
	modPathEmpty := filepath.Join(emptyDir, "go.mod")
	err = os.WriteFile(modPathEmpty, appModuleGoMod, 0o644)
	require.NoError(t, err)
	otherFilePath := filepath.Join(emptyDir, "other.go")
	err = os.WriteFile(otherFilePath, []byte(`package foo; type Bar struct{}`), 0o644)
	require.NoError(t, err)

	foundEmpty, err := cosmosanalysis.FindEmbed(emptyDir, targets)
	require.NoError(t, err)
	require.Empty(t, foundEmpty)

	// Test with non-existent directory
	_, err = cosmosanalysis.FindEmbed(filepath.Join(tmpDir, "nonexistent"), targets)
	require.Error(t, err) // Expect an error because parser.ParseDir will fail
}

func TestFindEmbedInFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "app.go")
	err := os.WriteFile(filePath, embeddedTypeFile, 0o644)
	require.NoError(t, err)

	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, filePath, nil, 0)
	require.NoError(t, err)

	targets := []string{
		"github.com/cosmos/cosmos-sdk/runtime.App",
		"github.com/cosmos/cosmos-sdk/baseapp.BaseApp",
		"github.com/cosmos/cosmos-sdk/runtime.Server", // To test OtherEmbed
	}

	found := cosmosanalysis.FindEmbedInFile(fileNode, targets)
	require.ElementsMatch(t, []string{"App1", "App2", "AppPointer", "OtherEmbed"}, found)

	// Test with a node that is not an *ast.File (though the function expects it)
	// Create a simple ident node
	identNode := ast.NewIdent("SomeIdent")
	foundNotFile := cosmosanalysis.FindEmbedInFile(identNode, targets)
	require.Empty(t, foundNotFile) // Expect empty as it's not a file node

	// Test with a file that doesn't import/embed the target types
	otherContent := `package bar; type Bar struct{}`
	otherFilePath := filepath.Join(tmpDir, "other.go")
	err = os.WriteFile(otherFilePath, []byte(otherContent), 0o644)
	require.NoError(t, err)
	otherFileNode, err := parser.ParseFile(fset, otherFilePath, nil, 0)
	require.NoError(t, err)
	foundOther := cosmosanalysis.FindEmbedInFile(otherFileNode, targets)
	require.Empty(t, foundOther)
}
