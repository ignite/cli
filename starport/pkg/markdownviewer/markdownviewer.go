package markdownviewer

import (
	"os"

	"github.com/charmbracelet/glow/ui"
	"golang.org/x/crypto/ssh/terminal"
)

// View starts the Markdown viewer at path that .md files are located at.
func View(path string) error {
	conf, err := config(path)
	if err != nil {
		return err
	}

	p := ui.NewProgram(conf)

	p.EnterAltScreen()
	defer p.ExitAltScreen()

	p.EnableMouseCellMotion()
	defer p.DisableMouseCellMotion()

	return p.Start()
}

func config(path string) (ui.Config, error) {
	var width uint

	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
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
