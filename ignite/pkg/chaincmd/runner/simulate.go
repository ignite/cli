package chaincmdrunner

import (
	"context"
	"os"

	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
)

// Simulation run the chain simulation.
func (r Runner) Simulation(
	ctx context.Context,
	appPath string,
	enabled bool,
	verbose bool,
	config simulation.Config,
	period uint,
	genesisTime int64,
) error {
	return r.run(ctx, runOptions{stdout: os.Stdout},
		chaincmd.SimulationCommand(
			appPath,
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
			chaincmd.SimappWithSimulateEveryOperation(config.OnOperation),
			chaincmd.SimappWithPrintAllInvariants(config.AllInvariants),
			chaincmd.SimappWithEnable(enabled),
			chaincmd.SimappWithVerbose(verbose),
			chaincmd.SimappWithPeriod(period),
			chaincmd.SimappWithGenesisTime(genesisTime),
		))
}
