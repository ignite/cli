package gno

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
	"gotest.tools/v3/assert"
)

func TestPreloadPaths(t *testing.T) {
	tt := []struct {
		name  string
		files map[string]string
		want  []string
	}{
		{
			name:  "gnomod with module",
			files: map[string]string{"gnomod.toml": "module = \"gno.land/r/x\"\ngno = \"0.9\"\n"},
			want:  []string{"gno.land/r/x"},
		},
		{
			name:  "no gnomod",
			files: map[string]string{},
			want:  nil,
		},
		{
			name:  "invalid gnomod",
			files: map[string]string{"gnomod.toml": "not toml <<<"},
			want:  nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, content := range tc.files {
				assert.NilError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
			}
			assert.DeepEqual(t, tc.want, preloadPaths(dir))
		})
	}
}

func TestIsGnoFile(t *testing.T) {
	assert.Assert(t, isGnoFile("/x/realm.gno"))
	assert.Assert(t, isGnoFile("/x/gnomod.toml"))
	assert.Assert(t, !isGnoFile("/x/main.go"))
	assert.Assert(t, !isGnoFile("/x/readme.md"))
}

func TestWatchDirRecursiveSkipsDotDirs(t *testing.T) {
	dir := t.TempDir()
	assert.NilError(t, os.MkdirAll(filepath.Join(dir, "sub", ".git"), 0o755))
	assert.NilError(t, os.WriteFile(filepath.Join(dir, "a.gno"), []byte("x"), 0o644))

	fw, err := fsnotify.NewWatcher()
	assert.NilError(t, err)
	defer fw.Close()

	assert.NilError(t, watchDirRecursive(fw, dir))

	watched := map[string]bool{}
	for _, w := range fw.WatchList() {
		watched[w] = true
	}
	assert.Assert(t, watched[dir])
	assert.Assert(t, watched[filepath.Join(dir, "sub")])
	assert.Assert(t, !watched[filepath.Join(dir, "sub", ".git")], ".git should be skipped")
}
