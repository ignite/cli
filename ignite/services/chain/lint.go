package chain

import (
	"context"
	"fmt"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
)

var golangCiLintVersion = "latest"

// Lint runs the linting process for the chain.
// It uses golangci-lint to lint the chain's codebase.
func (c *Chain) Lint(ctx context.Context) error {
	if err := exec.Exec(ctx, []string{"go", "install", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", golangCiLintVersion)}); err != nil {
		return fmt.Errorf("failed to install golangci-lint: %w", err)
	}
	return exec.Exec(ctx, []string{"golangci-lint", "run", "./...", "--out-format=tab"}, exec.IncludeStdLogsToError())
}
