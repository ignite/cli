package plugins_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
	"github.com/ignite/cli/ignite/config/plugins/testdata"
)

func TestParse(t *testing.T) {
	// Arrange: Initialize a reader with the previous version
	r := bytes.NewReader(testdata.ConfigYAML)

	// Act
	cfg, err := pluginsconfig.Parse(r)

	// Assert
	require.NoError(t, err)

	// Assert: Parse must return the latest version
	require.Equal(t, testdata.GetConfig(t), cfg)
}
