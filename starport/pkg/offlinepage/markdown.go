package offlinepage

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const migrationTempDir = "migration"

var (
	// exts defines the extensions that are used
	exts = blackfriday.Tables |
		blackfriday.Autolink |
		blackfriday.Footnotes |
		blackfriday.HeadingIDs |
		blackfriday.Titleblock |
		blackfriday.FencedCode |
		blackfriday.LaxHTMLBlocks |
		blackfriday.HardLineBreak |
		blackfriday.Strikethrough |
		blackfriday.SpaceHeadings |
		blackfriday.AutoHeadingIDs |
		blackfriday.DefinitionLists |
		blackfriday.NoIntraEmphasis |
		blackfriday.BackslashLineBreak

	// flags defines the HTML rendering flags that are used
	flags = blackfriday.TOC |
		blackfriday.UseXHTML |
		blackfriday.CompletePage |
		blackfriday.FootnoteReturnLinks |
		blackfriday.Smartypants |
		blackfriday.SmartypantsDashes |
		blackfriday.SmartypantsFractions |
		blackfriday.SmartypantsQuotesNBSP |
		blackfriday.SmartypantsLatexDashes |
		blackfriday.SmartypantsAngledQuotes
)

// Markdown saves file system f markdown in converted html
// to a temporary path and returns that path.
func Markdown(f fs.FS) (string, error) {
	path, err := os.MkdirTemp("", migrationTempDir)
	if err != nil {
		return path, err
	}
	return fmt.Sprintf("file://%s", path), save(f, path)
}

// save saves the markdown file converted to html to path.
func save(f fs.FS, path string) error {
	return fs.WalkDir(f, ".", func(wpath string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		content, err := fs.ReadFile(f, wpath)
		if err != nil {
			return err
		}

		name := strings.ReplaceAll(d.Name(), ".md", ".html")
		out := filepath.Join(path, name)

		htmlContent := markdownToHTML(content)
		return os.WriteFile(out, htmlContent, 0644)
	})
}

// markdownToHTML creates a html content based in the markdown path
func markdownToHTML(src []byte) []byte {
	// create the html styles
	r := bfchroma.NewRenderer(
		bfchroma.Extend(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{Flags: flags}),
		),
		bfchroma.ChromaStyle(styles.GitHub),
		bfchroma.ChromaOptions(
			html.WithLineNumbers(true),
			html.LineNumbersInTable(true),
		),
	)

	// render the markdown to html with options
	unsafe := blackfriday.Run(src, blackfriday.WithRenderer(r), blackfriday.WithExtensions(exts))
	p := bluemonday.UGCPolicy()
	result := p.SanitizeBytes(unsafe)

	return result
}
