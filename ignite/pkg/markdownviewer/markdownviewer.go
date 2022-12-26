package markdownviewer

import (
	"os"

	"github.com/charmbracelet/glow/ui"
	"golang.org/x/term"
)

// View starts the Markdown viewer at path that .md files are located at.
func View(path string) error {
	conf, err := config(path)
	if err != nil {
		return err
	}

	p := ui.NewProgram(conf)

	// TODO: refactor to remove deprecated function calls
	//nolint:staticcheck,nolintlint
	p.EnterAltScreen()
	//nolint:staticcheck,nolintlint
	defer p.ExitAltScreen()

	//nolint:staticcheck,nolintlint
	p.EnableMouseCellMotion()
	//nolint:staticcheck,nolintlint
	defer p.DisableMouseCellMotion()

	_, err = p.Run()
	return err
}

func config(path string) (ui.Config, error) {
	var width uint

	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return ui.Config{}, err
	}
	width = uint(w)
	if width > 120 {
		width = 120
	}

	docTypes := ui.NewDocTypeSet()
	docTypes.Add(ui.LocalDoc)

	conf := ui.Config{
		WorkingDirectory:     path,
		DocumentTypes:        docTypes,
		GlamourStyle:         "auto",
		HighPerformancePager: true,
		GlamourEnabled:       true,
		GlamourMaxWidth:      width,
	}

	return conf, nil
}
