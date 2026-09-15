package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/services/gno"
)

// flag names used across gno commands.
const (
	flagGnoRemote     = "remote"
	flagGnoFrom       = "from"
	flagGnoChainID    = "chain-id"
	flagGnoGasWanted  = "gas-wanted"
	flagGnoGasFee     = "gas-fee"
	flagGnoPassphrase = "passphrase"
	flagGnoHome       = "home"
	flagGnoPkgPath    = "pkg-path"
	flagGnoMaxDeposit = "max-deposit"
	flagGnoMaxGas     = "max-gas"
	flagGnoRecover    = "recover"
	flagGnoOutput     = "output"
)

// gnoDefaultRemote is the default chain RPC address (local dev chain).
const gnoDefaultRemote = "127.0.0.1:26657"

// gnoTxBaseFlags registers the shared tx flags on a command.
func gnoTxBaseFlags(c *cobra.Command) {
	c.Flags().String(flagGnoFrom, "", "key name or bech32 address signing the tx (default: dev account test1)")
	c.Flags().String(flagGnoRemote, gnoDefaultRemote, "chain RPC address")
	c.Flags().String(flagGnoChainID, "dev", "chain id")
	c.Flags().Int64(flagGnoGasWanted, gno.DefaultGasWanted, "gas requested for the tx")
	c.Flags().String(flagGnoGasFee, "", "gas payment fee (default: 1000000ugnot)")
	c.Flags().String(flagGnoPassphrase, "", "passphrase to unlock the signing key (empty for dev keys)")
}

// gnoTxBaseFrom reads the shared tx flags into a gno.TxBaseOptions.
func gnoTxBaseFrom(cmd *cobra.Command) (tb gno.TxBaseOptions) {
	gasWanted, _ := cmd.Flags().GetInt64(flagGnoGasWanted)
	return gno.TxBaseOptions{
		From:       flagGetGnoFrom(cmd),
		Remote:     flagGetGnoRemote(cmd),
		ChainID:    flagGetGnoChainID(cmd),
		GasWanted:  gasWanted,
		GasFee:     flagGetGnoGasFee(cmd),
		Passphrase: flagGetGnoPassphrase(cmd),
	}
}

func flagGetGnoFrom(cmd *cobra.Command) string {
	from, _ := cmd.Flags().GetString(flagGnoFrom)
	return from
}

func flagGetGnoRemote(cmd *cobra.Command) string {
	remote, _ := cmd.Flags().GetString(flagGnoRemote)
	return remote
}

func flagGetGnoChainID(cmd *cobra.Command) string {
	chainID, _ := cmd.Flags().GetString(flagGnoChainID)
	return chainID
}

func flagGetGnoGasFee(cmd *cobra.Command) string {
	gasFee, _ := cmd.Flags().GetString(flagGnoGasFee)
	return gasFee
}

func flagGetGnoPassphrase(cmd *cobra.Command) string {
	passphrase, _ := cmd.Flags().GetString(flagGnoPassphrase)
	return passphrase
}

func flagGetGnoMaxGas(cmd *cobra.Command) int64 {
	maxGas, _ := cmd.Flags().GetInt64(flagGnoMaxGas)
	return maxGas
}
