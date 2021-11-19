package starportcmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/services/plugin/mocks"
)

func Test_ScaffoldPlugins(t *testing.T) {
	// Test fixtures
	testConfigs := []struct {
		IsInstalled bool
		Plugin      chainconfig.Plugin
	}{
		{
			IsInstalled: true,
			Plugin:      chainconfig.Plugin{Name: "test-0"},
		},

		{
			IsInstalled: false,
			Plugin:      chainconfig.Plugin{Name: "test-1"},
		},

		{
			IsInstalled: true,
			Plugin:      chainconfig.Plugin{Name: "test-2"},
		},
	}

	allPlugins := make([]chainconfig.Plugin, len(testConfigs))
	for i, test := range testConfigs {
		allPlugins[i] = test.Plugin
	}

	// Mocks
	mockLoader := mocks.Loader{}
	pluginLoader = &mockLoader

	for _, test := range testConfigs {
		mockLoader.On("IsInstalled", test.Plugin).Return(test.IsInstalled)
	}

	// Test
	cmds := NewScaffoldPlugins(allPlugins)

	// Asserts
	mockLoader.AssertExpectations(t)

	installedPlugins := make([]chainconfig.Plugin, 0)
	for _, test := range testConfigs {
		if test.IsInstalled {
			installedPlugins = append(installedPlugins, test.Plugin)
		}
	}

	for i, cmd := range cmds {
		assert.Equal(t, cmd.Use, installedPlugins[i].Name)
	}
}
