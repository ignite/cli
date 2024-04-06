package diff

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestComputeFS(t *testing.T) {
	require := require.New(t)

	origin := fstest.MapFS{
		"foo.txt": &fstest.MapFile{
			Data: []byte("hello"),
		},
		"bar.txt": &fstest.MapFile{
			Data: []byte("unmodified"),
		},
		"pkg/main.go": &fstest.MapFile{
			Data: []byte("package main"),
		},
	}
	modified := fstest.MapFS{
		"foo.txt": &fstest.MapFile{
			Data: []byte("world"),
		},
		"bar.txt": &fstest.MapFile{
			Data: []byte("unmodified"),
		},
		"new.txt": &fstest.MapFile{
			Data: []byte("new file"),
		},
		"pkg/main.go": &fstest.MapFile{
			Data: []byte("package main\nfunc main() {}"),
		},
	}

	unified, err := computeFS(origin, modified)
	require.NoError(err)
	require.Len(unified, 3)
	expectedFiles := []string{"foo.txt", "new.txt", "pkg/main.go"}
	for _, u := range unified {
		require.Contains(expectedFiles, u.From, "unexpected file in diff: %s", u.From)
	}

	// Test ignoring files
	unified, err = computeFS(origin, modified, "**.go")
	require.NoError(err)
	require.Len(unified, 2)
	expectedFiles = []string{"foo.txt", "new.txt"}
	for _, u := range unified {
		require.Contains(expectedFiles, u.From, "unexpected file in diff: %s", u.From)
	}
}
