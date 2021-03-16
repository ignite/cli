package chain

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	starporterrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/cosmosgen"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/giturl"
	"github.com/tendermint/starport/starport/pkg/gocmd"
	"github.com/tendermint/starport/starport/pkg/goenv"
)

// Build builds an app.
func (c *Chain) Build(ctx context.Context) error {
	if err := c.setup(ctx); err != nil {
		return err
	}

	if err := c.buildProto(ctx); err != nil {
		return err
	}

	steps, err := c.buildSteps()
	if err != nil {
		return err
	}

	return cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, steps...)
}

func (c *Chain) buildSteps() (steps step.Steps, err error) {
	chainID, err := c.ID()
	if err != nil {
		return nil, err
	}

	binary, err := c.Binary()
	if err != nil {
		return nil, err
	}

	ldflags := fmt.Sprintf(`-X github.com/cosmos/cosmos-sdk/version.Name=NewApp
-X github.com/cosmos/cosmos-sdk/version.ServerName=%sd
-X github.com/cosmos/cosmos-sdk/version.ClientName=%scli
-X github.com/cosmos/cosmos-sdk/version.Version=%s
-X github.com/cosmos/cosmos-sdk/version.Commit=%s
-X %s/cmd/%s/cmd.ChainID=%s`,
		c.app.Name,
		c.app.Name,
		c.sourceVersion.tag,
		c.sourceVersion.hash,
		c.app.ImportPath,
		binary,
		chainID,
	)
	var (
		buildErr = &bytes.Buffer{}
	)
	captureBuildErr := func(err error) error {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &CannotBuildAppError{errors.New(buildErr.String())}
		}
		return err
	}

	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				gocmd.Name(),
				"mod",
				"tidy",
			),
			step.PreExec(func() error {
				fmt.Fprintln(c.stdLog(logStarport).out, "📦 Installing dependencies...")
				return nil
			}),
			step.PostExec(captureBuildErr),
		).
		Add(c.stdSteps(logStarport)...).
		Add(step.Stderr(buildErr))...,
	))

	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				gocmd.Name(),
				"mod",
				"verify",
			),
			step.PostExec(captureBuildErr),
		).
		Add(c.stdSteps(logBuild)...).
		Add(step.Stderr(buildErr))...,
	))

	// install the app.
	steps.Add(step.New(
		step.PreExec(func() error {
			fmt.Fprintln(c.stdLog(logStarport).out, "🛠️  Building the app...")
			return nil
		}),
	))

	addInstallStep := func(binaryName, mainPath string) {
		installPath := filepath.Join(goenv.GetGOBIN(), binaryName)

		steps.Add(step.New(step.NewOptions().
			Add(
				// ldflags somehow won't work if directly execute go binary.
				// bash stays as a workaround for now.
				step.Exec(
					"bash", "-c", fmt.Sprintf("%s build -mod readonly -o %s -ldflags '%s'", gocmd.Name(), installPath, ldflags),
				),
				step.Workdir(mainPath),
				step.PostExec(captureBuildErr),
			).
			Add(c.stdSteps(logStarport)...).
			Add(step.Stderr(buildErr))...,
		))
	}

	cmdPath := filepath.Join(c.app.Path, "cmd")

	addInstallStep(binary, filepath.Join(cmdPath, c.app.D()))

	if c.Version.Major().Is(cosmosver.Launchpad) {
		addInstallStep(c.BinaryCLI(), filepath.Join(cmdPath, c.app.CLI()))
	}

	return steps, nil
}

func (c *Chain) buildProto(ctx context.Context) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	// If proto dir exists, compile the proto files.
	protoPath := filepath.Join(c.app.Path, conf.Build.Proto.Path)
	if _, err := os.Stat(protoPath); os.IsNotExist(err) {
		return nil
	}

	if err := cosmosgen.InstallDependencies(context.Background(), c.app.Path); err != nil {
		if err == cosmosgen.ErrProtocNotInstalled {
			return starporterrors.ErrStarportRequiresProtoc
		}
		return err
	}

	fmt.Fprintln(c.stdLog(logStarport).out, "🛠️  Building proto...")

	options := []cosmosgen.Option{
		cosmosgen.WithGoGeneration(c.app.ImportPath),
		cosmosgen.IncludeDirs(conf.Build.Proto.ThirdPartyPaths),
	}

	enableThirdPartyModuleCodegen := !c.protoBuiltAtLeastOnce && c.options.isThirdPartyModuleCodegenEnabled

	// generate Vuex code as well if it is enabled.
	if conf.Client.Vuex.Path != "" {
		storeRootPath := filepath.Join(c.app.Path, conf.Client.Vuex.Path, "generated")
		options = append(options,
			cosmosgen.WithVuexGeneration(
				enableThirdPartyModuleCodegen,
				func(m module.Module) string {
					return filepath.Join(storeRootPath, giturl.UserAndRepo(m.Pkg.GoImportName), m.Pkg.Name, "module")
				},
				storeRootPath,
			),
		)
	}

	if err := cosmosgen.Generate(ctx, c.app.Path, conf.Build.Proto.Path, options...); err != nil {
		return &CannotBuildAppError{err}
	}

	c.protoBuiltAtLeastOnce = true

	return nil
}
