package testdata

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
)

//go:embed plugins.yml
var ConfigYAML []byte

func GetConfig(t *testing.T) *pluginsconfig.Config {
	c := &pluginsconfig.Config{}

	err := yaml.NewDecoder(bytes.NewReader(ConfigYAML)).Decode(c)
	require.NoError(t, err)

	return c
}
