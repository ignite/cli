package gocmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/gocmd"
)

func TestAvailable(t *testing.T) {
	b := gocmd.Available()

	assert.True(t, b)
}

func TestIsMinVersion(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	b, err := gocmd.IsMinVersion("1.7")

	require.NoError(err)
	assert.True(b)

	b, err = gocmd.IsMinVersion("9999.999")

	require.NoError(err)
	assert.False(b)
}
