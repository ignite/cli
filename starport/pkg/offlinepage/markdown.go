package offlinepage

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var (
	// exts defines the extensions that are used
	exts = blackfriday.Tables |
		blackfriday.Autolink |
		blackfriday.Footnotes |
		blackfriday.HeadingIDs |
		blackfriday.FencedCode |
		blackfriday.TabSizeEight |
		blackfriday.HardLineBreak |
		blackfriday.Strikethrough |
		blackfriday.SpaceHeadings |
		blackfriday.AutoHeadingIDs |
		blackfriday.NoIntraEmphasis |
		blackfriday.DefinitionLists |
		blackfriday.BackslashLineBreak

	// flags defines the HTML rendering flags that are used
	flags = blackfriday.TOC |
		blackfriday.UseXHTML |
		blackfriday.Smartypants |
		blackfriday.CompletePage |
		blackfriday.SmartypantsDashes |
		blackfriday.SmartypantsFractions |
		blackfriday.SmartypantsQuotesNBSP |
		blackfriday.SmartypantsLatexDashes |
		blackfriday.SmartypantsAngledQuotes
)

// SaveTemp saves file system f markdown in converted html to a temporary path
// and returns that path.
func SaveTemp(f fs.FS) (string, error) {
	path, err := os.MkdirTemp("", "")
	if err != nil {
		return path, err
	}
	return path, save(f, path)
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

		htmlContent, err := markdown(content)
		if err != nil {
			return err
		}
		name := strings.ReplaceAll(d.Name(), ".md", ".html")
		out := filepath.Join(path, name)
		return os.WriteFile(out, htmlContent, 0644)
	})
}

// markdown creates a html content based in the markdown path
func markdown(src []byte) ([]byte, error) {
	// remove file description by regex
	noTileRegex := regexp.MustCompile(`---[\s\S]*?---`)
	noTitle := noTileRegex.ReplaceAllString(string(src), "")

	// create the html styles
	r := bfchroma.NewRenderer(
		bfchroma.EmbedCSS(),
		bfchroma.WithoutAutodetect(),
		bfchroma.Extend(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{Flags: flags}),
		),
		bfchroma.ChromaStyle(styles.Monokai),
		bfchroma.ChromaOptions(
			html.WithLineNumbers(true),
			html.WithAllClasses(true),
			html.WithClasses(true),
		),
	)

	// render the markdown to html with options
	unsafe := blackfriday.Run([]byte(noTitle), blackfriday.WithRenderer(r), blackfriday.WithExtensions(exts))
	p := bluemonday.UGCPolicy()
	result := p.SanitizeBytes(unsafe)

	return result, nil
}
