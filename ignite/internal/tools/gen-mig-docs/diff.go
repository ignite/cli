package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

var diffExceptions = []string{
	"*.md",
	"go.sum",
	"*_test.go",
	"*.pb.go",
	"*.pb.gw.go",
}

// Diff returns unified diff between all files in two directories recursively.
func Diff(dir1, dir2 string) ([]gotextdiff.Unified, error) {
	var (
		diffs  []gotextdiff.Unified
		marked map[string]struct{} = make(map[string]struct{})
	)

	dirsLCP := longestCommonPrefix(dir1, dir2)

	// Consider dir1 as reference and walk through all of the files comparing them with files in dir2.
	err := filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		for _, exception := range diffExceptions {
			if match, _ := filepath.Match(exception, info.Name()); match {
				return nil
			}
		}

		relPath, err := filepath.Rel(dir1, path)
		if err != nil {
			return err
		}
		marked[relPath] = struct{}{}

		file1, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		str1 := string(file1)

		file2, err := os.ReadFile(filepath.Join(dir2, relPath))
		if !os.IsNotExist(err) && err != nil {
			return err
		}
		str2 := string(file2)

		edits := myers.ComputeEdits(span.URIFromPath(relPath), str1, str2)
		if len(edits) > 0 {
			diffs = append(diffs, gotextdiff.ToUnified(strings.TrimPrefix(path, dirsLCP), strings.TrimPrefix(filepath.Join(dir2, relPath), dirsLCP), str1, edits))
		}
		return nil
	})
	if err != nil {
		return diffs, err
	}

	// Walk through all of the files in dir2 that were not compared with files in dir1.
	err = filepath.Walk(dir2, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		for _, exception := range diffExceptions {
			if match, _ := filepath.Match(exception, info.Name()); match {
				return nil
			}
		}

		relPath, err := filepath.Rel(dir2, path)
		if err != nil {
			return err
		}
		if _, ok := marked[relPath]; ok {
			return nil
		}

		str1 := ""

		file2, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		str2 := string(file2)

		edits := myers.ComputeEdits(span.URIFromPath(relPath), str1, str2)
		diffs = append(diffs, gotextdiff.ToUnified(strings.TrimPrefix(filepath.Join(dir1, relPath), dirsLCP), strings.TrimPrefix(path, dirsLCP), str1, edits))
		return nil
	})
	return diffs, nil
}

func longestCommonPrefix(strs ...string) string {
	longest := strings.Split(strs[0], string(filepath.Separator))

	cmp := func(a []string) {
		if len(a) < len(longest) {
			longest = longest[:len(a)]
		}
		for i := 0; i < len(longest); i++ {
			if a[i] != longest[i] {
				longest = longest[:i]
				return
			}
		}
	}

	for i := 1; i < len(strs); i++ {
		r := strings.Split(strs[i], string(filepath.Separator))
		cmp(r)
	}
	return filepath.Join(longest...)
}
