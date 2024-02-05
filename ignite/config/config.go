package config

import (
	"github.com/ignite/cli/v29/ignite/pkg/env"
	"github.com/ignite/cli/v29/ignite/pkg/xfilepath"
)

// DirPath returns the path of configuration directory of Ignite.
var DirPath = xfilepath.Mkdir(env.ConfigDir())
