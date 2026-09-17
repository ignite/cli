package gno

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ignite/cli/v29/ignite/pkg/env"
	"gotest.tools/v3/assert"
)

func TestStdlibsHashIsDeterministic(t *testing.T) {
	h1, err := stdlibsHash()
	assert.NilError(t, err)
	h2, err := stdlibsHash()
	assert.NilError(t, err)
	assert.Equal(t, h1, h2)
	assert.Assert(t, len(h1) == 64, "sha256 hex fingerprint")
}

func TestEnsureStdlibsExtractsTree(t *testing.T) {
	base := t.TempDir()
	t.Setenv(env.ConfigDirEnvVar, base) // redirect ignite config dir

	root, err := EnsureStdlibs()
	assert.NilError(t, err)
	assert.Assert(t, root != "")

	// the stdlib dir holds the expected packages, including .go native files.
	stdlibDir := filepath.Join(root, "gnovm", "stdlibs")
	for _, pkg := range []string{"errors", "strconv", "strings"} {
		_, err := os.Stat(filepath.Join(stdlibDir, pkg))
		assert.NilError(t, err, pkg)
	}
	_, err = os.Stat(filepath.Join(stdlibDir, "builtin", "shims.go"))
	assert.NilError(t, err)

	// a second call is a no-op success (cached).
	root2, err := EnsureStdlibs()
	assert.NilError(t, err)
	assert.Equal(t, root, root2)
}
