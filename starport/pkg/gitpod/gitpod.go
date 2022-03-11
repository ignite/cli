package gitpod

import (
	"bytes"
	"context"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

// IsOnGitpod reports whether if running on Gitpod or not.
func IsOnGitpod() bool {
	return os.Getenv("GITPOD_WORKSPACE_ID") != ""
}

func URLForPort(ctx context.Context, port string) (string, error) {
	buf := bytes.Buffer{}
	if err := exec.Exec(ctx, []string{"gp", "url", port}, exec.StepOption(step.Stdout(&buf))); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}
