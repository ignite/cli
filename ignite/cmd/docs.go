package ignitecmd

import (
	"github.com/ignite/cli/v29/docs"
	"github.com/ignite/cli/v29/ignite/pkg/localfs"
	"github.com/ignite/cli/v29/ignite/pkg/markdownviewer"
	"github.com/spf13/cobra"
)

func NewDocs() *cobra.Command {
	c := &cobra.Command{
		Use:   "docs",
		Short: "Show Ignite CLI docs",
		Args:  cobra.NoArgs,
		RunE:  docsHandler,
	}
	return c
}

func docsHandler(*cobra.Command, []string) error {
	path, cleanup, err := localfs.SaveTemp(docs.Docs)
	if err != nil {
		return err
	}
	defer cleanup()

	return markdownviewer.View(path)
}
