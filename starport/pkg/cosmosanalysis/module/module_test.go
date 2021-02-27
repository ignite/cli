package module

import (
	"github.com/kr/pretty"
)

func ExampleDiscover() {
	pretty.Println(Discover("/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon"))
	// outputs:
	// []module.Module{
	// 	{
	// 		Name: "moon",
	// 		Pkg:  protoanalysis.Package{
	// 			Name:         "test.moon.moon",
	// 			Path:         "/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon",
	// 			GoImportName: "github.com/test/moon/x/moon/types",
	// 			Messages:     {
	// 				{Name:"GenesisState", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/genesis.proto"},
	// 				{Name:"GenesisState", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/genesis.proto"},
	// 				{Name:"QueryGetUserRequest", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/query.proto"},
	// 				{Name:"QueryGetUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/query.proto"},
	// 				{Name:"QueryAllUserRequest", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/query.proto"},
	// 				{Name:"QueryAllUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/query.proto"},
	// 				{Name:"MsgCreateUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 				{Name:"MsgCreateUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 				{Name:"MsgUpdateUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 				{Name:"MsgUpdateUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 				{Name:"MsgDeleteUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 				{Name:"MsgDeleteUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 				{Name:"User", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/user.proto"},
	// 			},
	// 		},
	// 		Msgs: {
	// 			{Name:"MsgUpdateUser", URI:"test.moon.moon.MsgUpdateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 			{Name:"MsgDeleteUser", URI:"test.moon.moon.MsgDeleteUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 			{Name:"MsgCreateUser", URI:"test.moon.moon.MsgCreateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/moon/proto/moon/tx.proto"},
	// 		},
	// 	},
	// } nil
}
