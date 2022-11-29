package config

import "github.com/ignite/cli/ignite/pkg/xfilepath"

var (
	// DirPath returns the path of configuration directory of Ignite.
	DirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))
)
