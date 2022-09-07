package gocmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/goenv"
	"github.com/ignite/cli/ignite/pkg/xexec"
)

const (
	// CommandInstall represents go "install" command.
	CommandInstall = "install"

	// CommandBuild represents go "build" command.
	CommandBuild = "build"

	// CommandMod represents go "mod" command.
	CommandMod = "mod"

	// CommandModTidy represents go mod "tidy" command.
	CommandModTidy = "tidy"

	// CommandModVerify represents go mod "verify" command.
	CommandModVerify = "verify"

	// CommandEnv represends go "env" command.
	CommandEnv = "env"
)

const (
	FlagMod              = "-mod"
	FlagModValueReadOnly = "readonly"
	FlagLdflags          = "-ldflags"
	FlagOut              = "-o"
)

const (
	EnvGOOS      = "GOOS"
	EnvGOARCH    = "GOARCH"
	EnvGOVERSION = "GOVERSION"
)

// Available returns true if the go command is available.
func Available() bool {
	return xexec.IsCommandAvailable("go")
}

// IsMinVersion returns true if v is less or equal to current go version.
func IsMinVersion(v string) (bool, error) {
	minVersion, err := semver.ParseTolerant(v)
	if err != nil {
		return false, errors.Wrapf(err, "semver parse %s", v)
	}
	e, err := Env(EnvGOVERSION)
	if err != nil {
		return false, err
	}
	e = e[2:] // remove go prefix
	version, err := semver.ParseTolerant(e)
	if err != nil {
		return false, errors.Wrapf(err, "semver parse %s", e)
	}
	return minVersion.LTE(version), nil
}

// Env returns the output of "go env key" command.
func Env(key string) (string, error) {
	var b bytes.Buffer
	err := exec.Exec(context.Background(), []string{
		Name(),
		CommandEnv,
		key,
	}, exec.StepOption(step.Stdout(&b)))
	return b.String(), err
}

// Name returns the name of Go binary to use.
func Name() string {
	custom := os.Getenv("GONAME")
	if custom != "" {
		return custom
	}
	return "go"
}

// ModTidy runs go mod tidy on path with options.
func ModTidy(ctx context.Context, path string, options ...exec.Option) error {
	return exec.Exec(ctx, []string{Name(), CommandMod, CommandModTidy}, append(options, exec.StepOption(step.Workdir(path)))...)
}

// ModVerify runs go mod verify on path with options.
func ModVerify(ctx context.Context, path string, options ...exec.Option) error {
	return exec.Exec(ctx, []string{Name(), CommandMod, CommandModVerify}, append(options, exec.StepOption(step.Workdir(path)))...)
}

// BuildPath runs go install on cmd folder with options.
func BuildPath(ctx context.Context, output, binary, path string, flags []string, options ...exec.Option) error {
	binaryOutput, err := binaryPath(output, binary)
	if err != nil {
		return err
	}
	command := []string{
		Name(),
		CommandBuild,
		FlagOut, binaryOutput,
	}
	command = append(command, flags...)
	command = append(command, ".")
	return exec.Exec(ctx, command, append(options, exec.StepOption(step.Workdir(path)))...)
}

// BuildAll runs go build ./... on path with options.
func BuildAll(ctx context.Context, out, path string, flags []string, options ...exec.Option) error {
	command := []string{
		Name(),
		CommandBuild,
		FlagOut, out,
	}
	command = append(command, flags...)
	command = append(command, "./...")
	return exec.Exec(ctx, command, append(options, exec.StepOption(step.Workdir(path)))...)
}

// InstallAll runs go install ./... on path with options.
func InstallAll(ctx context.Context, path string, flags []string, options ...exec.Option) error {
	command := []string{
		Name(),
		CommandInstall,
	}
	command = append(command, flags...)
	command = append(command, "./...")
	return exec.Exec(ctx, command, append(options, exec.StepOption(step.Workdir(path)))...)
}

// Ldflags returns a combined ldflags set from flags.
func Ldflags(flags ...string) string {
	return strings.Join(flags, " ")
}

// BuildTarget builds a GOOS:GOARCH pair.
func BuildTarget(goos, goarch string) string {
	return fmt.Sprintf("%s:%s", goos, goarch)
}

// ParseTarget parses GOOS:GOARCH pair.
func ParseTarget(t string) (goos, goarch string, err error) {
	parsed := strings.Split(t, ":")
	if len(parsed) != 2 {
		return "", "", errors.New("invalid Go target, expected in GOOS:GOARCH format")
	}

	return parsed[0], parsed[1], nil
}

// PackageLiteral returns the string representation of package part of go get [package].
func PackageLiteral(path, version string) string {
	return fmt.Sprintf("%s@%s", path, version)
}

// binaryPath determines the path where binary will be located at.
func binaryPath(output, binary string) (string, error) {
	if output != "" {
		outputAbs, err := filepath.Abs(output)
		if err != nil {
			return "", err
		}
		return filepath.Join(outputAbs, binary), nil
	}
	return filepath.Join(goenv.Bin(), binary), nil
}
