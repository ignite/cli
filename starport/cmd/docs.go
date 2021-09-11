package starportcmd

import (
	"fmt"
	"io/fs"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/docs"
	"github.com/tendermint/starport/starport/pkg/localfs"
	"github.com/tendermint/starport/starport/pkg/markdownviewer"
	"github.com/tendermint/starport/starport/pkg/offlinepage"
)

const migrationDocsPath = "guide/migration"

func NewDocs() *cobra.Command {
	c := &cobra.Command{
		Use:   "docs",
		Short: "Show Starport docs",
		Args:  cobra.ExactArgs(0),
		RunE:  docsHandler,
	}

	c.AddCommand(Migration())

	return c
}

func docsHandler(cmd *cobra.Command, args []string) error {
	path, cleanup, err := localfs.SaveTemp(docs.Docs)
	if err != nil {
		return err
	}
	defer cleanup()

	return markdownviewer.View(path)
}

// Migration print the migration guide.
func Migration() *cobra.Command {
	c := &cobra.Command{
		Use:   "migration",
		Short: "Shows the migration guide",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := markdown()
			if err != nil {
				return err
			}
			fmt.Println(path)
			return nil
		},
	}
	return c
}

func markdown() (string, error) {
	sub, err := fs.Sub(docs.Docs, migrationDocsPath)
	if err != nil {
		return "", err
	}
	return offlinepage.Markdown(sub)
}
