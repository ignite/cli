package config

import (
	"github.com/ignite/cli/ignite/pkg/env"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

// DirPath returns the path of configuration directory of Ignite.
var DirPath = xfilepath.Mkdir(env.ConfigDir())
