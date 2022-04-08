package starportcmd

import (
	"github.com/ignite-hq/cli/docs"
	"github.com/ignite-hq/cli/ignite/pkg/localfs"
	"github.com/ignite-hq/cli/ignite/pkg/markdownviewer"
	"github.com/spf13/cobra"
)

func NewDocs() *cobra.Command {
	c := &cobra.Command{
		Use:   "docs",
		Short: "Show Starport docs",
		Args:  cobra.NoArgs,
		RunE:  docsHandler,
	}
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
