package gno

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// gnoModTmpl is the gnomod.toml written for scaffolded packages.
const gnoModTmpl = `module = "%s"
gno = "0.9"
`

// realmTmpl is a minimal stateful realm with a counter, an entrypoint that
// mutates state and a Render function for gnoweb.
const realmTmpl = `package %s

import (
	"strconv"
	"strings"
)

// counter is realm state: package-level variables are persisted on-chain.
var counter int

// Increment mutates realm state. The leading "realm" argument is the gno
// convention for methods that require access to realm context.
func Increment(_ realm) {
	counter++
}

// Get returns the current counter value.
func Get() int {
	return counter
}

// Render implements the gnoweb rendering entrypoint.
func Render(path string) string {
	return strings.Join([]string{
		"# %s",
		"",
		"Counter: " + strconv.Itoa(counter),
	}, "\n")
}
`

// realmTestTmpl covers the scaffolded entrypoints.
const realmTestTmpl = `package %s

import "testing"

func TestCounter(cur realm, t *testing.T) {
	counter = 0
	if got := Increment(cross(cur)); got != 1 {
		t.Fatalf("counter after increment = %d, want 1", got)
	}
	if Get() != 1 {
		t.Fatalf("Get() = %d, want 1", Get())
	}
}
`

// packageTmpl is a minimal stateless pure package.
const packageTmpl = `package %s

// Hello returns a greeting for name.
func Hello(name string) string {
	if name == "" {
		name = "world"
	}
	return "hello, " + name + "!"
}
`

const packageTestTmpl = `package %s

import "testing"

func TestHello(t *testing.T) {
	if got := Hello("gno"); got != "hello, gno!" {
		t.Fatalf("Hello = %q", got)
	}
}
`

// modulePathRegex validates the final segment of a package path.
var modulePathRegex = regexp.MustCompile(`^[a-z0-9/_-]+$`)

// ScaffoldKind is the kind of gno package to scaffold.
type ScaffoldKind int

const (
	// KindRealm scaffolds a stateful realm (gno.land/r/...).
	KindRealm ScaffoldKind = iota
	// KindPackage scaffolds a stateless pure package (gno.land/p/...).
	KindPackage
)

// ScaffoldOptions configures Scaffold.
type ScaffoldOptions struct {
	// Path is the target directory (default: ./<name>).
	Path string
}

// Scaffold creates a new gno realm or package. name is either a bare package
// name ("counter") or a full gno.land path ("gno.land/r/demo/counter").
// It returns the created directory and the module path.
func Scaffold(kind ScaffoldKind, name string, opts ScaffoldOptions) (dir, modulePath string, err error) {
	prefix := "r"
	testTmpl := realmTestTmpl
	var srcTmpl string
	switch kind {
	case KindRealm:
		srcTmpl = realmTmpl
		prefix = "r"
	case KindPackage:
		srcTmpl = packageTmpl
		testTmpl = packageTestTmpl
		prefix = "p"
	default:
		return "", "", fmt.Errorf("unknown scaffold kind: %d", kind)
	}

	modulePath, err = resolveModulePath(prefix, name)
	if err != nil {
		return "", "", err
	}
	pkgName := pathPkgName(modulePath)

	target := opts.Path
	if target == "" {
		target = pkgName
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return "", "", err
	}
	if entries, err := os.ReadDir(target); err == nil && len(entries) > 0 {
		return "", "", fmt.Errorf("directory %s is not empty", target)
	}

	srcContent, err := formatSource(srcTmpl, pkgName, kind)
	if err != nil {
		return "", "", err
	}

	files := map[string]string{
		"gnomod.toml":         fmt.Sprintf(gnoModTmpl, modulePath),
		pkgName + ".gno":      srcContent,
		pkgName + "_test.gno": fmt.Sprintf(testTmpl, pkgName),
	}
	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(target, fname), []byte(content), 0o644); err != nil {
			return "", "", err
		}
	}

	return target, modulePath, nil
}

// formatSource renders the source template for the package kind. The realm
// template takes the package name twice (var declarations + Render heading).
func formatSource(tmpl, pkgName string, kind ScaffoldKind) (string, error) {
	if kind == KindRealm {
		return fmt.Sprintf(tmpl, pkgName, pkgName), nil
	}
	return fmt.Sprintf(tmpl, pkgName), nil
}

// resolveModulePath turns a bare name or full path into a gno.land module
// path of the given prefix (r or p).
func resolveModulePath(prefix, name string) (string, error) {
	name = strings.Trim(name, "/")
	if name == "" {
		return "", fmt.Errorf("empty package name")
	}
	var modulePath string
	switch {
	case strings.HasPrefix(name, "gno.land/r/"):
		if prefix != "r" {
			return "", fmt.Errorf("%q is a realm path, expected a package path (gno.land/p/...)", name)
		}
		modulePath = name
	case strings.HasPrefix(name, "gno.land/p/"):
		if prefix != "p" {
			return "", fmt.Errorf("%q is a package path, expected a realm path (gno.land/r/...)", name)
		}
		modulePath = name
	default:
		if !modulePathRegex.MatchString(name) {
			return "", fmt.Errorf("invalid package name %q: allowed characters are a-z, 0-9, -, _ and /", name)
		}
		modulePath = "gno.land/" + prefix + "/" + name
	}
	if !modulePathRegex.MatchString(strings.TrimPrefix(strings.TrimPrefix(modulePath, "gno.land/"+prefix+"/"), "")) {
		return "", fmt.Errorf("invalid module path %q", modulePath)
	}
	return modulePath, nil
}

// pathPkgName returns the final segment of a module path.
func pathPkgName(modulePath string) string {
	parts := strings.Split(modulePath, "/")
	return parts[len(parts)-1]
}
