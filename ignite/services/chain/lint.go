package chain

import (
	"context"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/exec"
)

// Lint runs the linting process for the chain.
// It uses golangci-lint to lint the chain's codebase.
func (c *Chain) Lint(ctx context.Context) error {
	cmd := []string{
		"go",
		"tool",
		"github.com/golangci/golangci-lint/cmd/golangci-lint",
		"run",
		"./...",
		"--out-format=tab",
	}

	return exec.Exec(ctx, cmd, exec.IncludeStdLogsToError())
}
