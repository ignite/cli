package pluginsrpc

import (
	"testing"

	"github.com/stretchr/testify/require"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

func TestGetPluginId(t *testing.T) {
	testPlugin := chaincfg.Plugin{
		Name: "test",
	}
	require.Equal(t, "test", getPluginId(testPlugin))

	testPlugin = chaincfg.Plugin{
		Repo: "github.com/starport/test",
	}
	require.Equal(t, "test", getPluginId(testPlugin))
}

func TestListDirs(t *testing.T) {
	dirs, err := listDirs("./")
	require.NoError(t, err)
	var dirNames []string
	for _, dir := range dirs {
		dirNames = append(dirNames, dir.Name())
	}
	require.Equal(t, []string(nil), dirNames)
}

func TestListDirsMatch(t *testing.T) {
	dirs, err := listDirsMatch("./", "")
	require.NoError(t, err)
	var dirNames []string
	for _, dir := range dirs {
		dirNames = append(dirNames, dir.Name())
	}
	require.Equal(t, []string(nil), dirNames)
}

func TestListFilesMatch(t *testing.T) {
	files, err := listFilesMatch("./", "*ok.go")
	require.NoError(t, err)
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	require.Equal(t, []string{"hook.go"}, fileNames)
}
