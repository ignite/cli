package plugin

import (
	"fmt"
	"testing"

	"github.com/tendermint/starport/starport/chainconfig"
)

func Test_Loader(t *testing.T) {
	path, err := chainconfig.ConfigDirPath()

	fmt.Println(path, err)
}
