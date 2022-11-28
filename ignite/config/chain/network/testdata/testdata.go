package testdata

import (
	"bytes"
	_ "embed"
	"github.com/ignite/cli/ignite/config/chain/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var ConfigYAML []byte

func GetConfig(t *testing.T) *v1.Config {
	c := &v1.Config{}

	err := yaml.NewDecoder(bytes.NewReader(ConfigYAML)).Decode(c)
	require.NoError(t, err)

	err = c.SetDefaults()
	require.NoError(t, err)

	return c
}
