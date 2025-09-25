package gocmd_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
)

func TestList(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	ctx := context.Background()
	packages, err := gocmd.List(ctx, wd, []string{"-m", "-f={{.Path}}", "github.com/ignite/cli/v29"})
	assert.NoError(t, err)

	assert.Contains(t, packages, "github.com/ignite/cli/v29")
}
