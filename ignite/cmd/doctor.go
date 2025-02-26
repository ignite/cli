package ignitecmd

import (
	"github.com/spf13/cobra"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/doctor"
)

func NewDoctor() *cobra.Command {
	c := &cobra.Command{
		Use:    "doctor",
		Short:  "Fix chain configuration",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			session := cliui.New()
			defer session.End()
			appPath := flagGetPath(cmd)

			doc := doctor.New(doctor.CollectEvents(session.EventBus()))

			cacheStorage, err := newCache(cmd)
			if err != nil {
				return err
			}

			configPath, err := chainconfig.LocateDefault(appPath)
			if err != nil {
				return err
			}

			if err := doc.MigrateChainConfig(configPath); err != nil {
				return err
			}

			if err := doc.MigrateBufConfig(cmd.Context(), cacheStorage, appPath, configPath); err != nil {
				return err
			}

			if err := doc.MigratePluginsConfig(); err != nil {
				return err
			}

			return nil
		},
	}

	flagSetPath(c)
	return c
}
