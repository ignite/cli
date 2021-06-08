package chain

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/cosmosgen"
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

	ldflags := fmt.Sprintf(`-X github.com/cosmos/cosmos-sdk/version.Name=%[1]s
-X github.com/cosmos/cosmos-sdk/version.ServerName=%[2]sd
-X github.com/cosmos/cosmos-sdk/version.ClientName=%[2]scli
-X github.com/cosmos/cosmos-sdk/version.Version=%[3]s
-X github.com/cosmos/cosmos-sdk/version.Commit=%[4]s
-X %[5]s/cmd/%[6]s/cmd.ChainID=%[7]s`,
		strings.Title(c.app.Name),
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
			fmt.Fprintln(c.stdLog(logStarport).out, "🛠️  Building the blockchain...")
			return nil
		}),
	))

	addInstallStep := func(binaryName, mainPath string) {
		installPath := filepath.Join(goenv.Bin(), binaryName)

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

	return steps, nil
}

func (c *Chain) buildProto(ctx context.Context) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	if err := cosmosgen.InstallDependencies(context.Background(), c.app.Path); err != nil {
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
					parsedGitURL, _ := giturl.Parse(m.Pkg.GoImportName)
					return filepath.Join(storeRootPath, parsedGitURL.UserAndRepo(), m.Pkg.Name, "module")
				},
				storeRootPath,
			),
		)
	}
	if conf.Client.OpenAPI.Path != "" {
		options = append(options, cosmosgen.WithOpenAPIGeneration(conf.Client.OpenAPI.Path))
	}

	if err := cosmosgen.Generate(ctx, c.app.Path, conf.Build.Proto.Path, options...); err != nil {
		return &CannotBuildAppError{err}
	}

	c.protoBuiltAtLeastOnce = true

	return nil
}
