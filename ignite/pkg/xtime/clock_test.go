package xtime_test

import (
	"testing"
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/xtime"
	"github.com/stretchr/testify/require"
)

func TestClockSystem(t *testing.T) {
	c := xtime.NewClockSystem()
	require.False(t, c.Now().IsZero())
	require.Panics(t, func() { c.Add(time.Second) })
}

func TestClockMock(t *testing.T) {
	timeSample := time.Now()
	c := xtime.NewClockMock(timeSample)
	require.True(t, c.Now().Equal(timeSample))
	c.Add(time.Second)
	require.True(t, c.Now().Equal(timeSample.Add(time.Second)))
}
