package diff

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hexops/gotextdiff"

	"github.com/ignite/cli/v28/ignite/pkg/diff"
)

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
	".gitignore",
	".github/**",
	"**.html",
	"**.css",
	"**.js",
	"**.ts",
	"**.json",
}

func CalculateDiffs(fromDir, toDir string) (map[string][]gotextdiff.Unified, error) {
	paths := make([]string, 0)
	err := filepath.Walk(fromDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && path != fromDir {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	diffs := make(map[string][]gotextdiff.Unified)
	for _, s := range paths {
		name := filepath.Base(s)
		from := filepath.Join(fromDir, name)
		if err := os.MkdirAll(from, os.ModePerm); err != nil {
			return nil, err
		}
		to := filepath.Join(toDir, name)
		if err := os.MkdirAll(to, os.ModePerm); err != nil {
			return nil, err
		}

		computedDiff, err := diff.ComputeFS(
			os.DirFS(from),
			os.DirFS(to),
			diffIgnoreGlobs...,
		)
		if err != nil {
			return nil, err
		}

		diffs[name] = computedDiff
	}
	return subtractBaseDiffs(diffs), nil
}

// subtractBaseDiffs removes chain and module diffs from other diffs.
func subtractBaseDiffs(diffs map[string][]gotextdiff.Unified) map[string][]gotextdiff.Unified {
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
				a[i] = diff.Subtract(ad, bd)
			}
		}
	}
	return a
}

func SaveDiffs(diffs map[string][]gotextdiff.Unified, outputPath string) error {
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
