package xos_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/xos"
)

func TestRename(t *testing.T) {
	var (
		dir     = t.TempDir()
		oldpath = path.Join(dir, "old")
		newpath = path.Join(dir, "new")
		require = require.New(t)
	)
	err := os.WriteFile(oldpath, []byte("foo"), os.ModePerm)
	require.NoError(err)

	err = xos.Rename(oldpath, newpath)

	require.NoError(err)
	bz, err := os.ReadFile(newpath)
	require.NoError(err)
	require.Equal([]byte("foo"), bz)
	_, err = os.Open(oldpath)
	require.EqualError(err, fmt.Sprintf("open %s: no such file or directory", oldpath))
}
