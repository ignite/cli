package diff

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hexops/gotextdiff"

	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
)

type Diffs map[string][]gotextdiff.Unified

var diffIgnoreGlobs = []string{
	".git/**",
	"**.md",
	"go.sum",
	"**_test.go",
	"**.pb.go",
	"**.pb.gw.go",
	"**.pulsar.go",
	"**/node_modules/**",
	"**/openapi.yml",
	"**/openapi.json",
	".gitignore",
	".github/**",
	"**.html",
	"**.css",
	"**.js",
	"**.ts",
	"**.json",
}

// CalculateDiffs calculate the diff from two directories.
func CalculateDiffs(fromDir, toDir string) (Diffs, error) {
	paths, err := readRootFolders(fromDir)
	if err != nil {
		return nil, err
	}
	toPaths, err := readRootFolders(toDir)
	if err != nil {
		return nil, err
	}
	for key, value := range toPaths {
		paths[key] = value
	}

	diffs := make(Diffs)
	for path := range paths {
		from := filepath.Join(fromDir, path)
		if err := os.MkdirAll(from, os.ModePerm); err != nil {
			return nil, err
		}
		to := filepath.Join(toDir, path)
		if err := os.MkdirAll(to, os.ModePerm); err != nil {
			return nil, err
		}

		computedDiff, err := computeFS(
			os.DirFS(from),
			os.DirFS(to),
			diffIgnoreGlobs...,
		)
		if err != nil {
			return nil, err
		}

		diffs[path] = computedDiff
	}
	return subtractBaseDiffs(diffs), nil
}

// SaveDiffs save all migration diffs to the output path.
func SaveDiffs(diffs Diffs, outputPath string) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	for name, diffs := range diffs {
		output, err := os.Create(filepath.Join(outputPath, name+".diff"))
		if err != nil {
			return err
		}
		for _, d := range diffs {
			output.WriteString(fmt.Sprint(d))
			output.WriteString("\n")
		}
		if err := output.Close(); err != nil {
			return err
		}
	}

	return nil
}

// FormatDiffs format all diffs in a single markdown byte array.
func FormatDiffs(diffs Diffs) ([]byte, error) {
	if len(diffs) == 0 {
		return []byte{}, nil
	}
	buffer := &bytes.Buffer{}
	for name, diffs := range diffs {
		if len(diffs) == 0 {
			continue
		}
		buffer.WriteString(fmt.Sprintf("#### **%s diff**\n\n", xstrings.ToUpperFirst(name)))
		buffer.WriteString("```diff\n")
		for _, d := range diffs {
			buffer.WriteString(fmt.Sprint(d))
		}
		buffer.WriteString("```\n\n")
	}
	return buffer.Bytes(), nil
}

// readRootFolders return a map of all root folders from a directory.
func readRootFolders(dir string) (map[string]struct{}, error) {
	paths := make(map[string]struct{})
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			paths[entry.Name()] = struct{}{}
		}
	}
	return paths, nil
}

// subtractBaseDiffs removes chain and module diffs from other diffs.
func subtractBaseDiffs(diffs Diffs) Diffs {
	chainDiff := diffs["chain"]
	moduleDiff := diffs["module"]
	for name, d := range diffs {
		if name != "chain" && name != "module" {
			diffs[name] = subtractUnifieds(d, moduleDiff)
		}
	}
	diffs["module"] = subtractUnifieds(moduleDiff, chainDiff)
	return diffs
}

func subtractUnifieds(a, b []gotextdiff.Unified) []gotextdiff.Unified {
	for i, ad := range a {
		for _, bd := range b {
			if ad.From == bd.From && ad.To == bd.To {
				a[i] = subtract(ad, bd)
			}
		}
	}
	return a
}
