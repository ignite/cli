package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/docs"
	"github.com/tendermint/starport/starport/pkg/offlinepage"
)

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
	path, err := offlinepage.SaveTemp(docs.MigrationDocs)
	fmt.Println(path)
	return err
}
