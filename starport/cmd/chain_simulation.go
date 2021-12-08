package starportcmd

import (
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/spf13/cobra"
)

const (
	flagSimappGenesis          = "Genesis"
	flagSimappParams           = "Params"
	flagExportParamsPath       = "ExportParamsPath"
	flagExportParamsHeight     = "ExportParamsHeight"
	flagExportStatePath        = "ExportStatePath"
	flagExportStatsPath        = "ExportStatsPath"
	flagSeed                   = "Seed"
	flagInitialBlockHeight     = "InitialBlockHeight"
	flagNumBlocks              = "NumBlocks"
	flagBlockSize              = "BlockSize"
	flagLean                   = "Lean"
	flagCommit                 = "Commit"
	flagSimulateEveryOperation = "SimulateEveryOperation"
	flagPrintAllInvariants     = "PrintAllInvariants"

	flagEnabled     = "Enabled"
	flagVerbose     = "Verbose"
	flagPeriod      = "Period"
	flagGenesisTime = "GenesisTime"
)

// NewChainSimulation creates a new simulation command to run the blockchain simulation.
func NewChainSimulation() *cobra.Command {
	c := &cobra.Command{
		Use:   "simulation",
		Short: "Run the blockchain simulation node in development",
		Long:  "Run the blockchain simulation for all chain modules",
		Args:  cobra.ExactArgs(0),
		RunE:  chainSimulationHandler,
	}
	simappFlags(c)
	flagSetPath(c)
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func chainSimulationHandler(cmd *cobra.Command, args []string) error {
	var (
		// simulation flags
		enabled, _     = cmd.Flags().GetBool("Enabled")
		verbose, _     = cmd.Flags().GetBool("Verbose")
		period, _      = cmd.Flags().GetUint("Period")
		genesisTime, _ = cmd.Flags().GetInt64("GenesisTime")
		config         = newConfigFromFlags(cmd)
	)

	// create the chain
	c, err := newChainWithHomeFlags(cmd)
	if err != nil {
		return err
	}

	return c.Simulate(cmd.Context())
}

// newConfigFromFlags creates a simulation from the retrieved values of the flags.
func newConfigFromFlags(cmd *cobra.Command) simulation.Config {
	var (
		genesis, _                = cmd.Flags().GetString(flagGenesis)
		params, _                 = cmd.Flags().GetString(flagParams)
		exportParamsPath, _       = cmd.Flags().GetString(flagExportParamsPath)
		exportParamsHeight, _     = cmd.Flags().GetInt(flagExportParamsHeight)
		exportStatePath, _        = cmd.Flags().GetString(flagExportStatePath)
		exportStatsPath, _        = cmd.Flags().GetString(flagExportStatsPath)
		seed, _                   = cmd.Flags().GetInt64(flagSeed)
		initialBlockHeight, _     = cmd.Flags().GetInt(flagInitialBlockHeight)
		numBlocks, _              = cmd.Flags().GetInt(flagNumBlocks)
		blockSize, _              = cmd.Flags().GetInt(flagBlockSize)
		lean, _                   = cmd.Flags().GetBool(flagLean)
		commit, _                 = cmd.Flags().GetBool(flagCommit)
		simulateEveryOperation, _ = cmd.Flags().GetBool(flagSimulateEveryOperation)
		printAllInvariants, _     = cmd.Flags().GetBool(flagPrintAllInvariants)
	)
	return simulation.Config{
		GenesisFile:        genesis,
		ParamsFile:         params,
		ExportParamsPath:   exportParamsPath,
		ExportParamsHeight: exportParamsHeight,
		ExportStatePath:    exportStatePath,
		ExportStatsPath:    exportStatsPath,
		Seed:               seed,
		InitialBlockHeight: initialBlockHeight,
		NumBlocks:          numBlocks,
		BlockSize:          blockSize,
		Lean:               lean,
		Commit:             commit,
		OnOperation:        simulateEveryOperation,
		AllInvariants:      printAllInvariants,
	}
}

func simappFlags(c *cobra.Command) {
	// config fields
	c.Flags().String(flagGenesis, "", "custom simulation genesis file; cannot be used with params file")
	c.Flags().String(flagParams, "", "custom simulation params file which overrides any random params; cannot be used with genesis")
	c.Flags().String(flagExportParamsPath, "", "custom file path to save the exported params JSON")
	c.Flags().Int(flagExportParamsHeight, 0, "height to which export the randomly generated params")
	c.Flags().String(flagExportStatePath, "", "custom file path to save the exported app state JSON")
	c.Flags().String(flagExportStatsPath, "", "custom file path to save the exported simulation statistics JSON")
	c.Flags().Int64(flagSeed, 42, "simulation random seed")
	c.Flags().Int(flagInitialBlockHeight, 1, "initial block to start the simulation")
	c.Flags().Int(flagNumBlocks, 500, "number of new blocks to simulate from the initial block height")
	c.Flags().Int(flagBlockSize, 200, "operations per block")
	c.Flags().Bool(flagLean, false, "lean simulation log output")
	c.Flags().Bool(flagCommit, false, "have the simulation commit")
	c.Flags().Bool(flagSimulateEveryOperation, false, "run slow invariants every operation")
	c.Flags().Bool(flagPrintAllInvariants, false, "print all invariants if a broken invariant is found")

	// simulation flags
	c.Flags().Bool(flagEnabled, false, "enable the simulation")
	c.Flags().Bool(flagVerbose, false, "verbose log output")
	c.Flags().Uint(flagPeriod, 0, "run slow invariants only once every period assertions")
	c.Flags().Int64(flagGenesisTime, 0, "override genesis UNIX time instead of using a random UNIX time")
}
