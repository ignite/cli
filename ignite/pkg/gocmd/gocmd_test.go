package gocmd_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

<<<<<<< HEAD
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
=======
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
>>>>>>> 3919d6bb (feat(cosmosgen): fetch fallback buf token (#4805))
)

func TestList(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	ctx := context.Background()
	packages, err := gocmd.List(ctx, wd, []string{"-m", "-f={{.Path}}", "github.com/ignite/cli/v28"})
	assert.NoError(t, err)

	assert.Contains(t, packages, "github.com/ignite/cli/v28")
}
