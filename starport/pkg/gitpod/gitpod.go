package gitpod

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
)

// IsOnGitpod reports whether if running on Gitpod or not.
func IsOnGitpod() bool {
	return os.Getenv("GITPOD_WORKSPACE_ID") != ""
}

func GitPodPortUrl(port string) string {
	buf := bytes.Buffer{}
	ctx := context.Background()
	if err := cmdrunner.New(cmdrunner.DefaultStdout(&buf)).Run(ctx, step.New(step.Exec("gp", "url", port))); err != nil {
		return ""
	}
	output, err := io.ReadAll(&buf)
	if err != nil {
		return ""
	}
	return strings.Trim(string(output), "\n")
}
