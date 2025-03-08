package chaincmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
)

const (
	optionSimappGenesis            = "-Genesis"
	optionSimappParams             = "-Params"
	optionSimappExportParamsPath   = "-ExportParamsPath"
	optionSimappExportParamsHeight = "-ExportParamsHeight"
	optionSimappExportStatePath    = "-ExportStatePath"
	optionSimappExportStatsPath    = "-ExportStatsPath"
	optionSimappSeed               = "-Seed"
	optionSimappInitialBlockHeight = "-InitialBlockHeight"
	optionSimappNumBlocks          = "-NumBlocks"
	optionSimappBlockSize          = "-BlockSize"
	optionSimappLean               = "-Lean"
	optionSimappCommit             = "-Commit"
	optionSimappEnabled            = "-Enabled"
	optionSimappGenesisTime        = "-GenesisTime"

	commandGoTest    = "test"
	optionGoBenchmem = "-benchmem"
	optionGoSimsTags = "-tags='sims'"
)

// SimappOption for the SimulateCommand.
type SimappOption func([]string) []string

// SimappWithGenesis provides genesis option for the simapp command.
func SimappWithGenesis(genesis string) SimappOption {
	return func(command []string) []string {
		if len(genesis) > 0 {
			return append(command, optionSimappGenesis, genesis)
		}
		return command
	}
}

// SimappWithParams provides params option for the simapp command.
func SimappWithParams(params string) SimappOption {
	return func(command []string) []string {
		if len(params) > 0 {
			return append(command, optionSimappParams, params)
		}
		return command
	}
}

// SimappWithExportParamsPath provides exportParamsPath option for the simapp command.
func SimappWithExportParamsPath(exportParamsPath string) SimappOption {
	return func(command []string) []string {
		if len(exportParamsPath) > 0 {
			return append(command, optionSimappExportParamsPath, exportParamsPath)
		}
		return command
	}
}

// SimappWithExportParamsHeight provides exportParamsHeight option for the simapp command.
func SimappWithExportParamsHeight(exportParamsHeight int) SimappOption {
	return func(command []string) []string {
		if exportParamsHeight > 0 {
			return append(
				command,
				optionSimappExportParamsHeight,
				strconv.Itoa(exportParamsHeight),
			)
		}
		return command
	}
}

// SimappWithExportStatePath provides exportStatePath option for the simapp command.
func SimappWithExportStatePath(exportStatePath string) SimappOption {
	return func(command []string) []string {
		if len(exportStatePath) > 0 {
			return append(command, optionSimappExportStatePath, exportStatePath)
		}
		return command
	}
}

// SimappWithExportStatsPath provides exportStatsPath option for the simapp command.
func SimappWithExportStatsPath(exportStatsPath string) SimappOption {
	return func(command []string) []string {
		if len(exportStatsPath) > 0 {
			return append(command, optionSimappExportStatsPath, exportStatsPath)
		}
		return command
	}
}

// SimappWithSeed provides seed option for the simapp command.
func SimappWithSeed(seed int64) SimappOption {
	return func(command []string) []string {
		return append(command, optionSimappSeed, strconv.FormatInt(seed, 10))
	}
}

// SimappWithInitialBlockHeight provides initialBlockHeight option for the simapp command.
func SimappWithInitialBlockHeight(initialBlockHeight int) SimappOption {
	return func(command []string) []string {
		return append(command, optionSimappBlockSize, strconv.Itoa(initialBlockHeight))
	}
}

// SimappWithNumBlocks provides numBlocks option for the simapp command.
func SimappWithNumBlocks(numBlocks int) SimappOption {
	return func(command []string) []string {
		return append(command, optionSimappNumBlocks, strconv.Itoa(numBlocks))
	}
}

// SimappWithBlockSize provides blockSize option for the simapp command.
func SimappWithBlockSize(blockSize int) SimappOption {
	return func(command []string) []string {
		return append(command, optionSimappBlockSize, strconv.Itoa(blockSize))
	}
}

// SimappWithLean provides lean option for the simapp command.
func SimappWithLean(lean bool) SimappOption {
	return func(command []string) []string {
		if lean {
			return append(command, optionSimappLean)
		}
		return command
	}
}

// SimappWithCommit provides commit option for the simapp command.
func SimappWithCommit(commit bool) SimappOption {
	return func(command []string) []string {
		if commit {
			return append(command, optionSimappCommit)
		}
		return command
	}
}

// SimappWithEnable provides enable option for the simapp command.
func SimappWithEnable(enable bool) SimappOption {
	return func(command []string) []string {
		if enable {
			return append(command, optionSimappEnabled)
		}
		return command
	}
}

// SimappWithGenesisTime provides genesisTime option for the simapp command.
func SimappWithGenesisTime(genesisTime int64) SimappOption {
	return func(command []string) []string {
		return append(command, optionSimappGenesisTime, strconv.Itoa(int(genesisTime)))
	}
}

// SimulationCommand returns the cli command for simapp tests.
// simName must be a test defined within the application (defaults to TestFullAppSimulation).
func SimulationCommand(appPath string, simName string, options ...SimappOption) step.Option {
	if simName == "" {
		simName = "TestFullAppSimulation"
	}

	command := []string{
		commandGoTest,
		optionGoBenchmem,
		fmt.Sprintf("-run=^%s$", simName),
		optionGoSimsTags,
		filepath.Join(appPath, "app"),
	}

	// Apply the options provided by the user
	for _, applyOption := range options {
		command = applyOption(command)
	}
	return step.Exec(gocmd.Name(), command...)
}
