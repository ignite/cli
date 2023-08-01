package chaincmd

import (
	"path/filepath"
	"strconv"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gocmd"
)

const (
	optionSimappGenesis                = "-Genesis"
	optionSimappParams                 = "-Params"
	optionSimappExportParamsPath       = "-ExportParamsPath"
	optionSimappExportParamsHeight     = "-ExportParamsHeight"
	optionSimappExportStatePath        = "-ExportStatePath"
	optionSimappExportStatsPath        = "-ExportStatsPath"
	optionSimappSeed                   = "-Seed"
	optionSimappInitialBlockHeight     = "-InitialBlockHeight"
	optionSimappNumBlocks              = "-NumBlocks"
	optionSimappBlockSize              = "-BlockSize"
	optionSimappLean                   = "-Lean"
	optionSimappCommit                 = "-Commit"
	optionSimappSimulateEveryOperation = "-SimulateEveryOperation"
	optionSimappPrintAllInvariants     = "-PrintAllInvariants"
	optionSimappEnabled                = "-Enabled"
	optionSimappVerbose                = "-Verbose"
	optionSimappPeriod                 = "-Period"
	optionSimappGenesisTime            = "-GenesisTime"

	commandGoTest       = "test"
	optionGoBenchmem    = "-benchmem"
	optionGoSimappRun   = "-run=^$"
	optionGoSimappBench = "-bench=^BenchmarkSimulation"
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
		return append(command, optionSimappInitialBlockHeight, strconv.Itoa(initialBlockHeight))
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

// SimappWithSimulateEveryOperation provides simulateEveryOperation option for the simapp command.
func SimappWithSimulateEveryOperation(simulateEveryOperation bool) SimappOption {
	return func(command []string) []string {
		if simulateEveryOperation {
			return append(command, optionSimappSimulateEveryOperation)
		}
		return command
	}
}

// SimappWithPrintAllInvariants provides printAllInvariants option for the simapp command.
func SimappWithPrintAllInvariants(printAllInvariants bool) SimappOption {
	return func(command []string) []string {
		if printAllInvariants {
			return append(command, optionSimappPrintAllInvariants)
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

// SimappWithVerbose provides verbose option for the simapp command.
func SimappWithVerbose(verbose bool) SimappOption {
	return func(command []string) []string {
		if verbose {
			return append(command, optionSimappVerbose)
		}
		return command
	}
}

// SimappWithPeriod provides period option for the simapp command.
func SimappWithPeriod(period uint) SimappOption {
	return func(command []string) []string {
		return append(command, optionSimappPeriod, strconv.Itoa(int(period)))
	}
}

// SimappWithGenesisTime provides genesisTime option for the simapp command.
func SimappWithGenesisTime(genesisTime int64) SimappOption {
	return func(command []string) []string {
		return append(command, optionSimappGenesisTime, strconv.Itoa(int(genesisTime)))
	}
}

// SimulationCommand returns the cli command for simapp tests.
func SimulationCommand(appPath string, options ...SimappOption) step.Option {
	command := []string{
		commandGoTest,
		optionGoBenchmem,
		optionGoSimappRun,
		optionGoSimappBench,
		filepath.Join(appPath, "app"),
	}

	// Apply the options provided by the user
	for _, applyOption := range options {
		command = applyOption(command)
	}
	return step.Exec(gocmd.Name(), command...)
}
