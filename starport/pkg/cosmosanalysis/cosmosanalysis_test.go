package cosmosanalysis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis"
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
)

func TestFindImplementation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cosmosanalysis_test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	f1 := filepath.Join(tmpDir, "1.go")
	err = os.WriteFile(f1, file1, 0644)
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
	emptyDir, err := os.MkdirTemp("", "cosmosanalysis_test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(emptyDir) })
	found, err = cosmosanalysis.FindImplementation(emptyDir, expectedinterface)
	require.NoError(t, err)
	require.Empty(t, found)

	// find in file
	found, err = cosmosanalysis.FindImplementation(filepath.Join(tmpDir, "1.go"), expectedinterface)
	require.NoError(t, err)
	require.Len(t, found, 1)
	require.Contains(t, found, "Foo")

	// no file
	_, err = cosmosanalysis.FindImplementation(filepath.Join(tmpDir, "3.go"), expectedinterface)
	require.Error(t, err)
}
