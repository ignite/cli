package diff

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

type Diff struct {
	dir1, dir2     string
	lcp            string
	files1, files2 map[string]string
	edits          []gotextdiff.TextEdit
}

func ComputeDiff(dir1, dir2 string) (*Diff, error) {
	marked := make(map[string]struct{})
	diff := &Diff{
		dir1:   dir1,
		dir2:   dir2,
		lcp:    longestCommonPrefix(dir1, dir2),
		files1: make(map[string]string),
		files2: make(map[string]string),
		edits:  make([]gotextdiff.TextEdit, 0),
	}

	// Consider dir1 as reference and walk through all of the files comparing them with files in dir2.
	err := filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isException(path) {
			return nil
		}

		relPath, err := filepath.Rel(dir1, path)
		if err != nil {
			return err
		}
		marked[relPath] = struct{}{}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		diff.files1[relPath] = string(b)

		b, err = os.ReadFile(filepath.Join(dir2, relPath))
		// If the file does not exist in dir2, we consider it as an empty file.
		if !os.IsNotExist(err) && err != nil {
			return err
		}
		diff.files2[relPath] = string(b)

		edits := myers.ComputeEdits(span.URIFromPath(path), diff.files1[relPath], diff.files2[relPath])
		conv := span.NewContentConverter(span.URIFromPath(path).Filename(), []byte(diff.files1[relPath]))
		err = ensureEditsHavePositionAndOffset(edits, conv)
		if err != nil {
			return err
		}
		diff.edits = append(diff.edits, edits...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Walk through all of the files in dir2 that were not compared with files in dir1.
	err = filepath.Walk(dir2, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isException(path) {
			return nil
		}

		relPath, err := filepath.Rel(dir2, path)
		if err != nil {
			return err
		}
		if _, ok := marked[relPath]; ok {
			return nil
		}

		diff.files1[relPath] = ""

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		diff.files2[relPath] = string(b)

		edits := myers.ComputeEdits(span.URIFromPath(filepath.Join(dir1, relPath)), diff.files1[relPath], diff.files2[relPath])
		diff.edits = append(diff.edits, edits...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	diff.makeSureEditsAreSorted()

	return diff, err
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
	return "/" + filepath.Join(longest...)
}

func ensureEditsHavePositionAndOffset(edits []gotextdiff.TextEdit, conv span.Converter) error {
	var err error
	for i, e := range edits {
		edits[i].Span, err = e.Span.WithAll(conv)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Diff) makeSureEditsAreSorted() {
	gotextdiff.SortTextEdits(d.edits)
}

// Subtract removes all the common changes of base and d from d.
// Note that this function only works for cases where all the changes of base is included in d.
func (d *Diff) Subtract(base *Diff) error {
	baseFiles := base.groupEditsByFile()

	newEdits := make([]gotextdiff.TextEdit, 0, len(d.edits))
	for f, edits := range d.groupEditsByFile() {
		if _, ok := baseFiles[f]; !ok {
			newEdits = append(newEdits, edits...)
			continue
		}

		// Because both base and d are sorted, we use a merge sort like algorithm to subtract base from d.
		var i, j, traverseOffset int
		for i < len(edits) && j < len(baseFiles[f]) {
			e := edits[i]
			b := baseFiles[f][j]

			if e.Span.End().Offset() < b.Span.Start().Offset()+traverseOffset {
				newEdits = append(newEdits, e)
				i++
				traverseOffset += calculateOffsetChange(e)
				continue
			}

			// Ideally, this condition should never be met as we are assuming that all the changes of base is included in d.
			if e.Span.Start().Offset() > b.Span.End().Offset()+traverseOffset {
				j++
				continue
			}

			if spansHaveConflict(e.Span, b.Span, traverseOffset) {
				// If there is a conflict, we add the change of d and move to the next change.
				newEdits = append(newEdits, e)

				if e.Span.Start().Offset() < b.Span.Start().Offset()+traverseOffset {
					i++
					traverseOffset += len(e.NewText) - (b.Span.Start().Offset() - e.Span.Start().Offset())
				} else {
					j++
					traverseOffset += len(e.NewText) - (e.Span.Start().Offset() - b.Span.End().Offset())
				}
			}

			// Finally the two cases where either e is in the middle of b or b is in the middle of e.
			if e.Span.Start().Offset() >= b.Span.Start().Offset()+traverseOffset {
				aconv := span.NewContentConverter(e.Span.URI().Filename(), []byte(d.files1[f]))
				editParts, err := subtractEdits(e, b, aconv)
				if err != nil {
					return err
				}
				newEdits = append(newEdits, editParts...)
				traverseOffset += calculateOffsetChange(e) - calculateOffsetChange(b)
			}

			i++
			j++
		}
	}

	d.edits = newEdits
	return nil
}

func (d *Diff) groupEditsByFile() map[string][]gotextdiff.TextEdit {
	d.makeSureEditsAreSorted()

	fileEdits := make(map[string][]gotextdiff.TextEdit)
	for _, e := range d.edits {
		path, err := filepath.Rel(d.dir1, e.Span.URI().Filename())
		if err != nil {
			panic(err)
		}
		fileEdits[path] = append(fileEdits[path], e)
	}
	return fileEdits
}

func calculateOffsetChange(e gotextdiff.TextEdit) int {
	return len(e.NewText) - (e.Span.End().Offset() - e.Span.Start().Offset())
}

func areEditsEqual(a, b gotextdiff.TextEdit, offset int) bool {
	if a.Span.Start().Offset() != b.Span.Start().Offset()+offset {
		return false
	}
	if a.Span.End().Offset() != b.Span.End().Offset()+offset {
		return false
	}
	if a.NewText != b.NewText {
		return false
	}
	return true
}

func spansHaveConflict(a, b span.Span, offset int) bool {
	if isPointInSpan(a.Start(), b) && !isPointInSpan(a.End(), b) {
		return true
	}
	if isPointInSpan(a.End(), b) && !isPointInSpan(a.Start(), b) {
		return true
	}
	return false
}

func isPointInSpan(p span.Point, s span.Span) bool {
	return p.Offset() >= s.Start().Offset() && p.Offset() <= s.End().Offset()
}

func subtractEdits(a, b gotextdiff.TextEdit, aconv span.Converter) ([]gotextdiff.TextEdit, error) {
	edits := myers.ComputeEdits(a.Span.URI(), b.NewText, a.NewText)
	bconv := span.NewContentConverter(b.Span.URI().Filename(), []byte(b.NewText))
	for i, e := range edits {
		s, err := e.Span.WithOffset(bconv)
		if err != nil {
			return nil, err
		}
		edits[i].Span = moveSpan(s, a.Span.End().Offset())
	}

	err := ensureEditsHavePositionAndOffset(edits, aconv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure edits have position and offset")
	}

	return edits, nil
}

func moveSpan(s span.Span, offset int) span.Span {
	return span.New(s.URI(), span.NewPoint(0, 0, s.Start().Offset()+offset), span.NewPoint(0, 0, s.End().Offset()+offset))
}

func (d *Diff) ToUnified() []gotextdiff.Unified {
	unified := make([]gotextdiff.Unified, 0, len(d.edits))
	fileEdits := d.groupEditsByFile()
	for path, edits := range fileEdits {
		from, err := filepath.Rel(d.lcp, filepath.Join(d.dir1, path))
		if err != nil {
			panic(err)
		}
		to, err := filepath.Rel(d.lcp, filepath.Join(d.dir2, path))
		if err != nil {
			panic(err)
		}
		unified = append(unified, gotextdiff.ToUnified(from, to, d.files1[path], edits))
	}
	return unified
}
