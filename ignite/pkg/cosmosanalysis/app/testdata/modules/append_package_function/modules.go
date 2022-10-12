package app

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	foomodule "github.com/username/test/x/foo"
)

// App modules are defined in a function from a different
// file from where the variable is being referenced.
func basicModules() {
	return []module.AppModuleBasic{
		foomodule.AppModuleBasic{},
	}
}
