package xfilepath_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/xfilepath"
	"os"
	"path/filepath"
	"testing"
)

func TestJoin(t *testing.T) {
	retriever := xfilepath.Join(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", nil),
		xfilepath.Path("foobar/barfoo"),
	)
	path, err := retriever()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(
		"foo",
		"bar",
		"foobar",
		"barfoo",
	), path)

	retriever = xfilepath.Join(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", errors.New("foo")),
		xfilepath.Path("foobar/barfoo"),
	)
	_, err = retriever()
	require.Error(t, err)
}

func TestJoinFromHome(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	retriever := xfilepath.JoinFromHome(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", nil),
		xfilepath.Path("foobar/barfoo"),
	)
	path, err := retriever()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(
		home,
		"foo",
		"bar",
		"foobar",
		"barfoo",
	), path)

	retriever = xfilepath.JoinFromHome(
		xfilepath.Path("foo"),
		xfilepath.PathWithError("bar", errors.New("foo")),
		xfilepath.Path("foobar/barfoo"),
	)
	_, err = retriever()
	require.Error(t, err)
}
