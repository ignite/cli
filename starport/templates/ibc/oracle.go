package ibc

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/pkg/xstrings"
)

var (
	//go:embed oracle/static/* oracle/static/**/*
	fsOracleStatic embed.FS

	// ibcOracleStaticTemplate is the template to scaffold a new static oracle templates in an IBC module
	ibcOracleStaticTemplate = xgenny.NewEmbedWalker(fsOracleStatic, "oracle/static/")

	//go:embed oracle/dynamic/* oracle/dynamic/**/*
	fsOracleDynamic embed.FS

	// ibcOracleDynamicTemplate is the template to scaffold a new dynamic oracle templates in an IBC module
	ibcOracleDynamicTemplate = xgenny.NewEmbedWalker(fsOracleDynamic, "oracle/dynamic/")
)

// OracleOptions are options to scaffold an oracle in a IBC module
type OracleOptions struct {
	AppName    string
	ModuleName string
	ModulePath string
	OwnerName  string
	QueryName  multiformatname.Name
}

// NewOracle returns the generator to scaffold the implementation of the Oracle interface inside a module
func NewOracle(replacer placeholder.Replacer, opts *OracleOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(moduleOracleModify(replacer, opts))
	g.RunFn(protoQueryOracleModify(replacer, opts))
	g.RunFn(protoTxOracleModify(replacer, opts))
	g.RunFn(handlerTxOracleModify(replacer, opts))
	g.RunFn(clientCliQueryOracleModify(replacer, opts))
	g.RunFn(clientCliTxOracleModify(replacer, opts))
	g.RunFn(codecOracleModify(replacer, opts))

	err := box(g, opts)
	if err != nil {
		return g, err
	}
	g.RunFn(packetHandlerOracleModify(replacer, opts))

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("queryName", opts.QueryName)
	ctx.Set("title", strings.Title)

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{queryName}}", opts.QueryName.Snake))
	return g, nil
}

func box(g *genny.Generator, opts *OracleOptions) error {
	gs := genny.New()
	path := fmt.Sprintf("x/%s/oracle.go", opts.ModuleName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := gs.Box(ibcOracleStaticTemplate); err != nil {
			return err
		}
	}
	g.Merge(gs)
	return g.Box(ibcOracleDynamicTemplate)
}

func moduleOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module_ibc.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Recv packet dispatch
		templateRecv := `%[1]v
	oracleAck, err := am.handleOraclePacket(ctx, modulePacket)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: "+err.Error())
	} else if ack != oracleAck {
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, oracleAck.GetBytes(), nil
	}`
		replacementRecv := fmt.Sprintf(templateRecv, PlaceholderOraclePacketModuleRecv)
		content := replacer.ReplaceOnce(f.String(), PlaceholderOraclePacketModuleRecv, replacementRecv)

		// Ack packet dispatch
		templateAck := `%[1]v
	sdkResult := am.handleOracleAcknowledgment(ctx, ack, modulePacket)
	if sdkResult != nil {
		sdkResult.Events = ctx.EventManager().Events().ToABCIEvents()
		return sdkResult, nil
	}`
		replacementAck := fmt.Sprintf(templateAck, PlaceholderOraclePacketModuleAck)
		content = replacer.ReplaceOnce(content, PlaceholderOraclePacketModuleAck, replacementAck)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoQueryOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import the type
		templateImport := `%[1]v
import "%[2]v/%[3]v.proto";`
		replacementImport := fmt.Sprintf(templateImport, Placeholder, opts.ModuleName, opts.QueryName.Snake)
		content := replacer.Replace(f.String(), Placeholder, replacementImport)

		// Add the service
		templateService := `%[1]v

  	// %[2]vResult defines a rpc handler method for Msg%[2]vData.
  	rpc %[2]vResult(Query%[2]vRequest) returns (Query%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[3]v_result/{request_id}";
  	}

  	// Last%[2]vId query the last %[2]v result id
  	rpc Last%[2]vId(QueryLast%[2]vIdRequest) returns (QueryLast%[2]vIdResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/last_%[3]v_id";
  	}
`
		replacementService := fmt.Sprintf(templateService, Placeholder2,
			opts.QueryName.UpperCamel,
			opts.QueryName.Snake,
			opts.AppName,
			opts.ModuleName,
		)
		content = replacer.Replace(content, Placeholder2, replacementService)

		// Add the service messages
		templateMessage := `%[1]v
message Query%[2]vRequest {int64 request_id = 1;}

message Query%[2]vResponse {
  %[2]vResult result = 1;
}

message QueryLast%[2]vIdRequest {}

message QueryLast%[2]vIdResponse {int64 request_id = 1;}
`
		replacementMessage := fmt.Sprintf(templateMessage, Placeholder3, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, Placeholder3, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := strings.ReplaceAll(f.String(), `
import "gogoproto/gogo.proto";`, "")
		content = strings.ReplaceAll(content, `
import "cosmos/base/v1beta1/coin.proto";`, "")

		// Import
		templateImport := `%[1]v
import "gogoproto/gogo.proto";
import "cosmos/base/v1beta1/coin.proto";
import "%[2]v/%[3]v.proto";`
		replacementImport := fmt.Sprintf(templateImport, PlaceholderProtoTxImport, opts.ModuleName, opts.QueryName.Snake)
		content = replacer.Replace(content, PlaceholderProtoTxImport, replacementImport)

		// RPC
		templateRPC := `%[1]v
  rpc %[2]vData(Msg%[2]vData) returns (Msg%[2]vDataResponse);`
		replacementRPC := fmt.Sprintf(templateRPC, PlaceholderProtoTxRPC, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, PlaceholderProtoTxRPC, replacementRPC)

		templateMessage := `%[1]v
message Msg%[2]vData {
  string creator = 1;
  int64 oracle_script_id = 2 [
    (gogoproto.customname) = "OracleScriptID",
    (gogoproto.moretags) = "yaml:\"oracle_script_id\""
  ];
  string source_channel = 3;
  %[2]vCallData calldata = 4;
  uint64 ask_count = 5;
  uint64 min_count = 6;
  repeated cosmos.base.v1beta1.Coin fee_limit = 7 [
    (gogoproto.nullable) = false,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
  string request_key = 8;
  uint64 prepare_gas = 9;
  uint64 execute_gas = 10;
}

message Msg%[2]vDataResponse {
}
`
		replacementMessage := fmt.Sprintf(templateMessage, PlaceholderProtoTxMessage, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, PlaceholderProtoTxMessage, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func handlerTxOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := replacer.ReplaceOnce(f.String(), PlaceholderHandlerMsgServer, replacementMsgServer)

		templateHandlers := `%[1]v
		case *types.Msg%[2]vData:
					res, err := msgServer.%[2]vData(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
`
		replacementHandlers := fmt.Sprintf(templateHandlers, Placeholder, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, Placeholder, replacementHandlers)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliQueryOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v

	cmd.AddCommand(Cmd%[2]vResult())
	cmd.AddCommand(CmdLast%[2]vID())
`
		replacement := fmt.Sprintf(template, Placeholder, opts.QueryName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	cmd.AddCommand(CmdRequest%[2]vData())
`
		replacement := fmt.Sprintf(template, Placeholder, opts.QueryName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set import if not set yet
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), Placeholder, replacement)

		// Register the module packet
		templateRegistry := `%[1]v
cdc.RegisterConcrete(&Msg%[3]vData{}, "%[2]v/%[3]vData", nil)
`
		replacementRegistry := fmt.Sprintf(templateRegistry, Placeholder2, opts.ModuleName, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, Placeholder2, replacementRegistry)

		// Register the module packet interface
		templateInterface := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&Msg%[2]vData{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, Placeholder3, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, Placeholder3, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func packetHandlerOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/oracle.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Register the module packet
		templateRecv := `%[1]v
	var %[2]vResult types.%[3]vResult
	if err := obi.Decode(modulePacketData.Result, &%[2]vResult); err == nil {
		am.keeper.Set%[3]vResult(ctx, types.OracleRequestID(modulePacketData.RequestID), %[2]vResult)
		ack = channeltypes.NewResultAcknowledgement(
			types.ModuleCdc.MustMarshalJSON(
				packet.NewOracleRequestPacketAcknowledgement(modulePacketData.RequestID),
			),
		)

		// TODO: %[3]v oracle data reception logic
	}
`
		replacementRegistry := fmt.Sprintf(templateRecv, PlaceholderOracleModuleRecv,
			opts.QueryName.LowerCamel, opts.QueryName.UpperCamel)
		content := replacer.Replace(f.String(), PlaceholderOracleModuleRecv, replacementRegistry)

		// Register the module packet interface
		templateAck := `%[1]v
		var %[2]vData types.%[3]vCallData
		if err = obi.Decode(data.GetCalldata(), &%[2]vData); err == nil {
			am.keeper.SetLast%[3]vID(ctx, requestID)
			return &sdk.Result{}
		}
`
		replacementInterface := fmt.Sprintf(templateAck, PlaceholderOracleModuleAck,
			opts.QueryName.LowerCamel, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, PlaceholderOracleModuleAck, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
