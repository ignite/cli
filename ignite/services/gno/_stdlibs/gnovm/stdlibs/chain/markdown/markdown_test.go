package markdown

import (
	"strings"
	"testing"
)

func TestStripBidiAndZeroWidth(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain", "plain"},
		{"a\u200Bb", "ab"},               // ZWSP
		{"a\u202Eevil\u202Cb", "aevilb"}, // RLO + PDF
		{"\uFEFFstart", "start"},         // BOM
		{"clean\nbreak", "clean\nbreak"},
	}
	for _, c := range cases {
		if got := StripBidiAndZeroWidth(c.in); got != c.want {
			t.Errorf("StripBidiAndZeroWidth(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeBreaks(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a\nb", "a\nb"},
		{"a\rb", "a\nb"},
		{"a\r\nb", "a\nb"},
		{"a\r\n\rb", "a\n\nb"},
	}
	for _, c := range cases {
		if got := NormalizeBreaks(c.in); got != c.want {
			t.Errorf("NormalizeBreaks(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEscapeInline(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain", "plain"},
		{"a*b", `a\*b`},
		{"a|b", "a|b"}, // | NOT in inline set
		{"a=b", "a=b"}, // = NOT in inline set
		{"a\x00b", "a\xef\xbf\xbdb"},
		{"\\*", `\\\*`},
	}
	for _, c := range cases {
		if got := EscapeInline(c.in); got != c.want {
			t.Errorf("EscapeInline(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEscapeTitle(t *testing.T) {
	cases := []struct{ in, want string }{
		{`he said "hi"`, `he said \"hi\"`},
		{`it's a (test)`, `it\'s a \(test\)`},
	}
	for _, c := range cases {
		if got := EscapeTitle(c.in); got != c.want {
			t.Errorf("EscapeTitle(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPercentEncodeURL(t *testing.T) {
	cases := []struct{ in, want string }{
		{"http://a.com/x", "http://a.com/x"},
		{"a b", "a%20b"},
		{"a<b", "a%3Cb"},
		{"a%20b", "a%20b"},   // already encoded, preserved
		{"a%zzb", "a%25zzb"}, // bare %
		{"&#x6a;avascript:alert", "%26#x6a;avascript:alert"},
		{"&#106;avascript:alert", "%26#106;avascript:alert"},
		// Leading zeros make goldmark parse the reference as octal
		// (strconv.ParseUint with base 0), so 0000152 is also 'j'.
		{"&#0000152;avascript:alert", "%26#0000152;avascript:alert"},
		{"javascript&colon;alert", "javascript%26colon;alert"},
		{"&bogus;", "%26bogus;"},
		{"&#12345678;", "%26#12345678;"},
		{"&name", "&name"},
		{"&#;", "&#;"},
		{"?x=1&y=2", "?x=1&y=2"},
		{"/r/boards$help&func=DeleteThread&boardID=2", "/r/boards$help&func=DeleteThread&boardID=2"},
		// UnescapePunctuations runs before the reference resolvers, so a
		// `\` on the `#` or the `;` must not hide the shape from us.
		{`&\#x6a;avascript:alert`, `%26%5C#x6a;avascript:alert`},
		{`&#x6a\;avascript:alert`, `%26#x6a%5C;avascript:alert`},
		{`&\#x6a\;avascript:alert`, `%26%5C#x6a%5C;avascript:alert`},
		{`&\#106;avascript:alert`, `%26%5C#106;avascript:alert`},
		{`javascript&colon\;alert`, `javascript%26colon%5C;alert`},
		{`&\name;`, `&%5Cname;`}, // `\n` is not punctuation — no reference to hide
		// A renderer runs its resolvers over each other's OUTPUT, so a
		// numeric reference can manufacture the `&` its entity-name pass
		// then consumes: `&#38;colon;` becomes `&colon;` becomes `:`.
		// Encoding the `&` of the numeric reference is what stops the
		// chain — there is no second `&` to encode.
		{"javascript&#38;colon;alert", "javascript%26#38;colon;alert"},
		{"javascript&#x26;colon;alert", "javascript%26#x26;colon;alert"},
	}
	for _, c := range cases {
		if got := PercentEncodeURL(c.in); got != c.want {
			t.Errorf("PercentEncodeURL(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestBlockEncodesEntityReferencesInInlineLinkURLs(t *testing.T) {
	for _, sanitize := range []struct {
		name string
		fn   func(string) string
	}{
		{"Block", EscapeBlockHazards},
		{"BlockRich", EscapeBlockHazardsRich},
	} {
		t.Run(sanitize.name, func(t *testing.T) {
			for _, c := range []struct{ in, want string }{
				{"[x](&#x6a;avascript:alert)\n", "[x](%26#x6a;avascript:alert)\n"},
				{"[x](&#106;avascript:alert)\n", "[x](%26#106;avascript:alert)\n"},
				{"[x](&#0000152;avascript:alert)\n", "[x](%26#0000152;avascript:alert)\n"}, // octal

				{"[x](javascript&colon;alert)\n", "[x](javascript%26colon;alert)\n"},
				{"[x](?x=1&y=2)\n", "[x](?x=1&y=2)\n"},

				// A `\` on the `#` or the `;` survives the sanitizer but
				// not the renderer: UnescapePunctuations runs before the
				// reference resolvers, so these reach a resolver as
				// `&#x6a;` / `&colon;` unless the `&` is encoded here.
				// The backslash itself is left alone — the span keeps its
				// escapes so `\)` still holds the destination together.
				{"[x](&\\#x6a;avascript:alert)\n", "[x](%26\\#x6a;avascript:alert)\n"},
				{"[x](&#x6a\\;avascript:alert)\n", "[x](%26#x6a\\;avascript:alert)\n"},
				{"[x](&\\#x6a\\;avascript:alert)\n", "[x](%26\\#x6a\\;avascript:alert)\n"},
				{"[x](javascript&colon\\;alert)\n", "[x](javascript%26colon\\;alert)\n"},
				{"[x](&\\#x64;ata:text/html,x)\n", "[x](%26\\#x64;ata:text/html,x)\n"},

				// Chained resolution: `&#38;colon;` reaches the entity-name
				// pass as `&colon;`, which resolves to `:`. See the
				// TestPercentEncodeURL cases.
				{"[x](javascript&#38;colon;alert)\n", "[x](javascript%26#38;colon;alert)\n"},
			} {
				if got := sanitize.fn(c.in); got != c.want {
					t.Errorf("%s(%q) = %q, want %q", sanitize.name, c.in, got, c.want)
				}
			}
		})
	}
}

func TestMatchCharsetN(t *testing.T) {
	// Build bitmaps for [a-z][a-z0-9]* with bounds [1, 64].
	var firstLo, firstHi uint64
	for c := byte('a'); c <= 'z'; c++ {
		if c < 64 {
			firstLo |= 1 << c
		} else {
			firstHi |= 1 << (c - 64)
		}
	}
	restLo, restHi := firstLo, firstHi
	for c := byte('0'); c <= '9'; c++ {
		restLo |= 1 << c
	}
	cases := []struct {
		s    string
		want bool
	}{
		{"abc", true},
		{"abc123", true},
		{"1abc", false}, // first must be letter
		{"", false},     // below minLen
		{"a_b", false},  // _ not in set
	}
	for _, c := range cases {
		if got := MatchCharsetN(c.s, firstLo, firstHi, restLo, restHi, 1, 64); got != c.want {
			t.Errorf("MatchCharsetN(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}

func TestCodeFence(t *testing.T) {
	cases := []struct {
		content string
		min     int
		want    string
	}{
		{"", 3, "```"},
		{"```", 3, "````"},             // 3 backticks → 4-fence
		{"a `b` c", 3, "```"},          // longest run is 1, min wins
		{"a ``` b ```` c", 3, "`````"}, // longest run 4, +1 = 5
		{"", 0, "`"},                   // min < 1 clamped to 1
	}
	for _, c := range cases {
		if got := CodeFence(c.content, c.min); got != c.want {
			t.Errorf("CodeFence(%q, %d) = %q, want %q", c.content, c.min, got, c.want)
		}
	}
}

func TestEscapeBlockHazards(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "hello world\n", "hello world\n"},
		{"atx-heading", "# evil heading\n", "\\# evil heading\n"},
		{"blockquote", "> caller\n", "\\> caller\n"},
		{"bullet-list", "- item\n", "\\- item\n"},
		{"ordered-list", "1. item\n", "1\\. item\n"},
		{"fence-open-autoclose", "```\nuser code\n", "```\nuser code\n```\n"},
		{"fence-roundtrip", "```\nuser code\n```\n", "```\nuser code\n```\n"},
		// CM §4.5: backtick fences with a backtick in the info string
		// do NOT open. Without this check, the sanitizer thinks a fence
		// opened and skips defenses on subsequent lines while goldmark
		// keeps parsing block markers. Line 2 must therefore get the
		// usual line-leader escape.
		{"fence-backtick-in-info-rejected", "```a`b\n<gno-card>\n", "```a`b\n\\<gno-card>\n"},
		{"fence-tilde-info-string-allowed", "~~~lang~tag\nx\n", "~~~lang~tag\nx\n~~~\n"},
		// Bracket walker MUST also reject the backtick-info-string
		// "fence" — otherwise it treats lines below as opaque fence
		// interior and skips LRD strip + bracket escape, letting an
		// caller smuggle realm-targeted ref-link definitions past
		// the walker.
		{"fence-walker-backtick-in-info-rejected", "```a`b\n\n[evil]: https://bad\n\n[evil]\n", "```a`b\n\n\n\\[evil\\]\n"},
		{"setext-h1", "title\n===\n", "title\n\\===\n"},
		{"ref-link-use", "[click][evil]\n", "\\[click\\]\\[evil\\]\n"},
		{"shortcut-ref", "[label]\n", "\\[label\\]\n"},
		{"footnote-ref", "[^name]\n", "\\[^name\\]\n"},
		{"inline-link-preserved", "[text](url)\n", "[text](url)\n"},
		{"inline-image-preserved", "![alt](src)\n", "![alt](src)\n"},
		{"multi-line-lrd-stripped", "[lab\nel]: https://e.com\n\nbody\n", "\nbody\n"},
		{"backslash-escaped-lrd-not-stripped", "[label\\]: url\n", "\\[label\\]: url\n"},
		{"lrd-strip", "[evil]: https://bad\n", ""},
		{"u2028-fold", "a\u2028b\n", "a\nb\n"},
		{"nel-fold", "a\u0085b\n", "a\nb\n"},
		{"ext-delimiter", "<gno-card>\n", "\\<gno-card>\n"},
		{"ext-delimiter-uppercase", "<GNO-CARD>\n", "\\<GNO-CARD>\n"},        // case-insensitive
		{"ext-delimiter-mixed-case", "<Gno-Columns>\n", "\\<Gno-Columns>\n"}, // case-insensitive
		{"ext-delimiter-close-uppercase", "</GNO-COLUMNS>\n", "\\</GNO-COLUMNS>\n"},
		{"ext-delimiter-not-matched", "<gnu-card>\n", "<gnu-card>\n"}, // not `gno-`
		{"gfm-table-row", "| a | b |\n", "\\| a | b |\n"},
		// CM §4.6 HTML block types 1-5 — escaped (blank-line-NON-terminating).
		{"html-type1-script", "<script>x</script>\n", "\\<script>x</script>\n"},
		{"html-type1-pre", "<pre>x</pre>\n", "\\<pre>x</pre>\n"},
		{"html-type1-style", "<style>x</style>\n", "\\<style>x</style>\n"},
		{"html-type1-textarea", "<textarea>x</textarea>\n", "\\<textarea>x</textarea>\n"},
		{"html-type1-case-insensitive", "<SCRIPT>x</SCRIPT>\n", "\\<SCRIPT>x</SCRIPT>\n"},
		{"html-type1-self-closing", "<script/>x\n", "\\<script/>x\n"},
		{"html-type1-bare-eol", "<script\n", "\\<script\n"},
		{"html-type1-not-name-prefix", "<scripta>x\n", "<scripta>x\n"}, // no boundary after name
		{"html-type2-comment", "<!-- comment -->\n", "\\<!-- comment -->\n"},
		{"html-type2-degenerate", "<!---->\n", "\\<!---->\n"},
		{"html-type3-pi", "<?php x ?>\n", "\\<?php x ?>\n"},
		{"html-type4-doctype", "<!DOCTYPE html>\n", "\\<!DOCTYPE html>\n"},
		{"html-type4-lowercase-not-matched", "<!doctype html>\n", "<!doctype html>\n"},
		{"html-type4-bang-eol", "<!\n", "<!\n"}, // bounds-safe, no follow char
		// Type 5 CDATA: the bracket walker (escapeBracketsOutsideLinks)
		// runs BEFORE the per-line HTML detector, so it mangles the
		// `[` / `]` in `<![CDATA[…]]>` to `\[` / `\]`. Goldmark's
		// Type 5 regex (`<\!\[CDATA\[`) no longer matches the mangled
		// form, so no Type 5 block opens — defense-in-depth via the
		// walker, even without the new detector firing.
		{"html-type5-cdata", "<![CDATA[x]]>\n", "<!\\[CDATA\\[x\\]\\]>\n"},
		{"html-3-space-indent", "   <!-- x -->\n", "\\   <!-- x -->\n"},             // \ at byte 0; spaces preserved after
		{"html-4-space-indent-not-matched", "    <!-- x -->\n", "    <!-- x -->\n"}, // 4+ = indented code
		{"html-inline-not-matched", "paragraph <!-- inline -->\n", "paragraph <!-- inline -->\n"},
		{"html-crlf", "<!-- x -->\r\n", "\\<!-- x -->\r\n"},        // CRLF passes through at native level (NormalizeBreaks runs in sanitize wrappers)
		{"html-already-escaped", "\\<script>x\n", "\\<script>x\n"}, // backslash at byte 0 → no re-escape
	}
	for _, c := range cases {
		if got := EscapeBlockHazards(c.in); got != c.want {
			t.Errorf("%s: EscapeBlockHazards(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
		// Idempotency: applying EscapeBlockHazards twice must be
		// byte-identical to applying it once.
		if twice := EscapeBlockHazards(c.want); twice != c.want {
			t.Errorf("%s: not idempotent: EscapeBlockHazards(want) = %q, want %q", c.name, twice, c.want)
		}
	}
}

func TestEscapeBlockHazardsRich(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "hello world\n", "hello world\n"},
		// Doc-spoof markers PRESERVED in Rich mode.
		{"atx-heading", "# heading\n", "# heading\n"},
		{"blockquote", "> quoted\n", "> quoted\n"},
		{"bullet-list-dash", "- item\n", "- item\n"},
		{"bullet-list-star", "* item\n", "* item\n"},
		{"bullet-list-plus", "+ item\n", "+ item\n"},
		{"ordered-list", "1. item\n", "1. item\n"},
		{"thematic-break-dash", "---\n", "---\n"},
		{"thematic-break-star", "***\n", "***\n"},
		{"thematic-break-underscore", "___\n", "___\n"},
		{"setext-h1", "title\n===\n", "title\n===\n"},
		{"setext-h2", "title\n---\n", "title\n---\n"},
		// GFM tables PRESERVED in Rich mode (line-leading `|` not escaped).
		{"gfm-table-row", "| a | b |\n", "| a | b |\n"},
		{"gfm-table-full", "| H1 | H2 |\n|---|---|\n| 1 | 2 |\n", "| H1 | H2 |\n|---|---|\n| 1 | 2 |\n"},
		// Realm-binding defenses STILL ON.
		{"ext-delimiter", "<gno-card>\n", "\\<gno-card>\n"},
		{"ext-delimiter-uppercase", "<GNO-CARD>\n", "\\<GNO-CARD>\n"},
		{"ext-delimiter-mixed-case", "<Gno-Columns>\n", "\\<Gno-Columns>\n"},
		{"ext-delimiter-not-matched", "<gnu-card>\n", "<gnu-card>\n"},
		{"ref-link-use", "[click][evil]\n", "\\[click\\]\\[evil\\]\n"},
		{"shortcut-ref", "[label]\n", "\\[label\\]\n"},
		{"footnote-ref", "[^name]\n", "\\[^name\\]\n"},
		{"lrd-strip", "[evil]: https://bad\n", ""},
		{"escaped-lrd-not-stripped", "[label\\]: url\n", "\\[label\\]: url\n"},
		{"inline-link-preserved", "[t](url)\n", "[t](url)\n"},
		{"inline-image-preserved", "![alt](src)\n", "![alt](src)\n"},
		{"fence-open-autoclose", "```\nuser\n", "```\nuser\n```\n"},
		{"fence-backtick-in-info-rejected", "```a`b\n<gno-card>\n", "```a`b\n\\<gno-card>\n"},
		{"fence-tilde-info-string-allowed", "~~~lang~tag\nx\n", "~~~lang~tag\nx\n~~~\n"},
		// Bracket walker MUST also reject the backtick-info-string
		// "fence" — see strict variant above for explanation.
		{"fence-walker-backtick-in-info-rejected", "```a`b\n\n[evil]: https://bad\n\n[evil]\n", "```a`b\n\n\n\\[evil\\]\n"},
		{"u2028-fold", "a\u2028b\n", "a\nb\n"},
		{"nel-fold", "a\u0085b\n", "a\nb\n"},
		// CM §4.6 HTML block types 1-5 — escaped in Rich mode too
		// (defense is mode-independent — see escapeBlockHazardsImpl).
		{"html-type1-script", "<script>x</script>\n", "\\<script>x</script>\n"},
		{"html-type1-case-insensitive", "<SCRIPT>x</SCRIPT>\n", "\\<SCRIPT>x</SCRIPT>\n"},
		{"html-type1-bare-eol", "<script\n", "\\<script\n"},
		{"html-type1-not-name-prefix", "<scripta>x\n", "<scripta>x\n"},
		{"html-type2-comment", "<!-- comment -->\n", "\\<!-- comment -->\n"},
		{"html-type3-pi", "<?php x ?>\n", "\\<?php x ?>\n"},
		{"html-type4-doctype", "<!DOCTYPE html>\n", "\\<!DOCTYPE html>\n"},
		{"html-type4-lowercase-not-matched", "<!doctype html>\n", "<!doctype html>\n"},
		{"html-type5-cdata", "<![CDATA[x]]>\n", "<!\\[CDATA\\[x\\]\\]>\n"},
		{"html-3-space-indent", "   <!-- x -->\n", "\\   <!-- x -->\n"},
		{"html-4-space-indent-not-matched", "    <!-- x -->\n", "    <!-- x -->\n"},
		{"html-inline-not-matched", "paragraph <!-- inline -->\n", "paragraph <!-- inline -->\n"},
		{"html-already-escaped", "\\<script>x\n", "\\<script>x\n"},
	}
	for _, c := range cases {
		if got := EscapeBlockHazardsRich(c.in); got != c.want {
			t.Errorf("%s: EscapeBlockHazardsRich(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
		// Idempotency.
		if twice := EscapeBlockHazardsRich(c.want); twice != c.want {
			t.Errorf("%s: not idempotent: EscapeBlockHazardsRich(want) = %q, want %q", c.name, twice, c.want)
		}
	}
}

// BenchmarkEscapeBlockHazardsPathological exercises bracket-walker
// inputs that maximize backtrack work in pass 1: repeated well-formed
// links (linear path), repeated unclosed-link prefixes terminated by a
// blank line (current super-linear path; see note), and a multi-line
// LRD label that must scan across newlines. The intent is a wall-clock
// budget guard against future regressions.
//
// NOTE: `unclosed-xN` currently scales worse than linear because every
// `[` re-scans forward to the blank-line terminator. For N=500/1k/2k/4k
// the per-op cost grows ~3x per 2x input. Acceptable for realistic
// realm input sizes (a few KB) but flagged here for follow-up if a
// caller needs to defend against pathological inputs.
func BenchmarkEscapeBlockHazardsPathological(b *testing.B) {
	links := strings.Repeat("[a](b)", 1000)
	unclosed := strings.Repeat("[a", 1000) + "\n\nbody"
	mlLRD := "[" + strings.Repeat("lab\nel", 500) + "]: url\n"
	cases := []struct {
		name string
		in   string
	}{
		{"links-x1000", links},
		{"unclosed-x1000", unclosed},
		{"multiline-lrd-label", mlLRD},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(c.in)))
			for i := 0; i < b.N; i++ {
				_ = EscapeBlockHazards(c.in)
			}
		})
	}
}

// BenchmarkEscapeBlockHazardsRichPathological mirrors the strict
// variant's pathological-input guard for EscapeBlockHazardsRich. The
// bracket walker, scan-budget cap, and fence state machine are
// shared with EscapeBlockHazards — Rich mode only skips two
// per-line escapes — so the same pathological inputs are the
// right wall-clock budget guards. Adding the parallel benchmark
// catches regressions specific to the Rich code path (e.g. if a
// future refactor breaks the bracket walker's budget enforcement
// for mode=0 only).
func BenchmarkEscapeBlockHazardsRichPathological(b *testing.B) {
	links := strings.Repeat("[a](b)", 1000)
	unclosed := strings.Repeat("[a", 1000) + "\n\nbody"
	mlLRD := "[" + strings.Repeat("lab\nel", 500) + "]: url\n"
	parenTitle := strings.Repeat("[a]: u\n(x\n", 100) // calibration shape
	cases := []struct {
		name string
		in   string
	}{
		{"links-x1000", links},
		{"unclosed-x1000", unclosed},
		{"multiline-lrd-label", mlLRD},
		{"paren-title-x100", parenTitle},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(c.in)))
			for i := 0; i < b.N; i++ {
				_ = EscapeBlockHazardsRich(c.in)
			}
		})
	}
}
