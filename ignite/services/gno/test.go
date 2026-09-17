package gno

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	gnopkg "github.com/gnolang/gno/gnovm/pkg/packages"
	"github.com/gnolang/gno/gnovm/pkg/test"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// Test runs the gno tests (_test.gno files) of the packages at dirs
// (default: the current directory), the `gno test` equivalent.
//
// Test is the integration entry point of the test runner: it executes real
// gno packages and is covered by the integration suite rather than unit
// tests.
func Test(dirs []string, stdout, stderr io.Writer) error {
	rootDir, err := EnsureStdlibs()
	if err != nil {
		return errors.Wrap(err, "preparing gno stdlibs")
	}
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	absDirs, err := absDirs(dirs)
	if err != nil {
		return err
	}

	loadConf := gnopkg.LoadConfig{
		Out:        stderr,
		Deps:       true,
		Test:       true,
		AllowEmpty: true,
	}
	pkgs, err := gnopkg.Load(loadConf, dirs...)
	if err != nil {
		return errors.Wrap(err, "loading packages")
	}

	// filter to the requested packages (deps are only loaded for the store)
	targets := filterTargetPkgs(pkgs, absDirs)
	if len(targets) == 0 {
		fmt.Fprintln(stderr, "no packages to test")
		return nil
	}

	opts := test.NewTestOptions(rootDir, stdout, stderr, pkgs)

	failed := runTargetPkgs(targets, opts, stderr)
	if failed > 0 {
		return errors.Errorf("FAIL: %d package(s)", failed)
	}
	return nil
}

// runTargetPkgs runs the tests of each target package and returns the count
// of failing packages.
func runTargetPkgs(targets gnopkg.PkgList, opts *test.TestOptions, stderr io.Writer) int {
	failed := 0
	for _, pkg := range targets {
		if !testPkg(pkg, opts, stderr) {
			failed++
		}
	}
	return failed
}

// testPkg runs the tests of one package and reports whether they passed.
func testPkg(pkg *gnopkg.Package, opts *test.TestOptions, stderr io.Writer) bool {
	pkgPath := pkg.ImportPath
	if mod, err := parseGnoModDir(pkg.Dir); err == nil {
		pkgPath = mod
	}

	started := time.Now()
	fmt.Fprintf(stderr, "=== RUN   %s\n", pkgPath)

	mpkg, err := gno.ReadMemPackage(pkg.Dir, pkgPath, gno.MPAnyAll)
	if err != nil {
		fmt.Fprintf(stderr, "--- FAIL: %s (%s)\n\t%v\n", pkgPath, time.Since(started).Round(time.Millisecond), err)
		return false
	}
	if err := test.Test(mpkg, pkg.Dir, opts); err != nil {
		fmt.Fprintf(stderr, "--- FAIL: %s (%s)\n\t%v\n", pkgPath, time.Since(started).Round(time.Millisecond), err)
		return false
	}
	fmt.Fprintf(stderr, "--- PASS: %s (%s)\n", pkgPath, time.Since(started).Round(time.Millisecond))
	return true
}

// filterTargetPkgs returns the explicitly requested packages among pkgs
// (loaded dependencies are excluded), matched by directory.
func filterTargetPkgs(pkgs gnopkg.PkgList, absDirs []string) (targets gnopkg.PkgList) {
	for _, p := range pkgs {
		if len(p.Match) == 0 {
			continue
		}
		for _, d := range absDirs {
			if p.Dir == d || strings.HasPrefix(p.Dir, d+string(filepath.Separator)) {
				targets = append(targets, p)
				break
			}
		}
	}
	return targets
}

// absDirs resolves each dir to an absolute path.
func absDirs(dirs []string) ([]string, error) {
	abs := make([]string, 0, len(dirs))
	for _, d := range dirs {
		a, err := filepath.Abs(d)
		if err != nil {
			return nil, err
		}
		abs = append(abs, a)
	}
	return abs, nil
}
