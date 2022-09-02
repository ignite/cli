package yaml_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	xyaml "github.com/ignite/cli/ignite/pkg/yaml"
)

func TestUnmarshalWithCustomMapType(t *testing.T) {
	// Arrange
	input := `
    foo:
      bar: baz
    `
	output := xyaml.Map{}

	// Act
	err := yaml.Unmarshal([]byte(input), &output)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, output["foo"])
	require.IsType(t, (map[string]interface{})(nil), output["foo"])
}

func TestUnmarshalWithNativeMapType(t *testing.T) {
	// Arrange
	input := `
    foo:
      bar: baz
    `
	output := make(map[string]interface{})

	// Act
	err := yaml.Unmarshal([]byte(input), &output)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, output["foo"])
	require.IsType(t, (map[interface{}]interface{})(nil), output["foo"])
}
