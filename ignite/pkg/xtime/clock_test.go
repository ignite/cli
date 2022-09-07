package xtime_test

import (
	"github.com/ignite/cli/ignite/pkg/xtime"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
