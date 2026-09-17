// Package markdown provides byte-level primitives for sanitizing
// user-supplied content before it flows into markdown rendering.
//
// Two functions are first-class public APIs (callable directly from
// realms for non-markdown data hygiene); six are documented as
// internal building blocks — realm authors should use the matching
// helpers in p/nt/md/sanitize, which combine these natives with
// the policy layer (scheme allowlists, charset bitmaps, length
// bounds, NUL substitution).
//
// All functions are pure (string → string or similar), single-pass,
// O(N) on input length, no compiled regex, no backtracking.
package markdown

// StripBidiAndZeroWidth removes Unicode bidi controls and zero-width
// characters: U+200B-D, U+200E-F, U+202A-E, U+2066-9, U+FEFF.
// Single rune-walk; output allocated to len(s) worst-case.
//
// PUBLIC. Useful beyond markdown: cleaning usernames, search queries,
// any user string that will be displayed or compared.
func StripBidiAndZeroWidth(s string) string // injected

// NormalizeBreaks unifies \r\n and \r to \n. Single byte scan.
//
// PUBLIC. Useful beyond markdown: normalizing user input before
// storage, comparison, or hashing. Does NOT fold U+2028/U+2029
// (those are Unicode separators, not CM §2.2 line endings).
func NormalizeBreaks(s string) string // injected

// EscapeInline backslash-escapes the inline-active CM punctuation set
// with \ always first. Replaces NUL with U+FFFD per CM §2.3.
//
// Set: \ * _ [ ] ( ) ~ > - + . ! ` # < &
// Intentionally excludes | (only table-active in GFM) and = (only
// block-active for setext underline).
//
// INTERNAL. Used by p/nt/md/sanitize.InlineText.
func EscapeInline(s string) string // injected

// EscapeTitle backslash-escapes the inline set plus " and ' (the
// remaining link-title delimiters; ( and ) are already in the inline
// set). \ always first. Replaces NUL with U+FFFD.
//
// INTERNAL. Used by p/nt/md/sanitize.LinkTitle.
func EscapeTitle(s string) string // injected

// PercentEncodeURL percent-encodes bytes unsafe in CM link destinations.
//
// Bytes encoded:
//   - 0x00-0x20, 0x7F (controls + space)
//   - " ' ( ) < > \ ` { | } ^
//   - 0x80-0xFF (UTF-8 high bytes — each byte encoded)
//   - Bare % not followed by two hex digits
//   - & when it opens an HTML entity reference (`&#…;` or `&name;`,
//     including the `&\#…;` / `&…\;` forms a renderer unescapes before
//     resolving). A renderer resolves those to raw bytes when it writes the
//     URL out, so leaving `&#x6a;avascript:` intact would hand `javascript:`
//     to any renderer that checks the scheme without resolving first.
//
// Preserved (URL-syntactic):
//   - A-Z a-z 0-9 - _ . ~ / : ? # [ ] @ ! $ + , ; = *
//   - & anywhere else, so query separators like `?x=1&y=2` survive
//
// The shape rule is deliberately blind to what a reference expands to, so
// harmless ones are encoded too: pass `?a=1&amp;b=2` and you get
// `?a=1%26amp;b=2`, which is ONE parameter valued `1&amp;b=2`, not two.
// Write the bare `&` (`?a=1&b=2`) — that is what CM link destinations take,
// and it round-trips. Exempting `&amp;` would be sound against this exact
// resolver chain (its output `&` is produced by the last pass, so nothing
// re-reads it) but not against a renderer that resolves to a fixed point,
// which would turn `&amp;#x6a;` back into `j`.
//
// INTERNAL. Does not apply scheme allowlist. Use the URL or ImageURL
// helpers in p/nt/md/sanitize for full filtering.
func PercentEncodeURL(s string) string // injected

// MatchCharsetN runs s against a DFA defined by 128-bit bitmaps for
// the first-char class and the rest-char class, plus length bounds
// (in bytes).
//
// Bitmap encoding: each pair (lo, hi) covers ASCII codepoints, with
// lo bit i for codepoint i (0-63) and hi bit i for codepoint i+64
// (64-127). Bytes >= 128 are rejected outright.
//
// INTERNAL. Used by p/nt/markdown/sanitize/v0.UserName, .FootnoteLabel,
// .LanguageName, .BechString matchers.
func MatchCharsetN(s string, firstLo, firstHi, restLo, restHi uint64, minLen, maxLen int) bool // injected

// CodeFence returns a string of N backticks where
// N = max(minCount, longestBacktickRunInContent + 1).
//
// minCount < 1 is clamped to 1; never panics. On empty content,
// returns strings.Repeat("`", max(minCount, 1)).
//
// INTERNAL. Used by p/nt/md/sanitize.CodeFence and moul/md
// InlineCode / CodeBlock / LanguageCodeBlock bug fixes.
func CodeFence(content string, minCount int) string // injected

// EscapeBlockHazards scans s line by line and applies block-context
// escaping in a single pass. Steps:
//  1. Strip CM link reference definitions.
//  2. Escape line-leading block markers (ATX heading, thematic break,
//     blockquote, list marker, code fence) by prefixing with \.
//  3. Escape setext heading underlines following a non-blank line.
//  4. Escape extension-delimiter lines (any line starting with
//     `<gno-…>` or `</gno-…>`) by prefixing with a backslash so the
//     opening `<` becomes a CM §2.4 inline escape.
//  5. Reference-link USE escape: [text][label] → [text]\[label].
//  6. Fold U+2028, U+2029, U+0085 NEL to \n.
//  7. Auto-close any open code fence at end-of-input.
//
// Preserves inline emphasis, code spans, inline links, GFM tables,
// soft line breaks within paragraphs.
//
// INTERNAL. Used by p/nt/md/sanitize.Block. Calling directly without
// first running StripBidiAndZeroWidth and NormalizeBreaks leaves
// some hazards uncovered.
func EscapeBlockHazards(s string) string // injected

// EscapeBlockHazardsRich is the permissive variant of
// EscapeBlockHazards. All realm-binding defenses stay on (bracket
// walker for link / LRD / ref / footnote / shortcut, extension
// delimiter `<gno-…>` escape, GFM table-row pipe escape, fence
// state, NUL / bidi / U+2028 fold), but doc-spoof escapes are
// skipped — ATX headings (#), blockquote markers (>), list
// markers (-, *, +, 1.), thematic breaks (---, ***, ___), and
// setext underlines (===, ---) are preserved.
//
// Cross-boundary setext promotion (user content starting with
// `===`/`---` reaching back to promote a realm chrome line) is
// NOT defended at this layer — the Gno-side p/nt/md/sanitize.BlockRich
// wrapper applies a separate qualifying-line escape before calling
// this native.
//
// INTERNAL. Used by p/nt/md/sanitize.BlockRich.
func EscapeBlockHazardsRich(s string) string // injected

// MaxForeignBlocksPerConvert is gnoweb's per-render cap on the number
// of <gno-foreign> sandbox blocks a single markdown render will admit.
// Beyond it, later openers fall through to raw HTML and are dropped by
// safe mode. A realm that emits many foreign blocks (e.g. one per
// comment) should keep its rendered total under this so content is not
// silently dropped. This is the single source of truth: the gnoweb
// foreign renderer enforces the same value.
func MaxForeignBlocksPerConvert() int // injected
