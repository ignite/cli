package gno

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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

func TestPrintServeBanner(t *testing.T) {
	opts := ServeOptions{ChainID: "dev", RPCListener: "tcp://127.0.0.1:26657"}

	t.Run("lists every deployed package", func(t *testing.T) {
		out := &bytes.Buffer{}
		printServeBanner(out, opts, []string{"gno.land/p/utils", "gno.land/r/counter"})

		assert.Assert(t, strings.Contains(out.String(), "gno.land/p/utils"), "banner should list utils: %s", out.String())
		assert.Assert(t, strings.Contains(out.String(), "gno.land/r/counter"), "banner should list counter: %s", out.String())
		assert.Assert(t, !strings.Contains(out.String(), "No gno package found"), "banner should not show the empty hint: %s", out.String())
	})

	t.Run("hints when no package is deployed", func(t *testing.T) {
		out := &bytes.Buffer{}
		printServeBanner(out, opts, nil)

		assert.Assert(t, strings.Contains(out.String(), "No gno package found"), "banner should show the empty hint: %s", out.String())
		assert.Assert(t, strings.Contains(out.String(), "gnowork.toml"), "empty hint should mention gnowork.toml: %s", out.String())
	})
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
