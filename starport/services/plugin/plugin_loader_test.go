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

func Test_IsExists(t *testing.T) {
	testConfigLoader := configLoader{}
	var configPath, err = chainconfig.ConfigDirPath()
	fmt.Println(configPath)
	exist := false
	_, errStat := os.Stat(configPath)
	if errStat == nil {
		exist = true
	}

	doesExists, err := (&testConfigLoader).IsExists(configPath)
	require.NoError(t, err)
	require.Equal(t, exist, doesExists)
}

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

	var starPortPath, _ = chainconfig.ConfigDirPath()
	starPortPathDirExist, _ := (&testConfigLoader).IsExists(starPortPath)
	if !starPortPathDirExist {
		mkErr := os.Mkdir(starPortPath, 0755)
		check(mkErr)
	}

	var pluginsPath = filepath.Join(starPortPath, "plugins")
	pluginDirExist, _ := (&testConfigLoader).IsExists(pluginsPath)

	if !pluginDirExist {
		mkErr := os.Mkdir(pluginsPath, 0755)
		check(mkErr)
	}

	defer os.RemoveAll(pluginsPath)

	var testPluginName = "testPlugin"
	var pluginsPathWithSample = filepath.Join(pluginsPath, testPluginName)
	pluginsPathWithSampleExist, _ := (&testConfigLoader).IsExists(pluginsPathWithSample)
	if !pluginsPathWithSampleExist {
		mkErr := os.Mkdir(pluginsPathWithSample, 0755)
		check(mkErr)
	}

	defer os.RemoveAll(pluginsPathWithSample)
	var testFileName = testPluginName + ".so"
	var filePath = filepath.Join(pluginsPathWithSample, testFileName)
	createEmptyFile(filePath)
	fileLists := (&testConfigLoader).Find(pluginsPathWithSample, ".so")
	fmt.Println(fileLists)
	var assertEqualValue = pluginsPathWithSample + "/testPlugin.so"
	require.Equal(t, []string{assertEqualValue}, fileLists)

}

func Test_IsInstalled(t *testing.T) {
	pluginSample := chainconfig.Plugin{
		Name:          "testPlugin",
		Description:   "testPlugin to Run check isInstalled",
		RepositoryURL: "github.com/this/is/fake/for/test",
	}

	testConfigLoader := configLoader{}
	createEmptyFile := func(name string) {
		d := []byte("")
		check(os.WriteFile(name, d, 0644))
	}

	var starPortPath, _ = chainconfig.ConfigDirPath()
	starPortPathDirExist, _ := (&testConfigLoader).IsExists(starPortPath)
	if !starPortPathDirExist {
		mkErr := os.Mkdir(starPortPath, 0755)
		check(mkErr)
	}

	var pluginsPath = filepath.Join(starPortPath, "plugins")
	pluginDirExist, _ := (&testConfigLoader).IsExists(pluginsPath)

	if !pluginDirExist {
		mkErr := os.Mkdir(pluginsPath, 0755)
		check(mkErr)
	}

	defer os.RemoveAll(pluginsPath)

	var testPluginName = pluginSample.Name
	var pluginsPathWithSample = filepath.Join(pluginsPath, testPluginName)
	pluginsPathWithSampleExist, _ := (&testConfigLoader).IsExists(pluginsPathWithSample)
	if !pluginsPathWithSampleExist {
		mkErr := os.Mkdir(pluginsPathWithSample, 0755)
		check(mkErr)
	}

	defer os.RemoveAll(pluginsPathWithSample)
	var testFileName = testPluginName + ".so"
	var filePath = filepath.Join(pluginsPathWithSample, testFileName)
	createEmptyFile(filePath)

	pluginLoader, _ := NewLoader()

	doesInstalledTrue := pluginLoader.IsInstalled(pluginSample)
	fmt.Println(doesInstalledTrue)
	require.Equal(t, true, doesInstalledTrue)

	os.Remove(filePath)
	doesInstalledFalse := pluginLoader.IsInstalled(pluginSample)
	fmt.Println(doesInstalledFalse)
	require.Equal(t, false, doesInstalledFalse)

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
		err := loader.checkMandatoryFunctions()

		// Asserts
		assert.Equal(t, test.ExpectErr, err)
	}
}
