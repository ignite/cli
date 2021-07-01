package module

// messageImplementation is the list of methods needed for a sdk.Msg implementation
// TODO(low priority): dynamically get these from the source code of underlying version of the sdk.
var messageImplementation = []string{
	"Route",
	"Type",
	"GetSigners",
	"GetSignBytes",
	"ValidateBasic",
}
