package starportcmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

func NewPluginsInstall() *cobra.Command {
	c := &cobra.Command{
		Use:   "install",
		Short: "install plugins listed in config",
		RunE:  pluginsInstallHandler,
	}

	flagSetPath(c)

	return c
}

func pluginsInstallHandler(cmd *cobra.Command, args []string) error {
	appPath := flagGetPath(cmd)
	confpath, err := chainconfig.LocateDefault(appPath)
	if err != nil {
		return err
	}
	conf, err := chainconfig.ParseFile(confpath)
	if err != nil {
		return err
	}

	pluginsPath := filepath.Join(appPath, "plugins")
	os.RemoveAll(pluginsPath)

	for _, plugin := range conf.Plugins {
		pluginPath := filepath.Join(pluginsPath, plugin.Name)
		url := "https://" + plugin.Repo

		_, err := git.PlainClone(pluginPath, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}

		ctx := context.Background()
		return cmdrunner.
			New(
				cmdrunner.DefaultWorkdir(pluginPath),
			).
			Run(ctx,
				step.New(
					step.Exec(
						"make",
						"plugin",
					),
				),
			)
	}

	return nil
}
