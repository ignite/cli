package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/doctor"
)

func NewDoctor() *cobra.Command {
	return &cobra.Command{
		Use:    "doctor",
		Short:  "Fix chain configuration",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New()
			defer session.End()

			doc := doctor.New(doctor.CollectEvents(session.EventBus()))

			if err := doc.MigrateConfig(cmd.Context()); err != nil {
				return err
			}

			return doc.FixDependencyTools(cmd.Context())
		},
	}
}
