package modules

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
	foomodule "github.com/username/test/x/foo"
)

// App modules are defined in a variable from a different file
// and package from where the variable is being referenced.
var Basic = []module.AppModuleBasic{
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	gov.NewAppModuleBasic([]govclient.ProposalHandler{
		paramsclient.ProposalHandler,
	}),
	foomodule.AppModuleBasic{},
}
