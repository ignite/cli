package chaincmdrunner

import (
	"context"
	"os"

	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
)

// Simulation run the chain simulation.
func (r Runner) Simulation(
	ctx context.Context,
	appPath, simName string,
	enabled bool,
	config simulation.Config,
	genesisTime int64,
) error {
	return r.run(ctx, runOptions{stdout: os.Stdout},
		chaincmd.SimulationCommand(
			appPath,
			simName,
			chaincmd.SimappWithGenesis(config.GenesisFile),
			chaincmd.SimappWithParams(config.ParamsFile),
			chaincmd.SimappWithExportParamsPath(config.ExportParamsPath),
			chaincmd.SimappWithExportParamsHeight(config.ExportParamsHeight),
			chaincmd.SimappWithExportStatePath(config.ExportStatePath),
			chaincmd.SimappWithExportStatsPath(config.ExportStatsPath),
			chaincmd.SimappWithSeed(config.Seed),
			chaincmd.SimappWithInitialBlockHeight(config.InitialBlockHeight),
			chaincmd.SimappWithNumBlocks(config.NumBlocks),
			chaincmd.SimappWithBlockSize(config.BlockSize),
			chaincmd.SimappWithLean(config.Lean),
			chaincmd.SimappWithCommit(config.Commit),
			chaincmd.SimappWithEnable(enabled),
			chaincmd.SimappWithGenesisTime(genesisTime),
		))
}
