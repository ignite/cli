package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	flagGenerateOnly = "generate-only"

	gasFlagAuto    = "auto"
	flagGasPrices  = "gas-prices"
	flagAdjustment = "gas-adjustment"
	flagGas        = "gas"
	flagFees       = "fees"
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
	c.PersistentFlags().String(flagFees, "", "fees to pay along with transaction; eg: 10uatom")

	c.AddCommand(NewNodeTxBank())

	return c
}

func flagSetGenerateOnly() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Bool(flagGenerateOnly, false, "build an unsigned transaction and write it to STDOUT")
	return fs
}

func getGenerateOnly(cmd *cobra.Command) bool {
	generateOnly, _ := cmd.Flags().GetBool(flagGenerateOnly)
	return generateOnly
}

func flagSetGasFlags() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagGasPrices, "", "gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)")
	fs.String(flagGas, gasFlagAuto, fmt.Sprintf("gas limit to set per-transaction; set to %q to calculate sufficient gas automatically", gasFlagAuto))
	fs.Float64(flagAdjustment, 0, "gas adjustment to set per-transaction")
	return fs
}

func getGasPrices(cmd *cobra.Command) string {
	gasPrices, _ := cmd.Flags().GetString(flagGasPrices)
	return gasPrices
}

func getGasAdjustment(cmd *cobra.Command) float64 {
	gasAdjustment, _ := cmd.Flags().GetFloat64(flagAdjustment)
	return gasAdjustment
}

func getGas(cmd *cobra.Command) string {
	gas, _ := cmd.Flags().GetString(flagGas)
	return gas
}

func getFees(cmd *cobra.Command) string {
	fees, _ := cmd.Flags().GetString(flagFees)
	return fees
}
