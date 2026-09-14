package gno

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func TestResolveModulePath(t *testing.T) {
	tt := []struct {
		name    string
		prefix  string
		in      string
		want    string
		wantErr bool
	}{
		{name: "bare realm name", prefix: "r", in: "counter", want: "gno.land/r/counter"},
		{name: "nested realm name", prefix: "r", in: "demo/counter", want: "gno.land/r/demo/counter"},
		{name: "full realm path", prefix: "r", in: "gno.land/r/demo/counter", want: "gno.land/r/demo/counter"},
		{name: "bare package name", prefix: "p", in: "utils", want: "gno.land/p/utils"},
		{name: "full package path", prefix: "p", in: "gno.land/p/utils", want: "gno.land/p/utils"},
		{name: "realm path for package kind", prefix: "p", in: "gno.land/r/counter", wantErr: true},
		{name: "package path for realm kind", prefix: "r", in: "gno.land/p/utils", wantErr: true},
		{name: "invalid characters", prefix: "r", in: "Foo", wantErr: true},
		{name: "empty", prefix: "r", in: "", wantErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveModulePath(tc.prefix, tc.in)
			if tc.wantErr {
				assert.ErrorContains(t, err, "")
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestScaffoldRealm(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "counter")

	dir, modulePath, err := Scaffold(KindRealm, "counter", ScaffoldOptions{Path: target})
	assert.NilError(t, err)
	assert.Equal(t, target, dir)
	assert.Equal(t, "gno.land/r/counter", modulePath)

	for _, f := range []string{"gnomod.toml", "counter.gno", "counter_test.gno"} {
		_, err := os.Stat(filepath.Join(target, f))
		assert.NilError(t, err, f)
	}

	// gnomod.toml is parseable and holds the module path.
	got, err := parseGnoModDir(target)
	assert.NilError(t, err)
	assert.Equal(t, "gno.land/r/counter", got)

	// scaffolding into a non-empty dir fails
	_, _, err = Scaffold(KindRealm, "counter", ScaffoldOptions{Path: target})
	assert.ErrorContains(t, err, "not empty")
}

func TestScaffoldPackage(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "utils")

	_, modulePath, err := Scaffold(KindPackage, "utils", ScaffoldOptions{Path: target})
	assert.NilError(t, err)
	assert.Equal(t, "gno.land/p/utils", modulePath)

	got, err := parseGnoModDir(target)
	assert.NilError(t, err)
	assert.Equal(t, "gno.land/p/utils", got)
}
