package chain

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	starporterrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosprotoc"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/goenv"
	"github.com/tendermint/starport/starport/pkg/xos"
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
	if err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, steps...); err != nil {
		return err
	}

	binaries, err := c.Binaries()
	if err != nil {
		return err
	}

	fmt.Fprintf(c.stdLog(logStarport).out, "🗃  Installed. Use with: %s\n", infoColor(strings.Join(binaries, ", ")))
	return nil
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
				"go",
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
				"go",
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
					"bash", "-c", fmt.Sprintf("go build -mod readonly -o %s -ldflags '%s'", installPath, ldflags),
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

	if err := cosmosprotoc.InstallDependencies(context.Background(), c.app.Path); err != nil {
		if err == cosmosprotoc.ErrProtocNotInstalled {
			return starporterrors.ErrStarportRequiresProtoc
		}
		return err
	}

	fmt.Fprintln(c.stdLog(logStarport).out, "🛠️  Building proto...")

	var (
		includePaths = xos.PrefixPathToList(conf.Build.Proto.ThirdPartyPaths, c.app.Path)
		targets      = []cosmosprotoc.Target{
			cosmosprotoc.WithGoGeneration(c.app.ImportPath),
		}
	)

	frontendPath := filepath.Join(c.app.Path, conf.Frontend.Path)
	if _, err := os.Stat(frontendPath); err == nil && conf.Build.Proto.JS.Out != "" {
		path := filepath.Join(c.app.Path, conf.Build.Proto.JS.Out)
		targets = append(targets, cosmosprotoc.WithJSGeneration(path))
	}

	if err := cosmosprotoc.Generate(ctx, c.app.Path, protoPath, includePaths, targets[0], targets[1:]...); err != nil {
		return &CannotBuildAppError{err}
	}

	return nil
}
