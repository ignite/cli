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
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	starportconf "github.com/tendermint/starport/starport/services/chain/conf"
)

// Build builds an app.
func (s *Chain) Build(ctx context.Context) error {
	if err := s.setup(ctx); err != nil {
		return err
	}
	conf, err := s.Config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

	steps, err := s.buildSteps(ctx, conf)
	if err != nil {
		return err
	}
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, steps...); err != nil {
		return err
	}

	fmt.Fprintf(s.stdLog(logStarport).out, "🗃  Installed. Use with: %s\n", infoColor(strings.Join(s.plugin.Binaries(), ", ")))
	return nil
}

func (s *Chain) buildSteps(ctx context.Context, conf starportconf.Config) (
	steps step.Steps, err error) {
	chainID, err := s.ID()
	if err != nil {
		return nil, err
	}

	ldflags := fmt.Sprintf(`-X github.com/cosmos/cosmos-sdk/version.Name=NewApp
-X github.com/cosmos/cosmos-sdk/version.ServerName=%sd
-X github.com/cosmos/cosmos-sdk/version.ClientName=%scli
-X github.com/cosmos/cosmos-sdk/version.Version=%s
-X github.com/cosmos/cosmos-sdk/version.Commit=%s
-X %s/cmd/%s/cmd.ChainID=%s`,
		s.app.Name,
		s.app.Name,
		s.version.tag,
		s.version.hash,
		s.app.ImportPath,
		s.app.D(),
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
				fmt.Fprintln(s.stdLog(logStarport).out, "\n📦 Installing dependencies...")
				return nil
			}),
			step.PostExec(captureBuildErr),
		).
		Add(s.stdSteps(logStarport)...).
		Add(step.Stderr(buildErr))...,
	))

	// If protocgen exists, compile the proto file
	protoScriptPath := filepath.Join(s.app.Path, "scripts/protocgen")
	if _, err := os.Stat(protoScriptPath); !os.IsNotExist(err) {
		steps.Add(step.New(step.NewOptions().
			Add(
				step.Exec(
					"/bin/bash",
					protoScriptPath,
				),
				step.PreExec(func() error {
					fmt.Fprintln(s.stdLog(logStarport).out, "🛠️  Building proto...")
					return nil
				}),
				step.PostExec(captureBuildErr),
			).
			Add(s.stdSteps(logStarport)...).
			Add(step.Stderr(buildErr))...,
		))
	}

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

	for _, binary := range s.plugin.Binaries() {
		steps.Add(step.New(step.NewOptions().
			Add(
				// ldflags somehow won't work if directly execute go binary.
				// bash stays as a workaround for now.
				step.Exec(
					"bash", "-c", fmt.Sprintf("go install -mod readonly -ldflags '%s'", ldflags),
				),
				step.Workdir(filepath.Join(s.app.Path, "cmd", binary)),
				step.PostExec(captureBuildErr),
			).
			Add(s.stdSteps(logStarport)...).
			Add(step.Stderr(buildErr))...,
		))
	}
	return steps, nil
}
