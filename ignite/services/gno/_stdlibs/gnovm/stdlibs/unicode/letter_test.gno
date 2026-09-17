// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unicode_test

import (
	"strings"
	"testing"
	uu "unicode"
)

var upperTest = []rune{
	0x41,
	0xc0,
	0xd8,
	0x100,
	0x139,
	0x14a,
	0x178,
	0x181,
	0x376,
	0x3cf,
	0x13bd,
	0x1f2a,
	0x2102,
	0x2c00,
	0x2c10,
	0x2c20,
	0xa650,
	0xa722,
	0xff3a,
	0x10400,
	0x1d400,
	0x1d7ca,
}

var notupperTest = []rune{
	0x40,
	0x5b,
	0x61,
	0x185,
	0x1b0,
	0x377,
	0x387,
	0x2150,
	0xab7d,
	0xffff,
	0x10000,
}

var letterTest = []rune{
	0x41,
	0x61,
	0xaa,
	0xba,
	0xc8,
	0xdb,
	0xf9,
	0x2ec,
	0x535,
	0x620,
	0x6e6,
	0x93d,
	0xa15,
	0xb99,
	0xdc0,
	0xedd,
	0x1000,
	0x1200,
	0x1312,
	0x1401,
	0x2c00,
	0xa800,
	0xf900,
	0xfa30,
	0xffda,
	0xffdc,
	0x10000,
	0x10300,
	0x10400,
	0x20000,
	0x2f800,
	0x2fa1d,
}

var notletterTest = []rune{
	0x20,
	0x35,
	0x375,
	0x619,
	0x700,
	0x1885,
	0xfffe,
	0x1ffff,
	0x10ffff,
}

// Contains all the special cased Latin-1 chars.
var spaceTest = []rune{
	0x09,
	0x0a,
	0x0b,
	0x0c,
	0x0d,
	0x20,
	0x85,
	0xA0,
	0x2000,
	0x3000,
}

type caseT struct {
	cas     int
	in, out rune
}

var caseTest = []caseT{
	// errors
	{-1, '\n', 0xFFFD},
	{uu.UpperCase, -1, -1},
	{uu.UpperCase, 1 << 30, 1 << 30},

	// ASCII (special-cased so test carefully)
	{uu.UpperCase, '\n', '\n'},
	{uu.UpperCase, 'a', 'A'},
	{uu.UpperCase, 'A', 'A'},
	{uu.UpperCase, '7', '7'},
	{uu.LowerCase, '\n', '\n'},
	{uu.LowerCase, 'a', 'a'},
	{uu.LowerCase, 'A', 'a'},
	{uu.LowerCase, '7', '7'},
	{uu.TitleCase, '\n', '\n'},
	{uu.TitleCase, 'a', 'A'},
	{uu.TitleCase, 'A', 'A'},
	{uu.TitleCase, '7', '7'},

	// Latin-1: easy to read the tests!
	{uu.UpperCase, 0x80, 0x80},
	{uu.UpperCase, 'Å', 'Å'},
	{uu.UpperCase, 'å', 'Å'},
	{uu.LowerCase, 0x80, 0x80},
	{uu.LowerCase, 'Å', 'å'},
	{uu.LowerCase, 'å', 'å'},
	{uu.TitleCase, 0x80, 0x80},
	{uu.TitleCase, 'Å', 'Å'},
	{uu.TitleCase, 'å', 'Å'},

	// 0131;LATIN SMALL LETTER DOTLESS I;Ll;0;L;;;;;N;;;0049;;0049
	{uu.UpperCase, 0x0131, 'I'},
	{uu.LowerCase, 0x0131, 0x0131},
	{uu.TitleCase, 0x0131, 'I'},

	// 0133;LATIN SMALL LIGATURE IJ;Ll;0;L;<compat> 0069 006A;;;;N;LATIN SMALL LETTER I J;;0132;;0132
	{uu.UpperCase, 0x0133, 0x0132},
	{uu.LowerCase, 0x0133, 0x0133},
	{uu.TitleCase, 0x0133, 0x0132},

	// 212A;KELVIN SIGN;Lu;0;L;004B;;;;N;DEGREES KELVIN;;;006B;
	{uu.UpperCase, 0x212A, 0x212A},
	{uu.LowerCase, 0x212A, 'k'},
	{uu.TitleCase, 0x212A, 0x212A},

	// From an UpperLower sequence
	// A640;CYRILLIC CAPITAL LETTER ZEMLYA;Lu;0;L;;;;;N;;;;A641;
	{uu.UpperCase, 0xA640, 0xA640},
	{uu.LowerCase, 0xA640, 0xA641},
	{uu.TitleCase, 0xA640, 0xA640},
	// A641;CYRILLIC SMALL LETTER ZEMLYA;Ll;0;L;;;;;N;;;A640;;A640
	{uu.UpperCase, 0xA641, 0xA640},
	{uu.LowerCase, 0xA641, 0xA641},
	{uu.TitleCase, 0xA641, 0xA640},
	// A64E;CYRILLIC CAPITAL LETTER NEUTRAL YER;Lu;0;L;;;;;N;;;;A64F;
	{uu.UpperCase, 0xA64E, 0xA64E},
	{uu.LowerCase, 0xA64E, 0xA64F},
	{uu.TitleCase, 0xA64E, 0xA64E},
	// A65F;CYRILLIC SMALL LETTER YN;Ll;0;L;;;;;N;;;A65E;;A65E
	{uu.UpperCase, 0xA65F, 0xA65E},
	{uu.LowerCase, 0xA65F, 0xA65F},
	{uu.TitleCase, 0xA65F, 0xA65E},

	// From another UpperLower sequence
	// 0139;LATIN CAPITAL LETTER L WITH ACUTE;Lu;0;L;004C 0301;;;;N;LATIN CAPITAL LETTER L ACUTE;;;013A;
	{uu.UpperCase, 0x0139, 0x0139},
	{uu.LowerCase, 0x0139, 0x013A},
	{uu.TitleCase, 0x0139, 0x0139},
	// 013F;LATIN CAPITAL LETTER L WITH MIDDLE DOT;Lu;0;L;<compat> 004C 00B7;;;;N;;;;0140;
	{uu.UpperCase, 0x013f, 0x013f},
	{uu.LowerCase, 0x013f, 0x0140},
	{uu.TitleCase, 0x013f, 0x013f},
	// 0148;LATIN SMALL LETTER N WITH CARON;Ll;0;L;006E 030C;;;;N;LATIN SMALL LETTER N HACEK;;0147;;0147
	{uu.UpperCase, 0x0148, 0x0147},
	{uu.LowerCase, 0x0148, 0x0148},
	{uu.TitleCase, 0x0148, 0x0147},

	// Lowercase lower than uppercase.
	// AB78;CHEROKEE SMALL LETTER GE;Ll;0;L;;;;;N;;;13A8;;13A8
	{uu.UpperCase, 0xab78, 0x13a8},
	{uu.LowerCase, 0xab78, 0xab78},
	{uu.TitleCase, 0xab78, 0x13a8},
	{uu.UpperCase, 0x13a8, 0x13a8},
	{uu.LowerCase, 0x13a8, 0xab78},
	{uu.TitleCase, 0x13a8, 0x13a8},

	// Last block in the 5.1.0 table
	// 10400;DESERET CAPITAL LETTER LONG I;Lu;0;L;;;;;N;;;;10428;
	{uu.UpperCase, 0x10400, 0x10400},
	{uu.LowerCase, 0x10400, 0x10428},
	{uu.TitleCase, 0x10400, 0x10400},
	// 10427;DESERET CAPITAL LETTER EW;Lu;0;L;;;;;N;;;;1044F;
	{uu.UpperCase, 0x10427, 0x10427},
	{uu.LowerCase, 0x10427, 0x1044F},
	{uu.TitleCase, 0x10427, 0x10427},
	// 10428;DESERET SMALL LETTER LONG I;Ll;0;L;;;;;N;;;10400;;10400
	{uu.UpperCase, 0x10428, 0x10400},
	{uu.LowerCase, 0x10428, 0x10428},
	{uu.TitleCase, 0x10428, 0x10400},
	// 1044F;DESERET SMALL LETTER EW;Ll;0;L;;;;;N;;;10427;;10427
	{uu.UpperCase, 0x1044F, 0x10427},
	{uu.LowerCase, 0x1044F, 0x1044F},
	{uu.TitleCase, 0x1044F, 0x10427},

	// First one not in the 5.1.0 table
	// 10450;SHAVIAN LETTER PEEP;Lo;0;L;;;;;N;;;;;
	{uu.UpperCase, 0x10450, 0x10450},
	{uu.LowerCase, 0x10450, 0x10450},
	{uu.TitleCase, 0x10450, 0x10450},

	// Non-letters with case.
	{uu.LowerCase, 0x2161, 0x2171},
	{uu.UpperCase, 0x0345, 0x0399},
}

func TestIsLetter(t *testing.T) {
	for _, r := range upperTest {
		if !uu.IsLetter(r) {
			t.Errorf("IsLetter(U+%04X) = false, want true", r)
		}
	}
	for _, r := range letterTest {
		if !uu.IsLetter(r) {
			t.Errorf("IsLetter(U+%04X) = false, want true", r)
		}
	}
	for _, r := range notletterTest {
		if uu.IsLetter(r) {
			t.Errorf("IsLetter(U+%04X) = true, want false", r)
		}
	}
}

func TestIsUpper(t *testing.T) {
	for _, r := range upperTest {
		if !uu.IsUpper(r) {
			t.Errorf("IsUpper(U+%04X) = false, want true", r)
		}
	}
	for _, r := range notupperTest {
		if uu.IsUpper(r) {
			t.Errorf("IsUpper(U+%04X) = true, want false", r)
		}
	}
	for _, r := range notletterTest {
		if uu.IsUpper(r) {
			t.Errorf("IsUpper(U+%04X) = true, want false", r)
		}
	}
}

func caseString(c int) string {
	switch c {
	case uu.UpperCase:
		return "uu.UpperCase"
	case uu.LowerCase:
		return "uu.LowerCase"
	case uu.TitleCase:
		return "uu.TitleCase"
	}
	return "ErrorCase"
}

func TestTo(t *testing.T) {
	for _, c := range caseTest {
		r := uu.To(c.cas, c.in)
		if c.out != r {
			t.Errorf("To(U+%04X, %s) = U+%04X want U+%04X", c.in, caseString(c.cas), r, c.out)
		}
	}
}

func TestToUpperCase(t *testing.T) {
	for _, c := range caseTest {
		if c.cas != uu.UpperCase {
			continue
		}
		r := uu.ToUpper(c.in)
		if c.out != r {
			t.Errorf("ToUpper(U+%04X) = U+%04X want U+%04X", c.in, r, c.out)
		}
	}
}

func TestToLowerCase(t *testing.T) {
	for _, c := range caseTest {
		if c.cas != uu.LowerCase {
			continue
		}
		r := uu.ToLower(c.in)
		if c.out != r {
			t.Errorf("ToLower(U+%04X) = U+%04X want U+%04X", c.in, r, c.out)
		}
	}
}

func TestToTitleCase(t *testing.T) {
	for _, c := range caseTest {
		if c.cas != uu.TitleCase {
			continue
		}
		r := uu.ToTitle(c.in)
		if c.out != r {
			t.Errorf("ToTitle(U+%04X) = U+%04X want U+%04X", c.in, r, c.out)
		}
	}
}

func TestIsSpace(t *testing.T) {
	for _, c := range spaceTest {
		if !uu.IsSpace(c) {
			t.Errorf("IsSpace(U+%04X) = false; want true", c)
		}
	}
	for _, c := range letterTest {
		if uu.IsSpace(c) {
			t.Errorf("IsSpace(U+%04X) = true; want false", c)
		}
	}
}

// Check that the optimizations for IsLetter etc. agree with the tables.
// We only need to check the Latin-1 range.
func TestLetterOptimizations(t *testing.T) {
	for i := rune(0); i <= uu.MaxLatin1; i++ {
		if uu.Is(uu.Letter, i) != uu.IsLetter(i) {
			t.Errorf("IsLetter(U+%04X) disagrees with Is(Letter)", i)
		}
		if uu.Is(uu.Upper, i) != uu.IsUpper(i) {
			t.Errorf("IsUpper(U+%04X) disagrees with Is(Upper)", i)
		}
		if uu.Is(uu.Lower, i) != uu.IsLower(i) {
			t.Errorf("IsLower(U+%04X) disagrees with Is(Lower)", i)
		}
		if uu.Is(uu.Title, i) != uu.IsTitle(i) {
			t.Errorf("IsTitle(U+%04X) disagrees with Is(Title)", i)
		}
		if uu.Is(uu.White_Space, i) != uu.IsSpace(i) {
			t.Errorf("IsSpace(U+%04X) disagrees with Is(White_Space)", i)
		}
		if uu.To(uu.UpperCase, i) != uu.ToUpper(i) {
			t.Errorf("ToUpper(U+%04X) disagrees with To(Upper)", i)
		}
		if uu.To(uu.LowerCase, i) != uu.ToLower(i) {
			t.Errorf("ToLower(U+%04X) disagrees with To(Lower)", i)
		}
		if uu.To(uu.TitleCase, i) != uu.ToTitle(i) {
			t.Errorf("ToTitle(U+%04X) disagrees with To(Title)", i)
		}
	}
}

func TestTurkishCase(t *testing.T) {
	lower := []rune("abcçdefgğhıijklmnoöprsştuüvyz")
	upper := []rune("ABCÇDEFGĞHIİJKLMNOÖPRSŞTUÜVYZ")
	for i, l := range lower {
		u := upper[i]
		if uu.TurkishCase.ToLower(l) != l {
			t.Errorf("lower(U+%04X) is U+%04X not U+%04X", l, uu.TurkishCase.ToLower(l), l)
		}
		if uu.TurkishCase.ToUpper(u) != u {
			t.Errorf("upper(U+%04X) is U+%04X not U+%04X", u, uu.TurkishCase.ToUpper(u), u)
		}
		if uu.TurkishCase.ToUpper(l) != u {
			t.Errorf("upper(U+%04X) is U+%04X not U+%04X", l, uu.TurkishCase.ToUpper(l), u)
		}
		if uu.TurkishCase.ToLower(u) != l {
			t.Errorf("lower(U+%04X) is U+%04X not U+%04X", u, uu.TurkishCase.ToLower(l), l)
		}
		if uu.TurkishCase.ToTitle(u) != u {
			t.Errorf("title(U+%04X) is U+%04X not U+%04X", u, uu.TurkishCase.ToTitle(u), u)
		}
		if uu.TurkishCase.ToTitle(l) != u {
			t.Errorf("title(U+%04X) is U+%04X not U+%04X", l, uu.TurkishCase.ToTitle(l), u)
		}
	}
}

var simpleFoldTests = []string{
	// SimpleFold(x) returns the next equivalent rune > x or wraps
	// around to smaller values.

	// Easy cases.
	"Aa",
	"δΔ",

	// ASCII special cases.
	"KkK",
	"Ssſ",

	// Non-ASCII special cases.
	"ρϱΡ",
	"ͅΙιι",

	// Extra special cases: has lower/upper but no case fold.
	"İ",
	"ı",

	// Upper comes before lower (Cherokee).
	"\u13b0\uab80",
}

func TestSimpleFold(t *testing.T) {
	for _, tt := range simpleFoldTests {
		cycle := []rune(tt)
		r := cycle[len(cycle)-1]
		for _, out := range cycle {
			if r := uu.SimpleFold(r); r != out {
				t.Errorf("SimpleFold(%#U) = %#U, want %#U", r, r, out)
			}
			r = out
		}
	}

	if r := uu.SimpleFold(-42); r != -42 {
		t.Errorf("SimpleFold(-42) = %v, want -42", r)
	}
}

/* REMOVED FOR GNO
// Running 'go test -calibrate' runs the calibration to find a plausible
// cutoff point for linear search of a range list vs. binary search.
// We create a fake table and then time how long it takes to do a
// sequence of searches within that table, for all possible inputs
// relative to the ranges (something before all, in each, between each, after all).
// This assumes that all possible runes are equally likely.
// In practice most runes are ASCII so this is a conservative estimate
// of an effective cutoff value. In practice we could probably set it higher
// than what this function recommends.

var calibrate = flag.Bool("calibrate", false, "compute crossover for linear vs. binary search")

func TestCalibrate(t *testing.T) {
	if !*calibrate {
		return
	}

	if runtime.GOARCH == "amd64" {
		fmt.Printf("warning: running calibration on %s\n", runtime.GOARCH)
	}

	// Find the point where binary search wins by more than 10%.
	// The 10% bias gives linear search an edge when they're close,
	// because on predominantly ASCII inputs linear search is even
	// better than our benchmarks measure.
	n := sort.Search(64, func(n int) bool {
		tab := fakeTable(n)
		blinear := func(b *testing.B) {
			tab := tab
			max := n*5 + 20
			for i := 0; i < b.N; i++ {
				for j := 0; j <= max; j++ {
					linear(tab, uint16(j))
				}
			}
		}
		bbinary := func(b *testing.B) {
			tab := tab
			max := n*5 + 20
			for i := 0; i < b.N; i++ {
				for j := 0; j <= max; j++ {
					binary(tab, uint16(j))
				}
			}
		}
		bmlinear := testing.Benchmark(blinear)
		bmbinary := testing.Benchmark(bbinary)
		fmt.Printf("n=%d: linear=%d binary=%d\n", n, bmlinear.NsPerOp(), bmbinary.NsPerOp())
		return bmlinear.NsPerOp()*100 > bmbinary.NsPerOp()*110
	})
	fmt.Printf("calibration: linear cutoff = %d\n", n)
}
*/

func fakeTable(n int) []uu.Range16 {
	var r16 []uu.Range16
	for i := 0; i < n; i++ {
		r16 = append(r16, uu.Range16{uint16(i*5 + 10), uint16(i*5 + 12), 1})
	}
	return r16
}

func linear(ranges []uu.Range16, r uint16) bool {
	for i := range ranges {
		range_ := &ranges[i]
		if r < range_.Lo {
			return false
		}
		if r <= range_.Hi {
			return (r-range_.Lo)%range_.Stride == 0
		}
	}
	return false
}

func binary(ranges []uu.Range16, r uint16) bool {
	// binary search over ranges
	lo := 0
	hi := len(ranges)
	for lo < hi {
		m := lo + (hi-lo)/2
		range_ := &ranges[m]
		if range_.Lo <= r && r <= range_.Hi {
			return (r-range_.Lo)%range_.Stride == 0
		}
		if r < range_.Lo {
			hi = m
		} else {
			lo = m + 1
		}
	}
	return false
}

func TestLatinOffset(t *testing.T) {
	maps := []map[string]*uu.RangeTable{
		uu.Categories,
		uu.FoldCategory,
		uu.FoldScript,
		uu.Properties,
		uu.Scripts,
	}
	for _, m := range maps {
		for name, tab := range m {
			i := 0
			for i < len(tab.R16) && tab.R16[i].Hi <= uu.MaxLatin1 {
				i++
			}
			if tab.LatinOffset != i {
				t.Errorf("%s: LatinOffset=%d, want %d", name, tab.LatinOffset, i)
			}
		}
	}
}

func TestSpecialCaseNoMapping(t *testing.T) {
	// Issue 25636
	// no change for rune 'A', zero delta, under upper/lower/title case change.
	noChangeForCapitalA := uu.CaseRange{'A', 'A', [uu.MaxCase]rune{0, 0, 0}}
	got := strings.ToLowerSpecial(uu.SpecialCase([]uu.CaseRange{noChangeForCapitalA}), "ABC")
	want := "Abc"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestNegativeRune(t *testing.T) {
	// Issue 43254
	// These tests cover negative rune handling by testing values which,
	// when cast to uint8 or uint16, look like a particular valid rune.
	// This package has Latin-1-specific optimizations, so we test all of
	// Latin-1 and representative non-Latin-1 values in the character
	// categories covered by IsGraphic, etc.
	nonLatin1 := []uint32{
		// Lu: LATIN CAPITAL LETTER A WITH MACRON
		0x0100,
		// Ll: LATIN SMALL LETTER A WITH MACRON
		0x0101,
		// Lt: LATIN CAPITAL LETTER D WITH SMALL LETTER Z WITH CARON
		0x01C5,
		// M: COMBINING GRAVE ACCENT
		0x0300,
		// Nd: ARABIC-INDIC DIGIT ZERO
		0x0660,
		// P: GREEK QUESTION MARK
		0x037E,
		// S: MODIFIER LETTER LEFT ARROWHEAD
		0x02C2,
		// Z: OGHAM SPACE MARK
		0x1680,
	}
	for i := 0; i < uu.MaxLatin1+len(nonLatin1); i++ {
		base := uint32(i)
		if i >= uu.MaxLatin1 {
			base = nonLatin1[i-uu.MaxLatin1]
		}

		// Note r is negative, but uint8(r) == uint8(base) and
		// uint16(r) == uint16(base).
		r := rune(base - 1<<31)
		if uu.Is(uu.Letter, r) {
			t.Errorf("Is(Letter, 0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsControl(r) {
			t.Errorf("IsControl(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsDigit(r) {
			t.Errorf("IsDigit(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsGraphic(r) {
			t.Errorf("IsGraphic(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsLetter(r) {
			t.Errorf("IsLetter(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsLower(r) {
			t.Errorf("IsLower(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsMark(r) {
			t.Errorf("IsMark(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsNumber(r) {
			t.Errorf("IsNumber(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsPrint(r) {
			t.Errorf("IsPrint(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsPunct(r) {
			t.Errorf("IsPunct(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsSpace(r) {
			t.Errorf("IsSpace(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsSymbol(r) {
			t.Errorf("IsSymbol(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsTitle(r) {
			t.Errorf("IsTitle(0x%x - 1<<31) = true, want false", base)
		}
		if uu.IsUpper(r) {
			t.Errorf("IsUpper(0x%x - 1<<31) = true, want false", base)
		}
	}
}
