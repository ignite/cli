package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/services/chain"
)

var appPath string

// NewServe creates a new serve command to serve a blockchain.
func NewServe() *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Launches a reloading server",
		Args:  cobra.ExactArgs(0),
		RunE:  serveHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func serveHandler(cmd *cobra.Command, args []string) error {
	path, err := gomodulepath.Parse(getModule(appPath))
	if err != nil {
		return err
	}
	app := chain.App{
		Name:       path.Root,
		Path:       appPath,
		ImportPath: path.RawPath,
	}

	c, err := chain.New(app, false, logLevel(cmd), chain.ChainWithKeyringBackend(chaincmd.KeyringBackendTest))
	if err != nil {
		return err
	}
	return c.Serve(cmd.Context())
}
