package module

import (
	"github.com/kr/pretty"
)

func ExampleDiscover() {
	pretty.Println(Discover("/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars"))
	// outputs:
	// []module.Module{
	//    {
	//        Name: "mars",
	//        Pkg:  protoanalysis.Package{
	//            Name:         "tendermint.mars.mars",
	//            Path:         "/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars",
	//            GoImportName: "github.com/tendermint/mars/x/mars/types",
	//            Messages:     {
	//                {Name:"GenesisState", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/genesis.proto"},
	//                {Name:"User", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/user.proto"},
	//                {Name:"MsgCreateUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//                {Name:"MsgCreateUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//                {Name:"MsgUpdateUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//                {Name:"MsgUpdateUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//                {Name:"MsgDeleteUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//                {Name:"MsgDeleteUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//                {Name:"QueryGetUserRequest", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//                {Name:"QueryGetUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//                {Name:"QueryAllUserRequest", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//                {Name:"QueryAllUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//            },
	//            Services: {
	//                {
	//                    Name:     "Msg",
	//                    RPCFuncs: {
	//                        {Name:"CreateUser", RequestType:"MsgCreateUser", ReturnsType:"MsgCreateUserResponse"},
	//                        {Name:"UpdateUser", RequestType:"MsgUpdateUser", ReturnsType:"MsgUpdateUserResponse"},
	//                        {Name:"DeleteUser", RequestType:"MsgDeleteUser", ReturnsType:"MsgDeleteUserResponse"},
	//                    },
	//                },
	//                {
	//                    Name:     "Query",
	//                    RPCFuncs: {
	//                        {Name:"User", RequestType:"QueryGetUserRequest", ReturnsType:"QueryGetUserResponse"},
	//                        {Name:"UserAll", RequestType:"QueryAllUserRequest", ReturnsType:"QueryAllUserResponse"},
	//                    },
	//                },
	//            },
	//        },
	//        Msgs: {
	//            {Name:"MsgUpdateUser", URI:"tendermint.mars.mars.MsgUpdateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//            {Name:"MsgDeleteUser", URI:"tendermint.mars.mars.MsgDeleteUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//            {Name:"MsgCreateUser", URI:"tendermint.mars.mars.MsgCreateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//        },
	//        Queries: {
	//            {Name:"User", FullName:"QueryUser"},
	//            {Name:"UserAll", FullName:"QueryUserAll"},
	//        },
	//        Types: {
	//            {Name:"User", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/user.proto"},
	//        },
	//    },
	// } nil
}
