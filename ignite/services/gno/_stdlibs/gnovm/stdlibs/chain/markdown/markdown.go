// Package markdown is the Go-side implementation of the chain/markdown
// gno stdlib. See markdown.gno for the public contract.
package markdown

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// maxForeignBlocksPerConvert caps the number of <gno-foreign> blocks a
// single markdown render admits (the gnoweb renderer drops openers
// beyond it). Single source of truth: the gnoweb foreign renderer reads
// it via MaxForeignBlocksPerConvert, and realms read the same value
// from gno via the native of the same name.
const maxForeignBlocksPerConvert = 256

// MaxForeignBlocksPerConvert returns maxForeignBlocksPerConvert. It
// backs the gno native of the same name (callable from realms) and is
// also the Go accessor the gnoweb foreign renderer uses.
func MaxForeignBlocksPerConvert() int { return maxForeignBlocksPerConvert }

// ---------- StripBidiAndZeroWidth ----------

func StripBidiAndZeroWidth(s string) string {
	// Fast path: all stripped codepoints lie in U+200B..U+FEFF, all of
	// which require a 3-byte UTF-8 sequence whose lead byte is 0xE2 or
	// 0xEF. Check the bytes directly (strings.ContainsAny is rune-based
	// and would mis-handle these continuation-expecting bytes).
	if !containsAnyByte(s, 0xe2, 0xef) {
		return s
	}
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		r, sz := utf8.DecodeRuneInString(s[i:])
		if isBidiOrZeroWidth(r) {
			i += sz
			continue
		}
		out = append(out, s[i:i+sz]...)
		i += sz
	}
	return string(out)
}

func isBidiOrZeroWidth(r rune) bool {
	switch {
	case r >= 0x200B && r <= 0x200F: // ZWSP, ZWNJ, ZWJ, LRM, RLM
		return true
	case r >= 0x202A && r <= 0x202E: // LRE, RLE, PDF, LRO, RLO
		return true
	case r >= 0x2066 && r <= 0x2069: // LRI, RLI, FSI, PDI
		return true
	case r == 0xFEFF: // BOM / ZWNBSP
		return true
	}
	return false
}

// ---------- NormalizeBreaks ----------

func NormalizeBreaks(s string) string {
	if !strings.ContainsAny(s, "\r") {
		return s
	}
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\r' {
			out = append(out, '\n')
			if i+1 < len(s) && s[i+1] == '\n' {
				i++ // skip the \n in \r\n
			}
			continue
		}
		out = append(out, s[i])
	}
	return string(out)
}

// ---------- EscapeInline / EscapeTitle ----------

var inlineEscapeSet [128]bool
var titleEscapeSet [128]bool

func init() {
	for _, c := range []byte{'\\', '*', '_', '[', ']', '(', ')', '~', '>', '-', '+', '.', '!', '`', '#', '<', '&'} {
		inlineEscapeSet[c] = true
	}
	titleEscapeSet = inlineEscapeSet
	titleEscapeSet['"'] = true
	titleEscapeSet['\''] = true
}

func EscapeInline(s string) string { return escapeWithSet(s, &inlineEscapeSet) }
func EscapeTitle(s string) string  { return escapeWithSet(s, &titleEscapeSet) }

func escapeWithSet(s string, set *[128]bool) string {
	// Pre-scan: skip allocation if no byte needs escaping and no NUL.
	work := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == 0 || (c < 128 && set[c]) {
			work = true
			break
		}
	}
	if !work {
		return s
	}
	out := make([]byte, 0, len(s)+8)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == 0 {
			out = append(out, "\xef\xbf\xbd"...) // U+FFFD
			continue
		}
		if c < 128 && set[c] {
			out = append(out, '\\')
		}
		out = append(out, c)
	}
	return string(out)
}

// ---------- PercentEncodeURL ----------

func PercentEncodeURL(s string) string {
	work := false
	for i := 0; i < len(s); i++ {
		if needsPercentEncode(s, i) {
			work = true
			break
		}
	}
	if !work {
		return s
	}
	out := make([]byte, 0, len(s)+16)
	for i := 0; i < len(s); i++ {
		if needsPercentEncode(s, i) {
			c := s[i]
			out = append(out, '%', hexDigit(c>>4), hexDigit(c&0xF))
			continue
		}
		out = append(out, s[i])
	}
	return string(out)
}

func needsPercentEncode(s string, i int) bool {
	c := s[i]
	if c <= 0x20 || c == 0x7F || c >= 0x80 {
		return true
	}
	switch c {
	case '"', '\'', '(', ')', '<', '>', '\\', '`', '{', '|', '}', '^':
		return true
	case '&':
		return startsEntityReference(s, i)
	case '%':
		// Bare % (not followed by two hex digits) gets encoded.
		if i+2 >= len(s) || !isHex(s[i+1]) || !isHex(s[i+2]) {
			return true
		}
	}
	return false
}

// skipEscapeBefore returns i advanced past a `\` that escapes want.
// UnescapePunctuations runs BEFORE the reference resolvers, so `&\#x6a;`
// and `&#x6a\;` both reach a resolver as `&#x6a;`; matching the reference
// shape literally would let the escape hide it from us and reveal it to
// the renderer. Only ASCII punctuation is escapable, which is why the
// byte is checked: in `&\name;` the `\n` survives unescaping, so no
// reference forms and the `&` stays bare.
func skipEscapeBefore(s string, i int, want byte) int {
	if i+1 < len(s) && s[i] == '\\' && s[i+1] == want {
		return i + 1
	}
	return i
}

// endsReference reports whether the `;` closing a reference sits at i,
// allowing for the `\;` form.
func endsReference(s string, i int) bool {
	i = skipEscapeBefore(s, i, ';')
	return i < len(s) && s[i] == ';'
}

func startsEntityReference(s string, i int) bool {
	i = skipEscapeBefore(s, i+1, '#')
	if i >= len(s) {
		return false
	}
	if s[i] == '#' {
		i++
		if i < len(s) && (s[i] == 'x' || s[i] == 'X') {
			i++
			start := i
			for i < len(s) && isHex(s[i]) {
				i++
			}
			return i > start && endsReference(s, i)
		}
		start := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		return i > start && endsReference(s, i)
	}
	if !isASCIIAlpha(s[i]) {
		return false
	}
	for i++; i < len(s) && isASCIIAlphaNumeric(s[i]); i++ {
	}
	return endsReference(s, i)
}

func percentEncodeEntityReferences(s string) string {
	first := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '&' && startsEntityReference(s, i) {
			first = i
			break
		}
	}
	if first < 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteString(s[:first])
	for i := first; i < len(s); i++ {
		if s[i] == '&' && startsEntityReference(s, i) {
			b.WriteString("%26")
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func isASCIIAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isASCIIAlphaNumeric(c byte) bool {
	return isASCIIAlpha(c) || c >= '0' && c <= '9'
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')
}

func hexDigit(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	return 'A' + (n - 10)
}

// ---------- MatchCharsetN ----------

func MatchCharsetN(s string, firstLo, firstHi, restLo, restHi uint64, minLen, maxLen int) bool {
	n := len(s)
	if n < minLen || n > maxLen {
		return false
	}
	if n == 0 {
		return minLen == 0
	}
	if !inBitmap(s[0], firstLo, firstHi) {
		return false
	}
	for i := 1; i < n; i++ {
		if !inBitmap(s[i], restLo, restHi) {
			return false
		}
	}
	return true
}

func inBitmap(c byte, lo, hi uint64) bool {
	if c >= 128 {
		return false
	}
	if c < 64 {
		return (lo>>uint(c))&1 != 0
	}
	return (hi>>uint(c-64))&1 != 0
}

// ---------- CodeFence ----------

func CodeFence(content string, minCount int) string {
	if minCount < 1 {
		minCount = 1
	}
	longest, cur := 0, 0
	for i := 0; i < len(content); i++ {
		if content[i] == '`' {
			cur++
			if cur > longest {
				longest = cur
			}
		} else {
			cur = 0
		}
	}
	n := max(longest+1, minCount)
	out := make([]byte, n)
	for i := range out {
		out[i] = '`'
	}
	return string(out)
}

// ---------- EscapeBlockHazards ----------

// blockHazardsMode selects which doc-spoof and table-injection
// defenses run inside escapeBlockHazardsImpl. Realm-binding defenses
// that stay unconditional in BOTH strict and Rich modes (bracket
// walker, <gno-…> extension delimiter, CM §4.6 HTML block types 1-5
// openers, fenced-code-block state machine, U+2028/U+2029/U+0085
// fold) are not represented as flags — they always fire.
type blockHazardsMode int

const (
	modeEscapeLineLeader blockHazardsMode = 1 << iota
	modeEscapeSetext
	modeEscapePipe

	// modeStrict is the strict-mode policy: all doc-spoof defenses
	// on. ANY new bit added to blockHazardsMode MUST be OR'd in here
	// so EscapeBlockHazards's strict posture stays the union of all
	// defenses. Forgetting to extend this constant silently weakens
	// strict callers (sanitize.Block) — there is no compile-time or
	// test-level catch for the omission besides per-defense unit
	// tests.
	modeStrict = modeEscapeLineLeader | modeEscapeSetext | modeEscapePipe
)

// EscapeBlockHazards is the strict variant (used by sanitize.Block).
// All doc-spoof, setext, and GFM-table-row defenses are on.
func EscapeBlockHazards(s string) string {
	return escapeBlockHazardsImpl(s, modeStrict)
}

// EscapeBlockHazardsRich is the permissive variant (used by
// sanitize.BlockRich). Line-leader (#, >, list markers, thematic
// breaks), setext-underline, and GFM table-row `|` escapes are all
// skipped — the user can compose multi-section markdown structure
// including tables. Realm-binding defenses (bracket walker,
// <gno-…> extension delimiters, CM §4.6 HTML block types 1-5
// openers, fenced-code-block state machine, U+2028/U+2029/U+0085
// fold) stay on — these are mode-independent security defenses,
// not stylistic preferences. NUL→U+FFFD replacement and
// bidi/zero-width strip run at the Gno layer (sanitize.BlockRich)
// before reaching this native, not here.
//
// Cross-paragraph promotion (user `===`/`---` setext or
// `|---|---|` table-separator at start or end of input reaching
// into adjacent realm chrome) is neutralized at the Gno layer by
// sanitize.BlockRich emitting `\n\n` (CM blank line, i.e.
// paragraph break) on BOTH sides of the user content — symmetric
// isolation against backward and forward attacks.
func EscapeBlockHazardsRich(s string) string {
	return escapeBlockHazardsImpl(s, 0)
}

func escapeBlockHazardsImpl(s string, mode blockHazardsMode) string {
	if s == "" {
		return s
	}
	// Fold Unicode separators U+2028, U+2029, U+0085 NEL to \n before
	// line splitting. These are not CM §2.2 line endings, so NormalizeBreaks
	// doesn't touch them; in block context we treat them as paragraph-
	// internal breaks.
	s = foldUnicodeSeparators(s)

	// Pass 1+2: bracket walker. Finds inline link / image / LRD / fence
	// spans on the whole input, then escapes any unescaped `[` / `]`
	// outside those spans, and deletes LRD spans entirely. Subsumes the
	// previous per-line LRD strip and ref-link-use escape.
	s = escapeBracketsOutsideLinks(s)

	var out strings.Builder
	out.Grow(len(s) + 16)

	lines := strings.Split(s, "\n")
	trailingNewline := strings.HasSuffix(s, "\n")
	if trailingNewline {
		lines = lines[:len(lines)-1] // strip artifact of split
	}

	var (
		inFence      bool
		fenceChar    byte
		fenceLen     int
		prevNonBlank bool
	)
	escapeLeader := mode&modeEscapeLineLeader != 0
	escapeSetext := mode&modeEscapeSetext != 0
	escapePipe := mode&modeEscapePipe != 0

	for idx, line := range lines {
		writeNL := idx < len(lines)-1 || trailingNewline

		if inFence {
			out.WriteString(line)
			if writeNL {
				out.WriteByte('\n')
			}
			if isCloseFence(line, fenceChar, fenceLen) {
				inFence = false
			}
			prevNonBlank = line != ""
			continue
		}

		// Extension delimiter lines: prefix with backslash so the line's
		// opening `<` becomes a CM §2.4 inline escape. A leading space
		// would be ineffective because gnoweb's extension parsers call
		// util.TrimLeftSpace before tag matching (see ext_columns.go,
		// ext_alert.go etc.) — the space gets stripped and the parser
		// sees the bare tag. Backslash survives util.TrimLeftSpace
		// (which only strips ASCII whitespace and form-feed) and
		// goldmark's Type-7 HTML block detection (which requires the
		// first non-whitespace char to be `<`, not `\`).
		if isExtDelimiter(line) {
			out.WriteByte('\\')
			out.WriteString(line)
			if writeNL {
				out.WriteByte('\n')
			}
			prevNonBlank = true
			continue
		}

		// CM §4.6 HTML block types 1-5: prefix with backslash so the
		// leading `<` becomes a CM §2.4 inline escape and goldmark's
		// HTML block parser (which requires the first non-whitespace
		// char to be `<`) refuses to open. Types 1-5 are the
		// blank-line-NON-terminating shapes (`<script>`, `<pre>`,
		// `<style>`, `<textarea>`, `<!--`, `<?…?>`, `<!DOCTYPE…>`,
		// `<![CDATA[…]]>`) — without this defense user content opening
		// any of them would consume realm chrome appended afterward
		// (until a type-specific close token or EOF). Types 6 and 7
		// close on a blank line, so sanitize.BlockRich's "\n\n"
		// envelope already bounds them; we don't escape those.
		// Fires unconditionally in both strict and Rich modes — this
		// is a security defense, not stylistic.
		if isHTMLBlockType1to5Opener(line) {
			out.WriteByte('\\')
			out.WriteString(line)
			if writeNL {
				out.WriteByte('\n')
			}
			prevNonBlank = true
			continue
		}

		// Setext underline (only if previous line was non-blank, and
		// the strict mode is enabled).
		if escapeSetext && prevNonBlank && isSetextUnderline(line) {
			out.WriteByte('\\')
			out.WriteString(line)
			if writeNL {
				out.WriteByte('\n')
			}
			prevNonBlank = true
			continue
		}

		// Block markers / fence open. Always called: fence detection is
		// unconditional. `escapeLeader` gates the doc-spoof markers
		// (#, >, list, HR) and `escapePipe` gates the GFM table-row
		// `|` escape; both are on in strict mode and off in Rich.
		escaped, fc, fl := escapeLineLeader(line, escapeLeader, escapePipe)
		out.WriteString(escaped)
		if writeNL {
			out.WriteByte('\n')
		}
		if fc != 0 {
			inFence = true
			fenceChar = fc
			fenceLen = fl
		}
		prevNonBlank = line != ""
	}

	// Auto-close any open code fence at end-of-input.
	if inFence {
		if out.Len() > 0 && out.String()[out.Len()-1] != '\n' {
			out.WriteByte('\n')
		}
		for i := 0; i < fenceLen; i++ {
			out.WriteByte(fenceChar)
		}
		out.WriteByte('\n')
	}

	return out.String()
}

// foldUnicodeSeparators replaces U+2028, U+2029 (3-byte UTF-8 starting
// with 0xE2 0x80) and U+0085 NEL (2-byte UTF-8 0xC2 0x85) with '\n'.
func foldUnicodeSeparators(s string) string {
	if !containsAnyByte(s, 0xe2, 0xc2) {
		return s
	}
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		r, sz := utf8.DecodeRuneInString(s[i:])
		if r == 0x2028 || r == 0x2029 || r == 0x0085 {
			out = append(out, '\n')
			i += sz
			continue
		}
		out = append(out, s[i:i+sz]...)
		i += sz
	}
	return string(out)
}

// isExtDelimiter recognises any gnoweb structural-extension delimiter
// shape via the open/close `<gno-...>` / `</gno-...>` prefixes. The
// wildcard is intentionally over-permissive: a malformed line like
// `<gno-x asdf` still trips it, which is fine — the only effect is
// that the line gets a leading backslash prepended (rendered as a
// literal `<` per CM §2.4). Future extensions (`<gno-card>`,
// `<gno-foreign>`, anything later) auto-cover without needing a
// sanitize-side update.
//
// Match is case-INsensitive (`<GNO-Card>`, `<Gno-COLUMNS>`, etc. all
// trip). Go's html.Tokenizer (used by the extension block parsers in
// gnoweb at ext_columns.go, ext_alert.go, etc.) lowercases tag names
// before the per-extension matcher runs, so an uppercase or mixed-
// case opener still opens the block. The sanitizer therefore must
// match the same byte-shape envelope the parsers do — otherwise
// `<GNO-columns>` slips past the sanitizer and opens a columns
// container in goldmark, swallowing realm chrome.
//
// Bare `|||` (the legacy `<gno-columns>` shorthand) is intentionally
// NOT matched here — the shorthand has been removed from the columns
// parser, so user content writing `|||` is now harmless paragraph
// text and doesn't need neutralisation.
func isExtDelimiter(line string) bool {
	trim := strings.TrimLeft(line, " \t")
	if len(trim) == 0 || trim[0] != '<' {
		return false
	}
	rest := trim[1:]
	// Optional `/` for close tags.
	if len(rest) > 0 && rest[0] == '/' {
		rest = rest[1:]
	}
	return hasCaseInsensitivePrefix(rest, "gno-")
}

// isHTMLBlockType1to5Opener reports whether line opens a CommonMark
// §4.6 HTML block of type 1, 2, 3, 4, or 5 — the types that do NOT
// close on a blank line. Types 6 and 7 close on a blank line per
// goldmark's `Continue` (parser/html_block.go), so the leading +
// trailing "\n\n" envelope in sanitize.BlockRich already bounds them;
// we only need to neutralize 1-5 here.
//
// Detection mirrors goldmark's regexes at parser/html_block.go:79-92
// (case-insensitive Type 1 tag name; case-sensitive Types 2/3/4/5;
// ASCII space-only indent of 0-3 columns — tabs are deliberately NOT
// allowed since goldmark uses `[ ]{0,3}` literal, divergent from
// isExtDelimiter's `TrimLeft(line, " \t")`).
//
// Type 4 uses `[A-Z]` (one uppercase letter) where goldmark's regex
// is `[A-Z]+` — equivalent for opener detection because any line
// goldmark would accept under `+` also satisfies the single-letter
// check.
func isHTMLBlockType1to5Opener(line string) bool {
	// 0-3 ASCII-space indent.
	i := 0
	for i < 3 && i < len(line) && line[i] == ' ' {
		i++
	}
	if i >= len(line) || line[i] != '<' {
		return false
	}
	rest := line[i+1:]
	if len(rest) == 0 {
		return false
	}
	// Type 2: <!--
	if strings.HasPrefix(rest, "!--") {
		return true
	}
	// Type 5: <![CDATA[ (case-sensitive)
	if strings.HasPrefix(rest, "![CDATA[") {
		return true
	}
	// Type 4: <![A-Z]
	if len(rest) >= 2 && rest[0] == '!' && rest[1] >= 'A' && rest[1] <= 'Z' {
		return true
	}
	// Type 3: <?
	if rest[0] == '?' {
		return true
	}
	// Type 1: <(script|pre|style|textarea) followed by \s, >, /, or EOL.
	for _, name := range [...]string{"script", "pre", "style", "textarea"} {
		if hasCaseInsensitivePrefix(rest, name) {
			after := rest[len(name):]
			if len(after) == 0 {
				return true
			}
			c := after[0]
			if c == ' ' || c == '\t' || c == '>' || c == '/' {
				return true
			}
		}
	}
	return false
}

// hasBacktickBeforeNewline reports whether s[from:] contains a `
// before the next `\n` or EOF. Used to enforce CM §4.5: a backtick
// fence opener whose info string contains another backtick does NOT
// open a fence. Three call sites in this file share this predicate
// (escapeLineLeader, isBlockInterrupt, findBracketSpans); keeping
// them on one helper prevents drift if the rule ever evolves.
func hasBacktickBeforeNewline(s string, from int) bool {
	for i := from; i < len(s) && s[i] != '\n'; i++ {
		if s[i] == '`' {
			return true
		}
	}
	return false
}

// hasCaseInsensitivePrefix reports whether s begins with prefix using
// ASCII-only case folding (A-Z → a-z). prefix MUST be lowercase ASCII;
// Type 1 tag names (`script`, `pre`, `style`, `textarea`) all qualify.
// No Unicode case folding — CM's HTML block 1 names are ASCII-only.
func hasCaseInsensitivePrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c != prefix[i] {
			return false
		}
	}
	return true
}

// isSetextUnderline reports whether line is shaped like ^ {0,3}=+[ \t]*$
// or ^ {0,3}-+[ \t]*$. Only valid following a non-blank line.
func isSetextUnderline(line string) bool {
	i := 0
	for i < len(line) && i < 3 && line[i] == ' ' {
		i++
	}
	if i >= len(line) {
		return false
	}
	marker := line[i]
	if marker != '=' && marker != '-' {
		return false
	}
	for i < len(line) && line[i] == marker {
		i++
	}
	for i < len(line) {
		if line[i] != ' ' && line[i] != '\t' {
			return false
		}
		i++
	}
	return true
}

// escapeLineLeader escapes a single line-leading block marker by
// prefixing it with \. Returns (possibly-escaped line, fence-char,
// fence-len) — fence char/len are returned when the line opens a
// code fence (length >= 3).
//
// The `escapeLeader` flag gates the doc-spoof markers (#, >, list,
// HR); when false, those markers are preserved verbatim. The
// `escapePipe` flag gates the GFM table-row `|` escape; when false,
// line-leading `|` is preserved so user content can render as a
// table. Fenced-code-block detection is always on (state machine
// tracking, not a defense).
func escapeLineLeader(line string, escapeLeader, escapePipe bool) (string, byte, int) {
	i := 0
	for i < len(line) && i < 3 && line[i] == ' ' {
		i++
	}
	if i >= len(line) {
		return line, 0, 0
	}
	c := line[i]
	switch c {
	case '#':
		// ATX heading: # to ###### followed by space, EOL, or tab
		if escapeLeader {
			j := i
			for j < len(line) && j-i < 6 && line[j] == '#' {
				j++
			}
			if j-i >= 1 && j-i <= 6 && (j == len(line) || line[j] == ' ' || line[j] == '\t') {
				return line[:i] + "\\" + line[i:], 0, 0
			}
		}
	case '>':
		if escapeLeader {
			return line[:i] + "\\" + line[i:], 0, 0
		}
	case '-', '*', '_':
		if escapeLeader {
			// Thematic break: 3+ of -, *, or _ optionally separated by spaces.
			if isThematicBreak(line, i) {
				return line[:i] + "\\" + line[i:], 0, 0
			}
			// Bullet list marker: - or * followed by space.
			if (c == '-' || c == '*') && i+1 < len(line) && (line[i+1] == ' ' || line[i+1] == '\t') {
				return line[:i] + "\\" + line[i:], 0, 0
			}
		}
	case '+':
		if escapeLeader && i+1 < len(line) && (line[i+1] == ' ' || line[i+1] == '\t') {
			return line[:i] + "\\" + line[i:], 0, 0
		}
	case '|':
		// GFM table-row line leader. gnoweb loads
		// extension.Table (see render_config.go), so a line-leading
		// `|` followed by a delimiter line forms a real table.
		// Block keeps tables out of user content (strict posture);
		// BlockRich preserves them so authors can compose tables.
		if escapePipe {
			return line[:i] + "\\" + line[i:], 0, 0
		}
	case '`', '~':
		// Fenced code block: 3+ of ` or ~. Code fences are legitimate
		// user content (users want to share code); do NOT escape them.
		// Just mark the fence as open so the EOF auto-close can fire if
		// the user never closes it.
		//
		// CM §4.5: for BACKTICK fences only, the info string MUST NOT
		// contain another backtick. Goldmark enforces this — if the
		// info string contains a `, no fence opens (the line is paragraph
		// text). The sanitizer MUST mirror this exactly, or it would
		// treat the line as a fence open while goldmark treats subsequent
		// lines as paragraph content where line-leader / extension /
		// HTML-block defenses normally apply. Without the check, an
		// caller can write "```a`b" to make the sanitizer believe a
		// fence opened and skip defenses on `<gno-…>`, `<!--`, `===`,
		// `#`, `|---|`, etc. on the following lines, while goldmark
		// happily opens those blocks. Tilde fences are NOT subject to
		// this rule (tildes ARE allowed in tilde-fence info strings).
		j := i
		for j < len(line) && line[j] == c {
			j++
		}
		if j-i >= 3 {
			if c == '`' && hasBacktickBeforeNewline(line, j) {
				return line, 0, 0
			}
			return line, c, j - i
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Ordered list marker: 1-9 digits followed by . or ) and space.
		if escapeLeader {
			j := i
			for j < len(line) && j-i < 9 && line[j] >= '0' && line[j] <= '9' {
				j++
			}
			if j < len(line) && (line[j] == '.' || line[j] == ')') &&
				j+1 < len(line) && (line[j+1] == ' ' || line[j+1] == '\t') {
				// Escape the digit run by prefixing the delimiter with \.
				return line[:j] + "\\" + line[j:], 0, 0
			}
		}
	}
	return line, 0, 0
}

func isThematicBreak(line string, start int) bool {
	c := line[start]
	if c != '-' && c != '*' && c != '_' {
		return false
	}
	count := 0
	for i := start; i < len(line); i++ {
		if line[i] == c {
			count++
		} else if line[i] != ' ' && line[i] != '\t' {
			return false
		}
	}
	return count >= 3
}

// isCloseFence reports whether line is a valid close fence for the
// current open fence (same char, length >= open, optional trailing
// whitespace, leading whitespace <= 3).
func isCloseFence(line string, fenceChar byte, fenceLen int) bool {
	i := 0
	for i < len(line) && i < 3 && line[i] == ' ' {
		i++
	}
	count := 0
	for i < len(line) && line[i] == fenceChar {
		count++
		i++
	}
	if count < fenceLen {
		return false
	}
	for i < len(line) {
		if line[i] != ' ' && line[i] != '\t' {
			return false
		}
		i++
	}
	return true
}

// containsAnyByte reports whether s contains any of the given bytes.
// Byte-level (not rune-level) — required for UTF-8 lead-byte fast paths.
func containsAnyByte(s string, bs ...byte) bool {
	for i := 0; i < len(s); i++ {
		if slices.Contains(bs, s[i]) {
			return true
		}
	}
	return false
}

// ---------- escapeBracketsOutsideLinks (replaces escapeRefLinkUse + isLRDDefinition) ----------
//
// Two-pass walker:
//   Pass 1: find spans for inline links [text](url), images ![alt](src),
//     LRD definitions [label]: url [title], and fenced code regions. Multi-
//     line and backslash-aware.
//   Pass 2: rewrite bytes — encode entity references in link/image URLs,
//     preserve fence spans, delete LRD spans, and escape any unescaped `[` / `]`
//     outside spans (with backslash-parity tracking).
//
// Closes three previously-documented residuals:
//   1. Shortcut-ref `[label]` collision with realm-emitted LRDs (both
//      brackets now escaped → literal text)
//   2. Multi-line LRD evasion `[lab\nel]: url` (label spans newlines)
//   3. False-positive over-strip `[label\]: url` (escape-aware `]` scan)

type spanKind byte

const (
	spanLink  spanKind = iota // [text](url) or ![alt](src) — encode entity refs in destination
	spanLRD                   // [label]: url ["title"] — delete entirely
	spanFence                 // fenced code block — preserve (opaque)
)

type bracketSpan struct {
	start, end       int // [start, end) half-open byte indices
	urlStart, urlEnd int // destination bounds for spanLink
	kind             spanKind
}

// findBracketSpans is pass 1 of the bracket walker. Returns sorted,
// non-overlapping spans.
func findBracketSpans(s string) []bracketSpan {
	var spans []bracketSpan
	i := 0
	atLineStart := true
	// nextClose is the byte index of the next `]` at or after i, or
	// len(s) if no `]` remains. Maintained lazily — we only need to
	// recompute it when i catches up. Used to short-circuit
	// scanLinkText on pathological input like `[[[[…` with no closer,
	// where every `[` would otherwise force a fresh forward walk to
	// the next blank line or EOF (O(n²) total). With this check,
	// each `[` either has a `]` ahead and we run the real scan, or
	// no `]` remains and we skip the scan in O(1).
	nextClose := strings.IndexByte(s, ']')
	if nextClose < 0 {
		nextClose = len(s)
	}
	// scanBudget bounds the *total* work scanLinkText can do across all
	// `[`s in this input, defending against the `[[[[…]]]]` shape where
	// every `[` could otherwise force a depth-balanced walk across most
	// of the input. 8×len(s) leaves plenty of headroom for legitimate
	// content (real links are short and rare relative to input size)
	// while capping worst-case work at O(n).
	scanBudget := 8 * len(s)
	for i < len(s) {
		// Advance nextClose past i if needed.
		if nextClose < i {
			rel := strings.IndexByte(s[i:], ']')
			if rel < 0 {
				nextClose = len(s)
			} else {
				nextClose = i + rel
			}
		}
		// Fence open at line-start (0-3 lead spaces, ≥3 backticks or tildes)?
		if atLineStart {
			j := i
			// skip up to 3 leading spaces
			for j < len(s) && j-i < 3 && s[j] == ' ' {
				j++
			}
			if j < len(s) && (s[j] == '`' || s[j] == '~') {
				fenceChar := s[j]
				k := j
				for k < len(s) && s[k] == fenceChar {
					k++
				}
				if k-j >= 3 {
					// CM §4.5: backtick fences with a backtick in the
					// info string do NOT open. Without this check the
					// walker treats subsequent lines as fence interior
					// (skipping LRD strip + bracket escape) while
					// goldmark treats them as paragraph — letting an
					// caller smuggle a ref-link definition past the
					// walker.
					if fenceChar == '`' && hasBacktickBeforeNewline(s, k) {
						i = k
						atLineStart = false
						continue
					}
					// Fence opens. Find close.
					fenceLen := k - j
					end := findFenceClose(s, k, fenceChar, fenceLen)
					spans = append(spans, bracketSpan{start: i, end: end, kind: spanFence})
					i = end
					atLineStart = (end > 0 && s[end-1] == '\n')
					continue
				}
			}
		}

		c := s[i]
		if c == '\\' && i+1 < len(s) {
			// Escaped byte — consume both
			i += 2
			atLineStart = false
			continue
		}
		if c == '\n' {
			atLineStart = true
			i++
			continue
		}
		// Try inline link / image / LRD at this `[` or `![`
		if c == '[' || (c == '!' && i+1 < len(s) && s[i+1] == '[') {
			start := i
			openOff := 0
			if c == '!' {
				openOff = 1
			}
			if nextClose >= len(s) {
				// No `]` remains anywhere ahead; scanLinkText will fail.
				// Skip the scan to keep the walker linear-time even on
				// pathological unclosed-bracket input.
				i++
				atLineStart = false
				continue
			}
			textEnd, ok := scanLinkText(s, i+openOff, &scanBudget)
			if !ok {
				i++
				atLineStart = false
				continue
			}
			// textEnd points at the closing `]`. Next byte determines what kind.
			if textEnd+1 < len(s) && s[textEnd+1] == '(' && c != '!' || // [text](url) — link
				textEnd+1 < len(s) && s[textEnd+1] == '(' && c == '!' { // ![alt](src) — image
				end, urlStart, urlEnd, ok := scanLinkURL(s, textEnd+2, &scanBudget)
				if ok {
					spans = append(spans, bracketSpan{start: start, end: end, urlStart: urlStart, urlEnd: urlEnd, kind: spanLink})
					i = end
					atLineStart = false
					continue
				}
			}
			// LRD candidate: [label]: url …  (only when `!` not prefixing AND only if line starts at line-start position)
			if c == '[' && atLineStartAt(s, start) {
				if end, ok := scanLRDTail(s, textEnd+1, &scanBudget); ok {
					spans = append(spans, bracketSpan{start: start, end: end, kind: spanLRD})
					i = end
					atLineStart = (end > 0 && s[end-1] == '\n')
					continue
				}
			}
			// Not a link, image, or LRD — advance past the `[` only
			i++
			atLineStart = false
			continue
		}
		// Any other byte
		if c != ' ' && c != '\t' {
			atLineStart = false
		}
		i++
	}
	return spans
}

// atLineStartAt reports whether byte index pos is at column 0-3 of a line
// (i.e. preceded only by 0-3 spaces since the last `\n` or start-of-input).
func atLineStartAt(s string, pos int) bool {
	// Walk back to find newline or start
	j := pos - 1
	spaces := 0
	for j >= 0 && s[j] == ' ' {
		spaces++
		if spaces > 3 {
			return false
		}
		j--
	}
	return j < 0 || s[j] == '\n'
}

// findFenceClose returns the byte index AFTER the closing fence line, or
// len(s) if no close is found (treat-as-opaque-to-EOF). `from` is the byte
// index right after the open fence's last fence char.
func findFenceClose(s string, from int, fenceChar byte, fenceLen int) int {
	// Skip rest of open fence's line
	i := from
	for i < len(s) && s[i] != '\n' {
		i++
	}
	if i < len(s) {
		i++ // past the newline
	}
	// Scan lines
	for i < len(s) {
		// line start
		j := i
		// 0-3 lead spaces
		for j < len(s) && j-i < 3 && s[j] == ' ' {
			j++
		}
		// count fence chars
		k := j
		for k < len(s) && s[k] == fenceChar {
			k++
		}
		if k-j >= fenceLen {
			// rest of line must be whitespace
			rest := k
			for rest < len(s) && s[rest] != '\n' && (s[rest] == ' ' || s[rest] == '\t') {
				rest++
			}
			if rest >= len(s) || s[rest] == '\n' {
				if rest < len(s) {
					rest++
				}
				return rest
			}
		}
		// Not a close fence — skip to next line
		for i < len(s) && s[i] != '\n' {
			i++
		}
		if i < len(s) {
			i++
		}
	}
	return len(s)
}

// scanLinkText scans `[…]` starting at s[i] where s[i] == '['. Returns
// the byte index of the closing `]` (depth-balanced, escape-aware) and
// true on success. `budget` is a shared work counter — each byte
// inspected decrements it, and the scan aborts (returns false) when
// it hits zero. This caps total pass-1 scan work at O(budget) across
// all calls, bounding pathological inputs like
// `[[[[…]]]]` where every `[` would otherwise force a fresh
// depth-balanced walk across most of the input (O(n²)).
func scanLinkText(s string, i int, budget *int) (int, bool) {
	if i >= len(s) || s[i] != '[' {
		return 0, false
	}
	j := i + 1
	depth := 1
	for j < len(s) {
		if *budget <= 0 {
			return 0, false
		}
		*budget--
		c := s[j]
		if c == '\\' && j+1 < len(s) {
			j += 2
			continue
		}
		switch c {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return j, true
			}
		case '\n':
			// Single \n OK inside link text; blank line aborts
			if j+1 < len(s) && s[j+1] == '\n' {
				return 0, false
			}
		}
		j++
	}
	return 0, false
}

// scanLinkURL scans `(url ["title"])` body starting AFTER the opening `(`.
// Returns the span end and destination bounds on success.
// Shares the pass-1 scan budget with scanLinkText so a caller can't
// chain `[a]([a]([a](…` (which would otherwise let scanLinkURL walk to
// EOF on every `[`, O(n²)).
func scanLinkURL(s string, i int, budget *int) (end, urlStart, urlEnd int, ok bool) {
	// Consume optional leading whitespace
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	if i < len(s) && s[i] == '\n' {
		i++
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
	}
	// URL body — angle-bracket form or plain
	if i < len(s) && s[i] == '<' {
		i++
		urlStart = i
		for i < len(s) {
			if *budget <= 0 {
				return 0, 0, 0, false
			}
			*budget--
			c := s[i]
			if c == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if c == '>' {
				urlEnd = i
				i++
				break
			}
			if c == '<' {
				return 0, 0, 0, false
			}
			if c == '\n' && i+1 < len(s) && s[i+1] == '\n' {
				return 0, 0, 0, false
			}
			i++
		}
	} else {
		// Plain URL with balanced parens
		urlStart = i
		urlDepth := 1
		for i < len(s) {
			if *budget <= 0 {
				return 0, 0, 0, false
			}
			*budget--
			c := s[i]
			if c == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if c == '(' {
				urlDepth++
			} else if c == ')' {
				urlDepth--
				if urlDepth == 0 {
					return i + 1, urlStart, i, true
				}
			} else if c == ' ' || c == '\t' || c == '\n' {
				// Whitespace after URL — title or close may follow
				break
			}
			i++
		}
		urlEnd = i
	}
	// Skip whitespace before title or close
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n') {
		if s[i] == '\n' && i+1 < len(s) && s[i+1] == '\n' {
			return 0, 0, 0, false
		}
		i++
	}
	// Optional title
	if i < len(s) && (s[i] == '"' || s[i] == '\'' || s[i] == '(') {
		closeQ := s[i]
		if closeQ == '(' {
			closeQ = ')'
		}
		i++
		for i < len(s) {
			if *budget <= 0 {
				return 0, 0, 0, false
			}
			*budget--
			c := s[i]
			if c == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if c == closeQ {
				i++
				break
			}
			if c == '\n' && i+1 < len(s) && s[i+1] == '\n' {
				return 0, 0, 0, false
			}
			i++
		}
		// Whitespace before close
		for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n') {
			i++
		}
	}
	if i < len(s) && s[i] == ')' {
		return i + 1, urlStart, urlEnd, true
	}
	return 0, 0, 0, false
}

// scanLRDTail scans the `: url [title]` portion of an LRD definition.
// `i` is the byte index right after the closing `]` of the label.
// Returns the byte index of the LRD region's end (after the URL or after
// the title, including a trailing newline if present) and true on success.
// Shares the pass-1 scan budget so a caller can't chain
// `[a]: u\n(x\n[b]: u\n(x\n…` (paren-title with no `)` anywhere — every
// LRD candidate would otherwise let the title scan walk to EOF, O(n²)).
func scanLRDTail(s string, i int, budget *int) (int, bool) {
	// Need ':' immediately
	if i >= len(s) || s[i] != ':' {
		return 0, false
	}
	i++
	// Need ≥1 space or tab, OR a newline (LRD URL can be on the next line)
	wsStart := i
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	if i < len(s) && s[i] == '\n' {
		i++
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
	}
	if i == wsStart {
		// `:` not followed by whitespace at all — not a valid LRD
		return 0, false
	}
	// URL bytes (no whitespace unless <...> form)
	urlStart := i
	if i < len(s) && s[i] == '<' {
		i++
		for i < len(s) {
			if *budget <= 0 {
				return 0, false
			}
			*budget--
			c := s[i]
			if c == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if c == '>' {
				i++
				break
			}
			if c == '<' || c == '\n' {
				return 0, false
			}
			i++
		}
	} else {
		for i < len(s) {
			if *budget <= 0 {
				return 0, false
			}
			*budget--
			c := s[i]
			if c == ' ' || c == '\t' || c == '\n' {
				break
			}
			if c == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			i++
		}
	}
	if i == urlStart {
		// Empty destination — not a valid LRD per CM §4.7
		return 0, false
	}
	end := i
	// Optional title: same line OR next line
	titleStart := i
	// Skip whitespace (incl. one newline)
	sawNL := false
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	if i < len(s) && s[i] == '\n' {
		sawNL = true
		i++
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
	}
	if i < len(s) && (s[i] == '"' || s[i] == '\'' || s[i] == '(') {
		closeQ := s[i]
		if closeQ == '(' {
			closeQ = ')'
		}
		i++
		titleOK := false
		for i < len(s) {
			if *budget <= 0 {
				break
			}
			*budget--
			c := s[i]
			if c == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if c == closeQ {
				i++
				titleOK = true
				break
			}
			if c == '\n' && i+1 < len(s) && s[i+1] == '\n' {
				break
			}
			if c == '\n' && isBlockInterrupt(s, i+1) {
				break
			}
			i++
		}
		if titleOK {
			end = i
		} else {
			// Title not closed cleanly — drop it from the span, keep URL
			i = titleStart
			end = titleStart
			_ = sawNL
		}
	}
	// Consume trailing newline so the next line starts fresh
	for end < len(s) && (s[end] == ' ' || s[end] == '\t') {
		end++
	}
	if end < len(s) && s[end] == '\n' {
		end++
	}
	return end, true
}

// isBlockInterrupt reports whether the line beginning at byte index
// `lineStart` opens a CommonMark block-level construct that interrupts
// a paragraph (and thus also interrupts an LRD title continuation).
// Recognized markers per CM §4.5 / §4.8 / §5.1: ATX heading (`#`),
// thematic break (`---`, `***`, `___`), blockquote (`>`), list marker
// (`-`, `*`, `+`, `1.`–`9.`), fenced code (` ``` ` or `~~~`).
// Does NOT include setext underline — that's tied to the preceding
// line and doesn't apply within an LRD title scan.
func isBlockInterrupt(s string, lineStart int) bool {
	i := lineStart
	// Up to 3 leading spaces allowed
	for j := 0; j < 3 && i < len(s) && s[i] == ' '; j++ {
		i++
	}
	if i >= len(s) {
		return false
	}
	c := s[i]
	switch c {
	case '#':
		// ATX heading: `#` then space/tab/newline/EOF
		j := i
		for j < len(s) && j < i+6 && s[j] == '#' {
			j++
		}
		if j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n') {
			return true
		}
		return j >= len(s)
	case '>':
		return true
	case '-', '*', '_':
		// Thematic break: three or more of the same char, optionally
		// separated by spaces/tabs, to end of line.
		j, count := i, 0
		for j < len(s) && s[j] != '\n' {
			if s[j] == c {
				count++
				j++
				continue
			}
			if s[j] == ' ' || s[j] == '\t' {
				j++
				continue
			}
			break
		}
		if count >= 3 && (j >= len(s) || s[j] == '\n') {
			return true
		}
		// List marker (`-` or `*`): followed by space/tab.
		if (c == '-' || c == '*') && i+1 < len(s) && (s[i+1] == ' ' || s[i+1] == '\t') {
			return true
		}
		return false
	case '+':
		if i+1 < len(s) && (s[i+1] == ' ' || s[i+1] == '\t') {
			return true
		}
		return false
	case '`', '~':
		// Fenced code: three or more of the same fence char. Mirror
		// CM §4.5: backtick fences with a backtick in the info string
		// do NOT open (goldmark rejects them). Without this check, the
		// bracket walker's LRD-title scan would terminate on a line
		// goldmark treats as paragraph text — a parser-state mismatch
		// of the same shape that drove the escapeLineLeader fix.
		j := i
		for j < len(s) && s[j] == c {
			j++
		}
		if j-i < 3 {
			return false
		}
		if c == '`' && hasBacktickBeforeNewline(s, j) {
			return false
		}
		return true
	}
	// Ordered list: 1–9 digits then `.` or `)` then space/tab.
	if c >= '0' && c <= '9' {
		j := i
		for j < len(s) && s[j] >= '0' && s[j] <= '9' && j-i < 9 {
			j++
		}
		if j < len(s) && (s[j] == '.' || s[j] == ')') && j+1 < len(s) && (s[j+1] == ' ' || s[j+1] == '\t' || s[j+1] == '\n') {
			return true
		}
	}
	return false
}

// rewriteWithSpans is pass 2 of the bracket walker. It walks bytes, encodes
// entity references in link/image destinations, leaves fence spans verbatim,
// deletes LRD spans, and prepends `\` to unescaped `[` / `]` outside spans.
func rewriteWithSpans(s string, spans []bracketSpan) string {
	if len(spans) == 0 && !strings.ContainsAny(s, "[]") {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 16)
	spanIdx := 0
	trailingBackslashes := 0
	i := 0
	for i < len(s) {
		// Skip over a span if we're entering one
		if spanIdx < len(spans) && i == spans[spanIdx].start {
			sp := spans[spanIdx]
			spanIdx++
			if sp.kind == spanLRD {
				// Delete entirely
				i = sp.end
				trailingBackslashes = 0
				continue
			}
			if sp.kind == spanLink {
				b.WriteString(s[sp.start:sp.urlStart])
				b.WriteString(percentEncodeEntityReferences(s[sp.urlStart:sp.urlEnd]))
				b.WriteString(s[sp.urlEnd:sp.end])
			} else {
				b.WriteString(s[sp.start:sp.end])
			}
			i = sp.end
			trailingBackslashes = 0
			continue
		}
		c := s[i]
		if c == '[' || c == ']' {
			if trailingBackslashes%2 == 0 {
				// Unescaped — prepend `\`
				b.WriteByte('\\')
			}
		}
		b.WriteByte(c)
		if c == '\\' {
			trailingBackslashes++
		} else {
			trailingBackslashes = 0
		}
		i++
	}
	return b.String()
}

// escapeBracketsOutsideLinks runs pass 1 + pass 2 over the input.
func escapeBracketsOutsideLinks(s string) string {
	spans := findBracketSpans(s)
	return rewriteWithSpans(s, spans)
}
