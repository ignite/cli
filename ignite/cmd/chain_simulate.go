package ignitecmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/ignite/cli/v29/ignite/services/chain"
)

const (
	flagSimappGenesis            = "genesis"
	flagSimappParams             = "params"
	flagSimappExportParamsPath   = "exportParamsPath"
	flagSimappExportParamsHeight = "exportParamsHeight"
	flagSimappExportStatePath    = "exportStatePath"
	flagSimappExportStatsPath    = "exportStatsPath"
	flagSimappSeed               = "seed"
	flagSimappInitialBlockHeight = "initialBlockHeight"
	flagSimappNumBlocks          = "numBlocks"
	flagSimappBlockSize          = "blockSize"
	flagSimappLean               = "lean"
	flagSimappGenesisTime        = "genesisTime"
	flagSimName                  = "simName"
)

// NewChainSimulate creates a new simulation command to run the blockchain simulation.
func NewChainSimulate() *cobra.Command {
	c := &cobra.Command{
		Use:   "simulate",
		Short: "Run simulation testing for the blockchain",
		Long:  "Run simulation testing for the blockchain. It sends many randomized-input messages of each module to a simulated node.",
		Args:  cobra.NoArgs,
		RunE:  chainSimulationHandler,
	}
	simappFlags(c)
	return c
}

func chainSimulationHandler(cmd *cobra.Command, _ []string) error {
	var (
		genesisTime, _ = cmd.Flags().GetInt64(flagSimappGenesisTime)
		simName, _     = cmd.Flags().GetString(flagSimName)
		config         = newConfigFromFlags(cmd)
		appPath        = flagGetPath(cmd)
	)
	// create the chain with path
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return err
	}
	c, err := chain.New(absPath)
	if err != nil {
		return err
	}

	config.ChainID, err = c.ID()
	if err != nil {
		return err
	}

	return c.Simulate(cmd.Context(),
		chain.SimappWithSimulationTestName(simName),
		chain.SimappWithGenesisTime(genesisTime),
		chain.SimappWithConfig(config),
	)
}

// newConfigFromFlags creates a simulation from the retrieved values of the flags.
func newConfigFromFlags(cmd *cobra.Command) simulation.Config {
	var (
		genesis, _            = cmd.Flags().GetString(flagSimappGenesis)
		params, _             = cmd.Flags().GetString(flagSimappParams)
		exportParamsPath, _   = cmd.Flags().GetString(flagSimappExportParamsPath)
		exportParamsHeight, _ = cmd.Flags().GetInt(flagSimappExportParamsHeight)
		exportStatePath, _    = cmd.Flags().GetString(flagSimappExportStatePath)
		exportStatsPath, _    = cmd.Flags().GetString(flagSimappExportStatsPath)
		seed, _               = cmd.Flags().GetInt64(flagSimappSeed)
		initialBlockHeight, _ = cmd.Flags().GetInt(flagSimappInitialBlockHeight)
		numBlocks, _          = cmd.Flags().GetInt(flagSimappNumBlocks)
		blockSize, _          = cmd.Flags().GetInt(flagSimappBlockSize)
		lean, _               = cmd.Flags().GetBool(flagSimappLean)
	)
	return simulation.Config{
		Commit:             true,
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
	}
}

func simappFlags(c *cobra.Command) {
	// config fields
	c.Flags().String(flagSimappGenesis, "", "custom simulation genesis file; cannot be used with params file")
	c.Flags().String(flagSimappParams, "", "custom simulation params file which overrides any random params; cannot be used with genesis")
	c.Flags().String(flagSimappExportParamsPath, "", "custom file path to save the exported params JSON")
	c.Flags().Int(flagSimappExportParamsHeight, 0, "height to which export the randomly generated params")
	c.Flags().String(flagSimappExportStatePath, "", "custom file path to save the exported app state JSON")
	c.Flags().String(flagSimappExportStatsPath, "", "custom file path to save the exported simulation statistics JSON")
	c.Flags().Int64(flagSimappSeed, 42, "simulation random seed")
	c.Flags().Int(flagSimappInitialBlockHeight, 1, "initial block to start the simulation")
	c.Flags().Int(flagSimappNumBlocks, 200, "number of new blocks to simulate from the initial block height")
	c.Flags().Int(flagSimappBlockSize, 30, "operations per block")
	c.Flags().Bool(flagSimappLean, false, "lean simulation log output")

	// simulation flags
	c.Flags().String(flagSimName, "TestFullAppSimulation", "name of the simulation to run")
	c.Flags().Int64(flagSimappGenesisTime, 0, "override genesis UNIX time instead of using a random UNIX time")
}
