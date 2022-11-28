package plugins_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/config/plugins"
)

func TestConfigDecode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	f, err := os.Open("testdata/plugins.yml")
	require.NoError(err)
	defer f.Close()
	var cfg plugins.Config

	err = cfg.Decode(f)

	require.NoError(err)
	expected := plugins.Config{
		Plugins: []plugins.Plugin{
			{
				Path: "/path/to/plugin1",
			},
			{
				Path: "/path/to/plugin2",
				With: map[string]string{"foo": "bar", "bar": "baz"},
			},
		},
	}
	assert.Equal(expected, cfg)
}
