package gno

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	gnopkg "github.com/gnolang/gno/gnovm/pkg/packages"
	"gotest.tools/v3/assert"

	"github.com/ignite/cli/v30/ignite/pkg/env"
)

// TestTestRunsRealmTests runs the real gno tests of a scaffolded realm
// through the Test entry point.
func TestTestRunsRealmTests(t *testing.T) {
	t.Setenv(env.ConfigDirEnvVar, t.TempDir())

	dir, _, err := Scaffold(KindRealm, "testcounter", ScaffoldOptions{Path: filepath.Join(t.TempDir(), "testcounter")})
	assert.NilError(t, err)
	t.Chdir(dir)

	var stdout, stderr bytes.Buffer
	assert.NilError(t, Test([]string{"."}, &stdout, &stderr))
	assert.Assert(t, strings.Contains(stderr.String(), "--- PASS"), "expected passing tests, got: %s", stderr.String())
}

func TestAbsDirs(t *testing.T) {
	dirs, err := absDirs([]string{"."})
	assert.NilError(t, err)
	assert.Assert(t, len(dirs) == 1)
	assert.Assert(t, dirs[0][0] == '/', "expected absolute path, got %s", dirs[0])
}

func TestFilterTargetPkgs(t *testing.T) {
	pkgs := gnopkg.PkgList{
		{Dir: "/w/realm", Match: []string{"./..."}},
		{Dir: "/w/realm/deep", Match: []string{"./..."}},
		{Dir: "/other/dep"}, // dependency, no Match
	}

	targets := filterTargetPkgs(pkgs, []string{"/w"})
	assert.Equal(t, 2, len(targets))

	targets = filterTargetPkgs(pkgs, []string{"/w/realm"})
	assert.Equal(t, 2, len(targets), "subdirectories are included")

	targets = filterTargetPkgs(pkgs, []string{"/nowhere"})
	assert.Equal(t, 0, len(targets))
}

func TestExprString(t *testing.T) {
	tt := []struct {
		src  string
		want string
	}{
		{"string", "string"},
		{"foo.Bar", "foo.Bar"},
		{"*Foo", "*Foo"},
		{"[]int", "[]int"},
		{"[4]byte", "[4]byte"},
		{"map[string]int", "map[string]int"},
		{"...string", "...string"},
	}
	for _, tc := range tt {
		if strings.HasPrefix(tc.src, "...") {
			// variadic types only parse in parameter position
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "x.gno", []byte("package p\nfunc f(a "+tc.src+") {}\n"), 0)
			assert.NilError(t, err, tc.src)
			fd := f.Decls[0].(*ast.FuncDecl)
			assert.Equal(t, tc.want, exprString(fd.Type.Params.List[0].Type), tc.src)
			continue
		}
		expr := parseExpr(t, tc.src)
		assert.Equal(t, tc.want, exprString(expr), tc.src)
	}
}

// parseExpr parses a single type expression source string.
func parseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	const preamble = "package p\nvar _ "
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.gno", []byte(preamble+src), 0)
	assert.NilError(t, err, src)
	return f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Type
}
