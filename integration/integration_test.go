// Package integration_test integration test Starport and scaffolded apps.
package integration_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/xexec"
)

const (
	relayerVersion = "3080dab10e37e6db9e691c343b1c0cdb6d845753"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if err := checkSystemRequirements(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func checkSystemRequirements() error {
	if !xexec.IsCommandAvailable("starport") {
		return errors.New("starport needs to be installed")
	}
	if !xexec.IsCommandAvailable("rly") {
		return errors.New("relayer needs to be installed")
	}
	version := &bytes.Buffer{}
	return cmdrunner.
		New().
		Run(context.Background(), step.New(
			step.Exec("rly", "version"),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if !strings.Contains(version.String(), relayerVersion) {
					return fmt.Errorf("relayer is not at the required version %q", relayerVersion)
				}
				return nil
			}),
			step.Stdout(version),
		))
}
