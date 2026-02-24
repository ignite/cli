package modulecreate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddModuleToAppConfig(t *testing.T) {
	content := readFixture(t, "../../app/files/app/app_config.go.plush")

	modified, err := addModuleToAppConfig(content, "blog")
	require.NoError(t, err)
	normalized := normalizedExpr(modified)
	require.Equal(t, 4, strings.Count(normalized, "blogmoduletypes.ModuleName"))
	require.Contains(t, normalized, "Config:appconfig.WrapAny(&blogmoduletypes.Module{}),")

	modified, err = addModuleToAppConfig(modified, "blog")
	require.NoError(t, err)
	require.Equal(t, 4, strings.Count(normalizedExpr(modified), "blogmoduletypes.ModuleName"))
}

func TestAddModuleToLegacyAppConfig(t *testing.T) {
	content := readFixture(t, "../../../pkg/cosmosanalysis/module/testdata/earth/app/app_config.go")

	modified, err := addModuleToAppConfig(content, "venus")
	require.NoError(t, err)
	require.Equal(t, 4, strings.Count(normalizedExpr(modified), "venusmoduletypes.ModuleName"))
}
