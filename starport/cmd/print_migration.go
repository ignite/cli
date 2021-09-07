package starportcmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
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

// Markdown creates a html docs server from the markdown
func Markdown(ctx context.Context, path string, port uint32) error {
	g, ctx := errgroup.WithContext(ctx)
	htmlText, err := markdownToHtml(path)
	if err != nil {
		return err
	}
	portSuffix := fmt.Sprintf(":%d", port)
	g.Go(func() error {
		return openBrowser("http://localhost" + portSuffix)
	})

	g.Go(func() error {
		http.HandleFunc("/", docHandler(htmlText))
		return http.ListenAndServe(portSuffix, nil)
	})
	return g.Wait()
}

// docHandler http method handler to serve the docs
func docHandler(htmlText string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := fmt.Fprint(w, htmlText)
		if err != nil {
			panic(err)
		}
	}
}

// markdownToHtml create a html content based in the markdown path
func markdownToHtml(path string) (string, error) {
	// read the markdown content
	src, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

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

	return string(result), nil
}

// openBrowser open the current OS browser with the giving url
func openBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

// PrintMigration print the migration guide.
func PrintMigration() *cobra.Command {
	c := &cobra.Command{
		Use:   "print-migration",
		Short: "Print the current migration guide",
		Run: func(cmd *cobra.Command, _ []string) {
			path, err := filepath.Abs("docs/guide/install.md")
			if err != nil {
				panic(err)
			}
			port := uint32(8080)
			err = Markdown(context.Background(), path, port)
			if err != nil {
				panic(err)
			}
		},
	}
	return c
}
