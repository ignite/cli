// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var htmlEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

var htmlUnescaper = strings.NewReplacer(
	"&amp;", "&",
	"&lt;", "<",
	"&gt;", ">",
	"&quot;", `"`,
	"&apos;", "'",
)

// The http package's old HTML escaping function.
func oldHTMLEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

var capitalLetters = strings.NewReplacer("a", "A", "b", "B")

// TestReplacer tests the replacer implementations.
func TestReplacer(t *testing.T) {
	type testCase struct {
		r       *strings.Replacer
		in, out string
	}
	var testCases []testCase

	// str converts 0xff to "\xff". This isn't just string(b) since that converts to UTF-8.
	str := func(b byte) string {
		return string([]byte{b})
	}
	var s []string

	// inc maps "\x00"->"\x01", ..., "a"->"b", "b"->"c", ..., "\xff"->"\x00".
	s = nil
	for i := 0; i < 256; i++ {
		s = append(s, str(byte(i)), str(byte(i+1)))
	}
	inc := strings.NewReplacer(s...)

	// Test cases with 1-byte old strings, 1-byte new strings.
	testCases = append(testCases,
		testCase{capitalLetters, "brad", "BrAd"},
		testCase{capitalLetters, strings.Repeat("a", (32<<10)+123), strings.Repeat("A", (32<<10)+123)},
		testCase{capitalLetters, "", ""},

		testCase{inc, "brad", "csbe"},
		testCase{inc, "\x00\xff", "\x01\x00"},
		testCase{inc, "", ""},

		testCase{strings.NewReplacer("a", "1", "a", "2"), "brad", "br1d"},
	)

	// repeat maps "a"->"a", "b"->"bb", "c"->"ccc", ...
	s = nil
	for i := 0; i < 256; i++ {
		n := i + 1 - 'a'
		if n < 1 {
			n = 1
		}
		s = append(s, str(byte(i)), strings.Repeat(str(byte(i)), n))
	}
	repeat := strings.NewReplacer(s...)

	// Test cases with 1-byte old strings, variable length new strings.
	testCases = append(testCases,
		testCase{htmlEscaper, "No changes", "No changes"},
		testCase{htmlEscaper, "I <3 escaping & stuff", "I &lt;3 escaping &amp; stuff"},
		testCase{htmlEscaper, "&&&", "&amp;&amp;&amp;"},
		testCase{htmlEscaper, "", ""},

		testCase{repeat, "brad", "bbrrrrrrrrrrrrrrrrrradddd"},
		testCase{repeat, "abba", "abbbba"},
		testCase{repeat, "", ""},

		testCase{strings.NewReplacer("a", "11", "a", "22"), "brad", "br11d"},
	)

	// The remaining test cases have variable length old strings.

	testCases = append(testCases,
		testCase{htmlUnescaper, "&amp;amp;", "&amp;"},
		testCase{htmlUnescaper, "&lt;b&gt;HTML&apos;s neat&lt;/b&gt;", "<b>HTML's neat</b>"},
		testCase{htmlUnescaper, "", ""},

		testCase{strings.NewReplacer("a", "1", "a", "2", "xxx", "xxx"), "brad", "br1d"},

		testCase{strings.NewReplacer("a", "1", "aa", "2", "aaa", "3"), "aaaa", "1111"},

		testCase{strings.NewReplacer("aaa", "3", "aa", "2", "a", "1"), "aaaa", "31"},
	)

	// gen1 has multiple old strings of variable length. There is no
	// overall non-empty common prefix, but some pairwise common prefixes.
	gen1 := strings.NewReplacer(
		"aaa", "3[aaa]",
		"aa", "2[aa]",
		"a", "1[a]",
		"i", "i",
		"longerst", "most long",
		"longer", "medium",
		"long", "short",
		"xx", "xx",
		"x", "X",
		"X", "Y",
		"Y", "Z",
	)
	testCases = append(testCases,
		testCase{gen1, "fooaaabar", "foo3[aaa]b1[a]r"},
		testCase{gen1, "long, longerst, longer", "short, most long, medium"},
		testCase{gen1, "xxxxx", "xxxxX"},
		testCase{gen1, "XiX", "YiY"},
		testCase{gen1, "", ""},
	)

	// gen2 has multiple old strings with no pairwise common prefix.
	gen2 := strings.NewReplacer(
		"roses", "red",
		"violets", "blue",
		"sugar", "sweet",
	)
	testCases = append(testCases,
		testCase{gen2, "roses are red, violets are blue...", "red are red, blue are blue..."},
		testCase{gen2, "", ""},
	)

	// gen3 has multiple old strings with an overall common prefix.
	gen3 := strings.NewReplacer(
		"abracadabra", "poof",
		"abracadabrakazam", "splat",
		"abraham", "lincoln",
		"abrasion", "scrape",
		"abraham", "isaac",
	)
	testCases = append(testCases,
		testCase{gen3, "abracadabrakazam abraham", "poofkazam lincoln"},
		testCase{gen3, "abrasion abracad", "scrape abracad"},
		testCase{gen3, "abba abram abrasive", "abba abram abrasive"},
		testCase{gen3, "", ""},
	)

	// foo{1,2,3,4} have multiple old strings with an overall common prefix
	// and 1- or 2- byte extensions from the common prefix.
	foo1 := strings.NewReplacer(
		"foo1", "A",
		"foo2", "B",
		"foo3", "C",
	)
	foo2 := strings.NewReplacer(
		"foo1", "A",
		"foo2", "B",
		"foo31", "C",
		"foo32", "D",
	)
	foo3 := strings.NewReplacer(
		"foo11", "A",
		"foo12", "B",
		"foo31", "C",
		"foo32", "D",
	)
	foo4 := strings.NewReplacer(
		"foo12", "B",
		"foo32", "D",
	)
	testCases = append(testCases,
		testCase{foo1, "fofoofoo12foo32oo", "fofooA2C2oo"},
		testCase{foo1, "", ""},

		testCase{foo2, "fofoofoo12foo32oo", "fofooA2Doo"},
		testCase{foo2, "", ""},

		testCase{foo3, "fofoofoo12foo32oo", "fofooBDoo"},
		testCase{foo3, "", ""},

		testCase{foo4, "fofoofoo12foo32oo", "fofooBDoo"},
		testCase{foo4, "", ""},
	)

	// genAll maps "\x00\x01\x02...\xfe\xff" to "[all]", amongst other things.
	allBytes := make([]byte, 256)
	for i := range allBytes {
		allBytes[i] = byte(i)
	}
	allString := string(allBytes)
	genAll := strings.NewReplacer(
		allString, "[all]",
		"\xff", "[ff]",
		"\x00", "[00]",
	)
	testCases = append(testCases,
		testCase{genAll, allString, "[all]"},
		testCase{genAll, "a\xff" + allString + "\x00", "a[ff][all][00]"},
		testCase{genAll, "", ""},
	)

	// Test cases with empty old strings.

	blankToX1 := strings.NewReplacer("", "X")
	blankToX2 := strings.NewReplacer("", "X", "", "")
	blankHighPriority := strings.NewReplacer("", "X", "o", "O")
	blankLowPriority := strings.NewReplacer("o", "O", "", "X")
	blankNoOp1 := strings.NewReplacer("", "")
	blankNoOp2 := strings.NewReplacer("", "", "", "A")
	blankFoo := strings.NewReplacer("", "X", "foobar", "R", "foobaz", "Z")
	testCases = append(testCases,
		testCase{blankToX1, "foo", "XfXoXoX"},
		testCase{blankToX1, "", "X"},

		testCase{blankToX2, "foo", "XfXoXoX"},
		testCase{blankToX2, "", "X"},

		testCase{blankHighPriority, "oo", "XOXOX"},
		testCase{blankHighPriority, "ii", "XiXiX"},
		testCase{blankHighPriority, "oiio", "XOXiXiXOX"},
		testCase{blankHighPriority, "iooi", "XiXOXOXiX"},
		testCase{blankHighPriority, "", "X"},

		testCase{blankLowPriority, "oo", "OOX"},
		testCase{blankLowPriority, "ii", "XiXiX"},
		testCase{blankLowPriority, "oiio", "OXiXiOX"},
		testCase{blankLowPriority, "iooi", "XiOOXiX"},
		testCase{blankLowPriority, "", "X"},

		testCase{blankNoOp1, "foo", "foo"},
		testCase{blankNoOp1, "", ""},

		testCase{blankNoOp2, "foo", "foo"},
		testCase{blankNoOp2, "", ""},

		testCase{blankFoo, "foobarfoobaz", "XRXZX"},
		testCase{blankFoo, "foobar-foobaz", "XRX-XZX"},
		testCase{blankFoo, "", "X"},
	)

	// single string replacer

	abcMatcher := strings.NewReplacer("abc", "[match]")

	testCases = append(testCases,
		testCase{abcMatcher, "", ""},
		testCase{abcMatcher, "ab", "ab"},
		testCase{abcMatcher, "abc", "[match]"},
		testCase{abcMatcher, "abcd", "[match]d"},
		testCase{abcMatcher, "cabcabcdabca", "c[match][match]d[match]a"},
	)

	// Issue 6659 cases (more single string replacer)

	noHello := strings.NewReplacer("Hello", "")
	testCases = append(testCases,
		testCase{noHello, "Hello", ""},
		testCase{noHello, "Hellox", "x"},
		testCase{noHello, "xHello", "x"},
		testCase{noHello, "xHellox", "xx"},
	)

	// No-arg test cases.

	nop := strings.NewReplacer()
	testCases = append(testCases,
		testCase{nop, "abc", "abc"},
		testCase{nop, "", ""},
	)

	// Run the test cases.

	for i, tc := range testCases {
		if s := tc.r.Replace(tc.in); s != tc.out {
			t.Errorf("%d. Replace(%q) = %q, want %q", i, tc.in, s, tc.out)
		}
		var buf bytes.Buffer
		n, err := tc.r.WriteString(&buf, tc.in)
		if err != nil {
			t.Errorf("%d. WriteString: %v", i, err)
			continue
		}
		got := buf.String()
		if got != tc.out {
			t.Errorf("%d. WriteString(%q) wrote %q, want %q", i, tc.in, got, tc.out)
			continue
		}
		if n != len(tc.out) {
			t.Errorf("%d. WriteString(%q) wrote correct string but reported %d bytes; want %d (%q)",
				i, tc.in, n, len(tc.out), tc.out)
		}
	}
}

var algorithmTestCases = []struct {
	r    *strings.Replacer
	want string
}{
	{capitalLetters, "*strings.byteReplacer"},
	{htmlEscaper, "*strings.byteStringReplacer"},
	{strings.NewReplacer("12", "123"), "*strings.singleStringReplacer"},
	{strings.NewReplacer("1", "12"), "*strings.byteStringReplacer"},
	{strings.NewReplacer("", "X"), "*strings.genericReplacer"},
	{strings.NewReplacer("a", "1", "b", "12", "cde", "123"), "*strings.genericReplacer"},
}

//// TestPickAlgorithm tests that strings.NewReplacer picks the correct algorithm.
//func TestPickAlgorithm(t *testing.T) {
//	for i, tc := range algorithmTestCases {
//		got := fmt.Sprintf("%T", tc.r.Replacer())
//		if got != tc.want {
//			t.Errorf("%d. algorithm = %s, want %s", i, got, tc.want)
//		}
//	}
//}

type errWriter struct{}

func (errWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("unwritable")
}

// TestWriteStringError tests that WriteString returns an error
// received from the underlying io.Writer.
func TestWriteStringError(t *testing.T) {
	for i, tc := range algorithmTestCases {
		n, err := tc.r.WriteString(errWriter{}, "abc")
		if n != 0 || err == nil || err.Error() != "unwritable" {
			t.Errorf("%d. WriteStringError = %d, %v, want 0, unwritable", i, n, err)
		}
	}
}

func BenchmarkGenericNoMatch(b *testing.B) {
	str := strings.Repeat("A", 100) + strings.Repeat("B", 100)
	generic := strings.NewReplacer("a", "A", "b", "B", "12", "123") // varying lengths forces generic
	for i := 0; i < b.N; i++ {
		generic.Replace(str)
	}
}

func BenchmarkGenericMatch1(b *testing.B) {
	str := strings.Repeat("a", 100) + strings.Repeat("b", 100)
	generic := strings.NewReplacer("a", "A", "b", "B", "12", "123")
	for i := 0; i < b.N; i++ {
		generic.Replace(str)
	}
}

func BenchmarkGenericMatch2(b *testing.B) {
	str := strings.Repeat("It&apos;s &lt;b&gt;HTML&lt;/b&gt;!", 100)
	for i := 0; i < b.N; i++ {
		htmlUnescaper.Replace(str)
	}
}

func benchmarkSingleString(b *testing.B, pattern, text string) {
	r := strings.NewReplacer(pattern, "[match]")
	b.SetBytes(int64(len(text)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Replace(text)
	}
}

func BenchmarkSingleMaxSkipping(b *testing.B) {
	benchmarkSingleString(b, strings.Repeat("b", 25), strings.Repeat("a", 10000))
}

func BenchmarkSingleLongSuffixFail(b *testing.B) {
	benchmarkSingleString(b, "b"+strings.Repeat("a", 500), strings.Repeat("a", 1002))
}

func BenchmarkSingleMatch(b *testing.B) {
	benchmarkSingleString(b, "abcdef", strings.Repeat("abcdefghijklmno", 1000))
}

func BenchmarkByteByteNoMatch(b *testing.B) {
	str := strings.Repeat("A", 100) + strings.Repeat("B", 100)
	for i := 0; i < b.N; i++ {
		capitalLetters.Replace(str)
	}
}

func BenchmarkByteByteMatch(b *testing.B) {
	str := strings.Repeat("a", 100) + strings.Repeat("b", 100)
	for i := 0; i < b.N; i++ {
		capitalLetters.Replace(str)
	}
}

func BenchmarkByteStringMatch(b *testing.B) {
	str := "<" + strings.Repeat("a", 99) + strings.Repeat("b", 99) + ">"
	for i := 0; i < b.N; i++ {
		htmlEscaper.Replace(str)
	}
}

func BenchmarkHTMLEscapeNew(b *testing.B) {
	str := "I <3 to escape HTML & other text too."
	for i := 0; i < b.N; i++ {
		htmlEscaper.Replace(str)
	}
}

func BenchmarkHTMLEscapeOld(b *testing.B) {
	str := "I <3 to escape HTML & other text too."
	for i := 0; i < b.N; i++ {
		oldHTMLEscape(str)
	}
}

func BenchmarkByteStringReplacerWriteString(b *testing.B) {
	str := strings.Repeat("I <3 to escape HTML & other text too.", 100)
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		htmlEscaper.WriteString(buf, str)
		buf.Reset()
	}
}

func BenchmarkByteReplacerWriteString(b *testing.B) {
	str := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100)
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		capitalLetters.WriteString(buf, str)
		buf.Reset()
	}
}

// BenchmarkByteByteReplaces compares byteByteImpl against multiple Replaces.
func BenchmarkByteByteReplaces(b *testing.B) {
	str := strings.Repeat("a", 100) + strings.Repeat("b", 100)
	for i := 0; i < b.N; i++ {
		strings.ReplaceAll(strings.ReplaceAll(str, "a", "A"), "b", "B")
	}
}

// BenchmarkByteByteMap compares byteByteImpl against Map.
func BenchmarkByteByteMap(b *testing.B) {
	str := strings.Repeat("a", 100) + strings.Repeat("b", 100)
	fn := func(r rune) rune {
		switch r {
		case 'a':
			return 'A'
		case 'b':
			return 'B'
		}
		return r
	}
	for i := 0; i < b.N; i++ {
		strings.Map(fn, str)
	}
}

var mapdata = []struct{ name, data string }{
	{"ASCII", "a b c d e f g h i j k l m n o p q r s t u v w x y z"},
	{"Greek", "α β γ δ ε ζ η θ ι κ λ μ ν ξ ο π ρ ς σ τ υ φ χ ψ ω"},
}

func BenchmarkMap(b *testing.B) {
	mapidentity := func(r rune) rune {
		return r
	}

	b.Run("identity", func(b *testing.B) {
		for _, md := range mapdata {
			b.Run(md.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					strings.Map(mapidentity, md.data)
				}
			})
		}
	})

	mapchange := func(r rune) rune {
		if 'a' <= r && r <= 'z' {
			return r + 'A' - 'a'
		}
		if 'α' <= r && r <= 'ω' {
			return r + 'Α' - 'α'
		}
		return r
	}

	b.Run("change", func(b *testing.B) {
		for _, md := range mapdata {
			b.Run(md.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					strings.Map(mapchange, md.data)
				}
			})
		}
	})
}
