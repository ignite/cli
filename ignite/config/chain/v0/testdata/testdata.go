package testdata

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	v0 "github.com/ignite/cli/ignite/config/chain/v0"
)

//go:embed config.yaml
var ConfigYAML []byte

func GetConfig(t *testing.T) *v0.Config {
	c := &v0.Config{}

	err := yaml.NewDecoder(bytes.NewReader(ConfigYAML)).Decode(c)
	require.NoError(t, err)

	err = c.SetDefaults()
	require.NoError(t, err)

	return c
}
