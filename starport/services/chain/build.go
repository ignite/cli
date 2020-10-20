package chain

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	starportconf "github.com/tendermint/starport/starport/services/chain/conf"
)

// Build builds an app.
func (s *Chain) Build(ctx context.Context) error {
	if err := s.setup(ctx); err != nil {
		return err
	}
	conf, err := s.config()
	if err != nil {
		return &CannotBuildAppError{err}
	}
	steps, binaries := s.buildSteps(ctx, conf)
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, steps...); err != nil {
		return err
	}
	fmt.Fprintf(s.stdLog(logStarport).out, "🗃  Installed. Use with: %s\n", infoColor(strings.Join(binaries, ", ")))
	return nil
}

func (s *Chain) buildSteps(ctx context.Context, conf starportconf.Config) (
	steps step.Steps, binaries []string) {
	ldflags := fmt.Sprintf(`'-X github.com/cosmos/cosmos-sdk/version.Name=NewApp 
	-X github.com/cosmos/cosmos-sdk/version.ServerName=%sd 
	-X github.com/cosmos/cosmos-sdk/version.ClientName=%scli 
	-X github.com/cosmos/cosmos-sdk/version.Version=%s 
	-X github.com/cosmos/cosmos-sdk/version.Commit=%s'`, s.app.Name, s.app.Name, s.version.tag, s.version.hash)
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
				fmt.Fprintln(s.stdLog(logStarport).out, "\n📦 Installing dependencies...")
				return nil
			}),
			step.PostExec(captureBuildErr),
		).
		Add(s.stdSteps(logStarport)...).
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
		Add(s.stdSteps(logBuild)...).
		Add(step.Stderr(buildErr))...,
	))

	// install the app.
	steps.Add(step.New(
		step.PreExec(func() error {
			fmt.Fprintln(s.stdLog(logStarport).out, "🛠️  Building the app...")
			return nil
		}),
	))
	installOptions, binaries := s.plugin.InstallCommands(ldflags)
	for _, execOption := range installOptions {
		execOption := execOption
		steps.Add(step.New(step.NewOptions().
			Add(
				execOption,
				step.PostExec(captureBuildErr),
			).
			Add(s.stdSteps(logStarport)...).
			Add(step.Stderr(buildErr))...,
		))
	}
	return steps, binaries
}
