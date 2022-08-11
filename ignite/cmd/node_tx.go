package ignitecmd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	flagGenerateOnly = "generate-only"

	gasFlagAuto       = "auto"
	flagGasPrices     = "gas-prices"
	flagGas           = "gas"
	flagFees          = "fees"
	flagBroadcastMode = "broadcast-mode"
)

func NewNodeTx() *cobra.Command {
	c := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}
	c.PersistentFlags().AddFlagSet(flagSetHome())
	c.PersistentFlags().AddFlagSet(flagSetKeyringBackend())
	c.PersistentFlags().AddFlagSet(flagSetAccountPrefixes())
	c.PersistentFlags().AddFlagSet(flagSetKeyringDir())
	c.PersistentFlags().AddFlagSet(flagSetGenerateOnly())
	c.PersistentFlags().AddFlagSet(flagSetGasFlags())
	c.PersistentFlags().String(flagFees, "", "Fees to pay along with transaction; eg: 10uatom")
	c.PersistentFlags().String(flagBroadcastMode, flags.BroadcastBlock, "Transaction broadcasting mode (sync|async|block), use sync if you encounter timeouts")

	c.AddCommand(NewNodeTxBank())

	return c
}

func flagSetGenerateOnly() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagGenerateOnly, false, "Build an unsigned transaction and write it to STDOUT")
	return fs
}

func getGenerateOnly(cmd *cobra.Command) bool {
	generateOnly, _ := cmd.Flags().GetBool(flagGenerateOnly)
	return generateOnly
}

func flagSetGasFlags() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagGasPrices, "", "Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)")
	fs.String(flagGas, gasFlagAuto, fmt.Sprintf("gas limit to set per-transaction; set to %q to calculate sufficient gas automatically", gasFlagAuto))
	return fs
}

func getGasPrices(cmd *cobra.Command) string {
	gasPrices, _ := cmd.Flags().GetString(flagGasPrices)
	return gasPrices
}

func getGas(cmd *cobra.Command) string {
	gas, _ := cmd.Flags().GetString(flagGas)
	return gas
}

func getFees(cmd *cobra.Command) string {
	fees, _ := cmd.Flags().GetString(flagFees)
	return fees
}

func getBroadcastMode(cmd *cobra.Command) string {
	broadcastMode, _ := cmd.Flags().GetString(flagBroadcastMode)
	return broadcastMode
}
