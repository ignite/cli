package services

import (
	"os"

	"github.com/tendermint/starport/starport/pkg/xfilepath"
)

var (
	// StarportConfPath returns the Starport Configuration directory
	StarportConfPath = xfilepath.JoinFromHome(xfilepath.Path(".starport"))
)

// InitConfig creates config directory if it is not yet created
func InitConfig() error {
	confPath, err := StarportConfPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(confPath, 0755);
		return err
	}
	return nil
}
