package datatype

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
)

var (
	// DataCoin coin data type definition.
	DataCoin = DataType{
		Name:                    Coin,
		DataType:                func(string) string { return "sdk.Coin" },
		CollectionsKeyValueName: func(string) string { return collectionValueComment },
		DefaultTestValue:        "10token",
		ValueLoop:               "sdk.NewInt64Coin(`token`, int64(i+100))",
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
		Name:                    Coins,
		DataType:                func(string) string { return "sdk.Coins" },
		CollectionsKeyValueName: func(string) string { return collectionValueComment },
		DefaultTestValue:        "20stake",
		ValueLoop:               "sdk.NewCoins(sdk.NewInt64Coin(`token`, int64(i%1+100)), sdk.NewInt64Coin(`stake`, int64(i%2+100)))",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf(`repeated cosmos.base.v1beta1.Coin %s = %d [(gogoproto.nullable) = false]`,
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

	// DataDecCoin decimal coin data type definition.
	DataDecCoin = DataType{
		Name:                    DecCoin,
		DataType:                func(string) string { return "sdk.DecCoin" },
		CollectionsKeyValueName: func(string) string { return collectionValueComment },
		DefaultTestValue:        "100001token",
		ValueLoop:               "sdk.NewInt64DecCoin(`token`, int64(i+100))",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf("cosmos.base.v1beta1.DecCoin %s = %d [(gogoproto.nullable) = false]",
				name, index)
		},
		GenesisArgs: func(multiformatname.Name, int) string { return "" },
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := sdk.ParseDecCoins(args[%d])
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
				name, "cosmos.base.v1beta1.DecCoin", index, protoutil.WithFieldOptions(option),
			)
		},
	}

	// DataDecCoinSlice is a decimal coin array data type definition.
	DataDecCoinSlice = DataType{
		Name:                    DecCoins,
		DataType:                func(string) string { return "sdk.DecCoins" },
		CollectionsKeyValueName: func(string) string { return collectionValueComment },
		DefaultTestValue:        "20000002stake",
		ValueLoop:               "sdk.NewDecCoins(sdk.NewInt64DecCoin(`token`, int64(i%1+100)), sdk.NewInt64DecCoin(`stake`, int64(i%2+100)))",
		ProtoType: func(_, name string, index int) string {
			return fmt.Sprintf(`repeated cosmos.base.v1beta1.DecCoin %s = %d [(gogoproto.nullable) = false]`,
				name, index)
		},
		GenesisArgs: func(multiformatname.Name, int) string { return "" },
		CLIArgs: func(name multiformatname.Name, _, prefix string, argIndex int) string {
			return fmt.Sprintf(`%s%s, err := sdk.ParseDecCoins(args[%d])
					if err != nil {
						return err
					}`, prefix, name.UpperCamel, argIndex)
		},
		GoCLIImports: []GoImport{{Name: "github.com/cosmos/cosmos-sdk/types", Alias: "sdk"}},
		ProtoImports: []string{"gogoproto/gogo.proto", "cosmos/base/v1beta1/coin.proto"},
		NonIndex:     true,
		ToProtoField: func(_, name string, index int) *proto.NormalField {
			optionNullable := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
			optionCast := protoutil.NewOption("gogoproto.castrepeated", "github.com/cosmos/cosmos-sdk/types.DecCoins", protoutil.Custom())
			return protoutil.NewField(name, "cosmos.base.v1beta1.DecCoin", index,
				protoutil.WithFieldOptions(optionNullable),
				protoutil.WithFieldOptions(optionCast),
				protoutil.Repeated(),
			)
		},
	}
)
