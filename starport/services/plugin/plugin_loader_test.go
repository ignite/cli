package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/chainconfig"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Test_Find(t *testing.T) {
	testConfigLoader := configLoader{}
	createEmptyFile := func(name string) {
		d := []byte("")
		check(os.WriteFile(name, d, 0644))
	}

	starPortPath, _ := chainconfig.ConfigDirPath()
	starPortPathDirExist, _ := testConfigLoader.IsExists(starPortPath)
	if !starPortPathDirExist {
		mkErr := os.Mkdir(starPortPath, 0755)
		check(mkErr)
	}

	pluginsPath := filepath.Join(starPortPath, "plugins")
	pluginDirExist, _ := testConfigLoader.IsExists(pluginsPath)

	if !pluginDirExist {
		mkErr := os.Mkdir(pluginsPath, 0755)
		check(mkErr)
	}

	defer os.RemoveAll(pluginsPath)

	testPluginName := "testPlugin"
	pluginsPathWithSample := filepath.Join(pluginsPath, testPluginName)
	pluginsPathWithSampleExist, _ := testConfigLoader.IsExists(pluginsPathWithSample)
	if !pluginsPathWithSampleExist {
		mkErr := os.Mkdir(pluginsPathWithSample, 0755)
		check(mkErr)
	}

	defer os.RemoveAll(pluginsPathWithSample)

	testFileName := testPluginName + ".so"
	filePath := filepath.Join(pluginsPathWithSample, testFileName)
	createEmptyFile(filePath)
	fileLists := testConfigLoader.Find(pluginsPathWithSample, ".so")
	fmt.Println(fileLists)

	assertEqualValue := pluginsPathWithSample + "/testPlugin.so"
	require.Equal(t, []string{assertEqualValue}, fileLists)

}

func Test_IsInstalled(t *testing.T) {
	// Prepare test
	pluginConfig := chainconfig.Plugin{
		Name:          "unittest-dummy-plugin",
		Description:   "test plugin to Run check isInstalled",
		RepositoryURL: "github.com/dummy-user/plugin-repo",
	}

	chainID := "mars"

	testConfigLoader := configLoader{}
	starportHome, _ := chainconfig.ConfigDirPath()

	if exist, _ := testConfigLoader.IsExists(starportHome); !exist {
		err := os.Mkdir(starportHome, 0755)
		if err != nil {
			panic(err)
		}
	}

	pluginsHome := filepath.Join(starportHome, "plugins")
	pluginDirExist, _ := testConfigLoader.IsExists(pluginsHome)
	if !pluginDirExist {
		mkErr := os.Mkdir(pluginsHome, 0755)
		check(mkErr)
	}

	// Create dummy plugin files.
	pluginRepoHome := filepath.Join(pluginsHome, chainID, "plugin-repo")
	pluginPath := filepath.Join(pluginRepoHome, pluginConfig.Name)

	if exist, _ := testConfigLoader.IsExists(pluginPath); !exist {
		err := os.MkdirAll(pluginPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	defer os.RemoveAll(pluginRepoHome)

	// Test
	pluginLoader, _ := NewLoader(chainID)
	isPluginInstalled := pluginLoader.IsInstalled(pluginConfig)
	require.Equal(t, false, isPluginInstalled)

	// Create dummy
	os.WriteFile(fmt.Sprintf("%s/%s.so", pluginPath, pluginConfig.Name), []byte(""), 0644)

	isPluginInstalled = pluginLoader.IsInstalled(pluginConfig)
	require.Equal(t, true, isPluginInstalled)
}

func Test_CheckMandatory(t *testing.T) {
	tests := []struct {
		Desc      string
		Spec      *starportplugin
		ExpectErr error
	}{
		{
			Desc: "Required function is not exist",
			Spec: &starportplugin{
				name: "dummy",
				funcSpecs: map[string]FuncSpec{
					"Init": {
						Name:       "Init",
						ParamTypes: []reflect.Type{},
					},
				},
			},
			ExpectErr: ErrPluginWrongSpec,
		},

		{
			Desc: "Required function has diff parameters",
			Spec: &starportplugin{
				name: "dummy",
				funcSpecs: map[string]FuncSpec{
					"Init": {
						Name:       "Init",
						ParamTypes: []reflect.Type{},
					},

					"Help": {
						Name:       "Help",
						ParamTypes: []reflect.Type{},
					},
				},
			},
			ExpectErr: ErrPluginWrongSpec,
		},

		{
			Desc: "Success",
			Spec: &starportplugin{
				name: "dummy",
				funcSpecs: map[string]FuncSpec{
					"Init": {
						Name:       "Init",
						ParamTypes: []reflect.Type{},
					},

					"Help": {
						Name:       "Help",
						ParamTypes: []reflect.Type{reflect.TypeOf("")},
					},
				},
			},
			ExpectErr: nil,
		},
	}

	for _, test := range tests {
		// Prepare test
		loader := configLoader{}
		loader.pluginSpec = test.Spec

		// Test
		err := loader.checkMandatoryFunctions(test.Spec.funcSpecs)

		// Asserts
		assert.Equal(t, test.ExpectErr, err)
	}
}
