package diff

import (
	"sort"

	"github.com/hexops/gotextdiff"
)

// subtract two unified diffs from each other.
func subtract(a, b gotextdiff.Unified) gotextdiff.Unified {
	return gotextdiff.Unified{
		From:  a.From,
		To:    a.To,
		Hunks: subtractHunks(a.Hunks, b.Hunks),
	}
}

func subtractHunks(src, base []*gotextdiff.Hunk) []*gotextdiff.Hunk {
	sortHunks(src)
	sortHunks(base)

	res := make([]*gotextdiff.Hunk, 0, len(src))
	offset := 0
	for i, j := 0, 0; i < len(src) || j < len(base); {
		if i >= len(src) {
			break
		}
		if j >= len(base) {
			res = append(res, src[i])
			i++
			continue
		}

		s := src[i]
		b := base[j]

		switch {
		case beforeHunk(s, b, offset):
			res = append(res, s)
			offset += calculateHunkOffsetChange(s.Lines)
			i++
		case beforeHunk(b, s, -offset):
			j++
		case hunksOverlap(s, b, offset):
			if s.FromLine < b.FromLine {
				res = append(res, s)
				offset += calculateHunkOffsetChange(s.Lines) - calculateHunkOffsetChange(b.Lines)
				i++
			} else {
				offset += calculateHunkOffsetChange(s.Lines) - calculateHunkOffsetChange(b.Lines)
				j++
			}
		default:
			h := subtractHunk(s, b)
			if !isHunkEmpty(h) {
				res = append(res, subtractHunk(s, b))
			}
			offset += calculateHunkOffsetChange(s.Lines) - calculateHunkOffsetChange(b.Lines)
			i++
			j++
		}

	}

	return res
}

func sortHunks(hunks []*gotextdiff.Hunk) {
	sort.Slice(hunks, func(i, j int) bool {
		return hunks[i].FromLine < hunks[j].FromLine
	})
}

// beforeHunk returns true if a comes before b.
func beforeHunk(a, b *gotextdiff.Hunk, offset int) bool {
	return a.ToLine-calculateEndEqualLines(a) < b.FromLine+calculateStartEqualLines(b)+offset
}

func calculateStartEqualLines(h *gotextdiff.Hunk) int {
	lines := 0
	for _, l := range h.Lines {
		if l.Kind == gotextdiff.Equal {
			lines++
		} else {
			break
		}
	}
	return lines
}

func calculateEndEqualLines(h *gotextdiff.Hunk) int {
	lines := 0
	for i := len(h.Lines) - 1; i >= 0; i-- {
		if h.Lines[i].Kind == gotextdiff.Equal {
			lines++
		} else {
			break
		}
	}
	return lines
}

func calculateHunkOffsetChange(lines []gotextdiff.Line) int {
	offset := 0
	for _, l := range lines {
		if l.Kind == gotextdiff.Insert {
			offset++
		} else if l.Kind == gotextdiff.Delete {
			offset--
		}
	}
	return offset
}

func hunksOverlap(a, b *gotextdiff.Hunk, offset int) bool {
	if !isLineInHunk(a.FromLine, b, offset) && isLineInHunk(a.ToLine, b, offset) {
		return true
	}
	if isLineInHunk(a.FromLine, b, offset) && !isLineInHunk(a.ToLine, b, offset) {
		return true
	}
	return false
}

func isLineInHunk(line int, h *gotextdiff.Hunk, offset int) bool {
	return line-calculateStartEqualLines(h) > h.FromLine+offset && line+calculateEndEqualLines(h) < h.ToLine+offset
}

func subtractHunk(a, b *gotextdiff.Hunk) *gotextdiff.Hunk {
	lines := subtractLines(a.Lines, b.Lines)
	return &gotextdiff.Hunk{
		FromLine: a.FromLine,
		ToLine:   a.ToLine + calculateHunkOffsetChange(a.Lines) - calculateHunkOffsetChange(lines),
		Lines:    lines,
	}
}

func subtractLines(a, b []gotextdiff.Line) []gotextdiff.Line {
	res := make([]gotextdiff.Line, 0, len(a))
	for _, la := range a {
		rep := false
		for _, lb := range b {
			if la.Kind != gotextdiff.Equal && la.Kind == lb.Kind && la.Content == lb.Content {
				rep = true
				break
			}
		}

		if !rep {
			res = append(res, la)
		}
	}

	return res
}

func isHunkEmpty(h *gotextdiff.Hunk) bool {
	effectiveLines := 0
	for _, l := range h.Lines {
		if l.Kind != gotextdiff.Equal {
			effectiveLines++
		}
	}
	return effectiveLines == 0
}
