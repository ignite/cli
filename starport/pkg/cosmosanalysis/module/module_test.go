package module

import (
	"github.com/kr/pretty"
)

func ExampleDiscover() {
	pretty.Println(Discover("/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon"))
	// outputs:
	// 	[]msgs.Module{
	// 		{
	// 			Name:            "moon",
	// 			TypesImportPath: "github.com/ilker/moon/x/moon/types",
	// 			Msgs:            {
	// 				{Name:"MsgUpdateUser", URI:"ilker.moon.moon.MsgUpdateUser"},
	// 				{Name:"MsgDeleteUser", URI:"ilker.moon.moon.MsgDeleteUser"},
	// 				{Name:"MsgCreateUser", URI:"ilker.moon.moon.MsgCreateUser"},
	// 			},
	// 		},
	// 		{
	// 			Name:            "elips",
	// 			TypesImportPath: "github.com/ilker/moon/x/elips/types",
	// 			Msgs:            {
	// 				{Name:"MsgCreateBlog", URI:"ilker.moon.elips.MsgCreateBlog"},
	// 				{Name:"MsgUpdateBlog", URI:"ilker.moon.elips.MsgUpdateBlog"},
	// 				{Name:"MsgDeleteBlog", URI:"ilker.moon.elips.MsgDeleteBlog"},
	// 			},
	// 		},
	// 	} nil
}
