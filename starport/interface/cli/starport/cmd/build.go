package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/services/chain"
)

// NewBuild returns a new build command to build a blockchain app.
func NewBuild() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Builds an app and installs binaries",
		Args:  cobra.ExactArgs(0),
		RunE:  buildHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func buildHandler(cmd *cobra.Command, args []string) error {
	path, err := gomodulepath.Parse(getModule(appPath))
	if err != nil {
		return err
	}
	app := chain.App{
		Name:       path.Root,
		Path:       appPath,
		ImportPath: path.RawPath,
	}

	c, err := chain.New(app, logLevel(cmd), chain.WithKeyringBackend(chaincmd.KeyringBackendTest))
	if err != nil {
		return err
	}
	return c.Build(cmd.Context())
}
