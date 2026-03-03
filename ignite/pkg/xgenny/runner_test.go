package xgenny_test

import (
	"context"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
)

func TestMultipleGen(t *testing.T) {
	var (
		runner         = xgenny.NewRunner(context.Background(), t.TempDir())
		firstGen       = genny.New()
		secondGen      = genny.New()
		firstRunCount  int
		secondRunCount int
	)

	firstGen.RunFn(func(_ *genny.Runner) error {
		firstRunCount++
		return nil
	})

	secondGen.RunFn(func(_ *genny.Runner) error {
		secondRunCount++
		return nil
	})

	require.NoError(t, runner.Run(firstGen, secondGen))
	require.Equal(t, 1, firstRunCount, "first generator should run only once")
	require.Equal(t, 1, secondRunCount, "second generator should run only once")
}
