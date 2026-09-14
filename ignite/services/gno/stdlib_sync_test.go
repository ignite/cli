package gno

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/mod/modfile"
)

// gnoModule is the module whose stdlibs are embedded.
const gnoModule = "github.com/gnolang/gno"

// gnoModuleVersion returns the gno version pinned in the ignite go.mod.
func gnoModuleVersion(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	mod, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		t.Fatalf("parsing go.mod: %v", err)
	}
	for _, r := range mod.Require {
		if r.Mod.Path == gnoModule {
			return r.Mod.Version
		}
	}
	t.Fatalf("%s not found in go.mod", gnoModule)
	return ""
}

// gnoModuleDir downloads the gno module at the pinned version and returns
// its source directory in the module cache.
func gnoModuleDir(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("go", "mod", "download", "-json", gnoModule)
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("downloading %s: %v", gnoModule, err)
	}
	// minimal extraction of the "Dir" field without a JSON dependency
	for _, line := range strings.Split(string(out), "\n") {
		if d, ok := strings.CutPrefix(strings.TrimSpace(line), "\"Dir\":"); ok {
			return strings.Trim(strings.TrimSpace(d), "\",")
		}
	}
	t.Fatalf("no Dir in `go mod download -json %s` output", gnoModule)
	return ""
}

// hashTree returns a deterministic fingerprint of an fs.FS subtree (file
// names relative to root, and contents), used to detect drift between the
// embedded stdlibs and the gno module sources.
func hashTree(fsys fs.FS, root string) (string, error) {
	h := sha256.New()
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		h.Write([]byte(filepath.ToSlash(rel)))
		h.Write(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// TestStdlibsInSync guards the embedded gno stdlibs against drift: they must
// be a byte-for-byte copy of gnovm/stdlibs from the gno module version
// pinned in go.mod. When this fails after a gno version bump, run:
//
//	make sync-gno-stdlibs
func TestStdlibsInSync(t *testing.T) {
	version := gnoModuleVersion(t)
	modDir := gnoModuleDir(t)

	want, err := hashTree(os.DirFS(modDir), "gnovm/stdlibs")
	if err != nil {
		t.Fatalf("hashing module stdlibs: %v", err)
	}
	got, err := hashTree(stdlibsFS, stdlibsRoot)
	if err != nil {
		t.Fatalf("hashing embedded stdlibs: %v", err)
	}

	if got != want {
		t.Errorf(
			"embedded stdlibs (%s) are out of sync with %s@%s\ncurrent hash:   %s\nexpected hash:  %s\n\nrun `make sync-gno-stdlibs` and commit the result",
			stdlibsRoot, gnoModule, version, got, want,
		)
	} else {
		t.Logf("embedded stdlibs match %s@%s (%s)", gnoModule, version, fmtShortHash(got))
	}
}

func fmtShortHash(h string) string {
	if len(h) > 12 {
		return h[:12]
	}
	return h
}
