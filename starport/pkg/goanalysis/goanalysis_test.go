package goanalysis

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindTypes(t *testing.T) {
	srcDirPath, _ := filepath.Abs(filepath.Join("..", "..", "services", "chain"))
	typesFound, err := FindTypes(srcDirPath, []string{"App", "Hello", "Chain"})

	require.Nil(t, err, "Error should be nil")
	require.Equal(t, []Type{
		{Name: "App"},
		{Name: "Chain"},
	}, typesFound)
}
