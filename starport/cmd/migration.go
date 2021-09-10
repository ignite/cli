package starportcmd

import (
	"io/fs"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/docs"
	"github.com/tendermint/starport/starport/pkg/offlinepage"
)

const migrationDocsPath = "guide/migration"

// NewMigration print the migration guide.
func NewMigration() *cobra.Command {
	c := &cobra.Command{
		Use:   "show-migration",
		Short: "Shows the current migration guide",
		Args:  cobra.ExactArgs(0),
		RunE:  migrationHandler,
	}
	return c
}

func migrationHandler(cmd *cobra.Command, args []string) error {
	sub, err := fs.Sub(docs.Docs, migrationDocsPath)
	if err != nil {
		return err
	}
	path, err := offlinepage.SaveTemp(sub)
	if err != nil {
		return err
	}
	return open.Run(path)
}
