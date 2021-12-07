package starportcmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/services/plugin"
	"github.com/tendermint/starport/starport/services/plugin/mocks"
)

func Test_ScaffoldPlugins_NotInstalled(t *testing.T) {
	// Test fixtures
	testPlugin := chainconfig.Plugin{Name: "test-1"}

	// Mocks
	mockLoader := mocks.Loader{}
	pluginLoader = &mockLoader

	// Intentionally returns false
	mockLoader.On("IsInstalled", testPlugin).Return(false)

	// Test
	cmds := NewScaffoldPlugins("mars", []chainconfig.Plugin{testPlugin})

	// Asserts
	mockLoader.AssertExpectations(t)
	assert.Zero(t, len(cmds))
}

func Test_ScaffoldPlugins_Installed(t *testing.T) {
	// Test fixtures
	testPlugin := chainconfig.Plugin{Name: "test-0"}

	// Mocks
	mockLoader := mocks.Loader{}
	pluginLoader = &mockLoader

	mockLoader.On("IsInstalled", testPlugin).Return(true)

	mockPlugin := mocks.StarportPlugin{}
	mockFuncs := []plugin.FuncSpec{
		{Name: "func1"},
		{Name: "func2"},
	}

	mockLoader.On("LoadPlugin", testPlugin, pluginHome).Return(&mockPlugin, nil)

	mockPlugin.On("List").Return(mockFuncs)
	mockPlugin.On("Help", "func1").Return("")
	mockPlugin.On("Help", "func2").Return("")

	// Test
	cmds := NewScaffoldPlugins("mars", []chainconfig.Plugin{testPlugin})

	// Asserts
	mockLoader.AssertExpectations(t)
	assert.Equal(t, 1, len(cmds))
	assert.Equal(t, testPlugin.Name, cmds[0].Use)

	for i, cmd := range cmds[0].Commands() {
		assert.Equal(t, mockFuncs[i].Name, cmd.Use)
	}
}
