package module

import (
	"github.com/kr/pretty"
)

func ExampleDiscover() {
	pretty.Println(Discover("/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars"))
	// outputs:
	// []module.Module{
	//   {
	//       Name: "mars",
	//       Pkg:  protoanalysis.Package{
	//           Name:         "test.mars.mars",
	//           Path:         "/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars",
	//           GoImportName: "github.com/test/mars/x/mars/types",
	//           Messages:     {
	//               {Name:"GenesisState", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/genesis.proto"},
	//               {Name:"QueryGetUserRequest", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//               {Name:"QueryGetUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//               {Name:"QueryAllUserRequest", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//               {Name:"QueryAllUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/query.proto"},
	//               {Name:"User", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/user.proto"},
	//               {Name:"MsgCreateUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//               {Name:"MsgCreateUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//               {Name:"MsgUpdateUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//               {Name:"MsgUpdateUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//               {Name:"MsgDeleteUser", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//               {Name:"MsgDeleteUserResponse", Path:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//           },
	//           Services: {
	//               {
	//                   Name:     "Query",
	//                   RPCFuncs: {
	//                       {
	//                           Name:            "User",
	//                           RequestType:     "QueryGetUserRequest",
	//                           ReturnsType:     "QueryGetUserResponse",
	//                           HTTPAnnotations: protoanalysis.HTTPAnnotations{
	//                               URLParams:   {"id"},
	//                               URLHasQuery: false,
	//                           },
	//                       },
	//                       {
	//                           Name:            "UserAll",
	//                           RequestType:     "QueryAllUserRequest",
	//                           ReturnsType:     "QueryAllUserResponse",
	//                           HTTPAnnotations: protoanalysis.HTTPAnnotations{
	//                               URLParams:   nil,
	//                               URLHasQuery: true,
	//                           },
	//                       },
	//                   },
	//               },
	//                {
	//                   Name:     "Msg",
	//                   RPCFuncs: {
	//                       {
	//                           Name:            "CreateUser",
	//                           RequestType:     "MsgCreateUser",
	//                           ReturnsType:     "MsgCreateUserResponse",
	//                           HTTPAnnotations: protoanalysis.HTTPAnnotations{},
	//                       },
	//                       {
	//                           Name:            "UpdateUser",
	//                           RequestType:     "MsgUpdateUser",
	//                           ReturnsType:     "MsgUpdateUserResponse",
	//                           HTTPAnnotations: protoanalysis.HTTPAnnotations{},
	//                       },
	//                       {
	//                           Name:            "DeleteUser",
	//                           RequestType:     "MsgDeleteUser",
	//                           ReturnsType:     "MsgDeleteUserResponse",
	//                           HTTPAnnotations: protoanalysis.HTTPAnnotations{},
	//                       },
	//                   },
	//               },
	//           },
	//       },
	//       Msgs: {
	//           {Name:"MsgUpdateUser", URI:"test.mars.mars.MsgUpdateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//           {Name:"MsgCreateUser", URI:"test.mars.mars.MsgCreateUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//           {Name:"MsgDeleteUser", URI:"test.mars.mars.MsgDeleteUser", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/tx.proto"},
	//       },
	//       Queries: {
	//           {
	//               Name:            "User",
	//               FullName:        "QueryUser",
	//               HTTPAnnotations: protoanalysis.HTTPAnnotations{
	//                   URLParams:   {"id"},
	//                   URLHasQuery: false,
	//               },
	//           },
	//           {
	//               Name:            "UserAll",
	//               FullName:        "QueryUserAll",
	//               HTTPAnnotations: protoanalysis.HTTPAnnotations{
	//                   URLParams:   nil,
	//                   URLHasQuery: true,
	//               },
	//           },
	//       },
	//       Types: {
	//           {Name:"User", FilePath:"/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars/proto/mars/user.proto"},
	//        },
	//    },
	// } nil

}
