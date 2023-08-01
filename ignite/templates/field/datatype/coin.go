package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
)

var (
	// DataCoin coin data type definition.
	DataCoin = DataType{
		DataType:         func(string) string { return "sdk.Coin" },
		DefaultTestValue: "10token",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("cosmos.base.v1beta1.Coin %s = %d [(gogoproto.nullable) = false]",
				name, index)
		},
		GenesisArgs: func(multiformatname.Name, int) string { return "" },
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := sdk.ParseCoinNormalized(args[%d])
					if err != nil {
						return err
					}`, prefix, name.UpperCamel, argIndex)
		},
		GoCLIImports: []GoImport{{Name: "github.com/cosmos/cosmos-sdk/types", Alias: "sdk"}},
		ProtoImports: []string{"gogoproto/gogo.proto", "cosmos/base/v1beta1/coin.proto"},
		NonIndex:     true,
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			option := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
			return protoutil.NewField(
				name, "cosmos.base.v1beta1.Coin", index, protoutil.WithFieldOptions(option),
			)
		},
	}

	// DataCoinSlice is a coin array data type definition.
	DataCoinSlice = DataType{
		DataType:         func(string) string { return "sdk.Coins" },
		DefaultTestValue: "10token,20stake",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("repeated cosmos.base.v1beta1.Coin %s = %d [(gogoproto.nullable) = false]",
				name, index)
		},
		GenesisArgs: func(multiformatname.Name, int) string { return "" },
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := sdk.ParseCoinsNormalized(args[%d])
					if err != nil {
						return err
					}`, prefix, name.UpperCamel, argIndex)
		},
		GoCLIImports: []GoImport{{Name: "github.com/cosmos/cosmos-sdk/types", Alias: "sdk"}},
		ProtoImports: []string{"gogoproto/gogo.proto", "cosmos/base/v1beta1/coin.proto"},
		NonIndex:     true,
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			option := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
			return protoutil.NewField(
				name, "cosmos.base.v1beta1.Coin", index, protoutil.WithFieldOptions(option), protoutil.Repeated(),
			)
		},
	}
)
