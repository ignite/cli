package starportcmd

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/chain"
)

// NewRelayer creates a new command called chain that holds IBC Relayer related
// sub commands.
func NewRelayer() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain",
		Short: "Relay connects blockchains via IBC protocol",
	}
	c.AddCommand(NewRelayerInfo())
	c.AddCommand(NewRelayerAdd())
	return c
}

// NewRelayerInfo creates a command that shows self chain information.
func NewRelayerInfo() *cobra.Command {
	c := &cobra.Command{
		Use:   "me",
		Short: "Retrieves self chain information to share with other chains",
		RunE:  relayerInfoHandler,
	}
	c.Flags().AddFlagSet(flagSetHomes())
	return c
}

// NewRelayerAdd creates a command to connect added chain with relayer.
func NewRelayerAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [another]",
		Short: "Adds another chain by its chain information",
		Args:  cobra.MinimumNArgs(1),
		RunE:  relayerAddHandler,
	}
	c.Flags().AddFlagSet(flagSetHomes())
	return c
}

func relayerInfoHandler(cmd *cobra.Command, args []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	// Check if custom home is provided
	home, cliHome, err := getHomeFlags(cmd)
	if err != nil {
		return err
	}
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	if cliHome != "" {
		chainOption = append(chainOption, chain.CLIHomePath(cliHome))
	}

	c, err := chain.New(appPath, chainOption...)
	if err != nil {
		return err
	}
	info, err := c.RelayerInfo()
	if err != nil {
		return err
	}
	fmt.Println(info)
	return nil
}

func relayerAddHandler(cmd *cobra.Command, args []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	// Check if custom home is provided
	home, cliHome, err := getHomeFlags(cmd)
	if err != nil {
		return err
	}
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	if cliHome != "" {
		chainOption = append(chainOption, chain.CLIHomePath(cliHome))
	}

	c, err := chain.New(appPath, chainOption...)
	if err != nil {
		return err
	}
	if err := c.RelayerAdd(args[0]); err != nil {
		return err
	}
	return nil
}
