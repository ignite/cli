package exec

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
)

func TestExecSuccess(t *testing.T) {
	err := Exec(context.Background(), []string{"go", "version"})
	require.NoError(t, err)
}

func TestExecReturnsDetailedError(t *testing.T) {
	err := Exec(context.Background(), []string{"command-that-does-not-exist-ignite-test"})
	require.Error(t, err)

	var detailed *Error
	require.ErrorAs(t, err, &detailed)
	require.Contains(t, detailed.Error(), "error while running command")
	require.Contains(t, detailed.Error(), "command-that-does-not-exist-ignite-test")
}

func TestExecIncludesStdLogsWhenConfigured(t *testing.T) {
	if _, err := os.Stat("/bin/sh"); err != nil {
		t.Skip("/bin/sh not available")
	}

	err := Exec(
		context.Background(),
		[]string{"/bin/sh", "-c", "echo stdout-log; exit 1"},
		IncludeStdLogsToError(),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "stdout-log")
}

func TestStepOptionAddsStepOption(t *testing.T) {
	cfg := &execConfig{}
	StepOption(step.Workdir("/tmp"))(cfg)
	require.Len(t, cfg.stepOptions, 1)
}
