package ibc

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/ignite/templates/module"
	"github.com/ignite/cli/ignite/templates/testutil"
)

//go:embed files/oracle/* files/oracle/**/*
var fsOracle embed.FS

// OracleOptions are options to scaffold an oracle query in a IBC module.
type OracleOptions struct {
	AppName    string
	AppPath    string
	ModuleName string
	ModulePath string
	QueryName  multiformatname.Name
	MsgSigner  multiformatname.Name
}

// NewOracle returns the generator to scaffold the implementation of the Oracle interface inside a module
//
// Deprecated: This function is no longer maintained.
func NewOracle(replacer placeholder.Replacer, opts *OracleOptions) (*genny.Generator, error) {
	g := genny.New()

	template := xgenny.NewEmbedWalker(fsOracle, "files/oracle/", opts.AppPath)

	g.RunFn(moduleOracleModify(replacer, opts))
	g.RunFn(protoQueryOracleModify(replacer, opts))
	g.RunFn(protoTxOracleModify(replacer, opts))
	g.RunFn(clientCliQueryOracleModify(replacer, opts))
	g.RunFn(clientCliTxOracleModify(replacer, opts))
	g.RunFn(codecOracleModify(replacer, opts))

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("queryName", opts.QueryName)
	ctx.Set("MsgSigner", opts.MsgSigner)
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName))

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{queryName}}", opts.QueryName.Snake))

	// Create the 'testutil' package with the test helpers
	if err := testutil.Register(g, opts.AppPath); err != nil {
		return g, err
	}

	if err := xgenny.Box(g, template); err != nil {
		return g, err
	}

	g.RunFn(packetHandlerOracleModify(replacer, opts))

	return g, nil
}

// Deprecated: This function is no longer maintained.
func moduleOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module_ibc.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Recv packet dispatch
		templateRecv := `oracleAck, err := im.handleOraclePacket(ctx, modulePacket)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: "+err.Error()))
	} else if ack != oracleAck {
		return oracleAck
	}
	%[1]v`
		replacementRecv := fmt.Sprintf(templateRecv, PlaceholderOraclePacketModuleRecv)
		content := replacer.ReplaceOnce(f.String(), PlaceholderOraclePacketModuleRecv, replacementRecv)

		// Ack packet dispatch
		templateAck := `sdkResult, err := im.handleOracleAcknowledgment(ctx, ack, modulePacket)
	if err != nil {
		return err
	}
	if sdkResult != nil {
		sdkResult.Events = ctx.EventManager().Events().ToABCIEvents()
		return nil
	}
	%[1]v`
		replacementAck := fmt.Sprintf(templateAck, PlaceholderOraclePacketModuleAck)
		content = replacer.ReplaceOnce(content, PlaceholderOraclePacketModuleAck, replacementAck)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Deprecated: This function is no longer maintained.
func protoQueryOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.AppName, opts.ModuleName, "query.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import the type
		templateImport := `import "%[2]v/%[3]v/%[4]v.proto";
%[1]v`
		replacementImport := fmt.Sprintf(templateImport, Placeholder, opts.AppName, opts.ModuleName, opts.QueryName.Snake)
		content := replacer.Replace(f.String(), Placeholder, replacementImport)

		// Add the service
		templateService := `
  	// %[2]vResult defines a rpc handler method for Msg%[2]vData.
  	rpc %[2]vResult(Query%[2]vRequest) returns (Query%[2]vResponse) {
		option (google.api.http).get = "/%[3]v/%[4]v/%[5]v_result/{request_id}";
  	}

  	// Last%[2]vId query the last %[2]v result id
  	rpc Last%[2]vId(QueryLast%[2]vIdRequest) returns (QueryLast%[2]vIdResponse) {
		option (google.api.http).get = "/%[3]v/%[4]v/last_%[5]v_id";
  	}
%[1]v`
		replacementService := fmt.Sprintf(templateService, Placeholder2,
			opts.QueryName.UpperCamel,
			opts.AppName,
			opts.ModuleName,
			opts.QueryName.Snake,
		)
		content = replacer.Replace(content, Placeholder2, replacementService)

		// Add the service messages
		templateMessage := `message Query%[2]vRequest {int64 request_id = 1;}

message Query%[2]vResponse {
  %[2]vResult result = 1;
}

message QueryLast%[2]vIdRequest {}

message QueryLast%[2]vIdResponse {int64 request_id = 1;}

%[1]v`
		replacementMessage := fmt.Sprintf(templateMessage, Placeholder3, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, Placeholder3, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Deprecated: This function is no longer maintained.
func protoTxOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.AppName, opts.ModuleName, "tx.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := strings.ReplaceAll(f.String(), `
import "gogoproto/gogo.proto";`, "")
		content = strings.ReplaceAll(content, `
import "cosmos/base/v1beta1/coin.proto";`, "")

		// Import
		templateImport := `import "gogoproto/gogo.proto";
import "cosmos/base/v1beta1/coin.proto";
import "%[2]v/%[3]v/%[4]v.proto";
%[1]v`
		replacementImport := fmt.Sprintf(templateImport, PlaceholderProtoTxImport, opts.AppName, opts.ModuleName, opts.QueryName.Snake)
		content = replacer.Replace(content, PlaceholderProtoTxImport, replacementImport)

		// RPC
		templateRPC := `  rpc %[2]vData(Msg%[2]vData) returns (Msg%[2]vDataResponse);
%[1]v`
		replacementRPC := fmt.Sprintf(templateRPC, PlaceholderProtoTxRPC, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, PlaceholderProtoTxRPC, replacementRPC)

		templateMessage := `message Msg%[2]vData {
  string %[3]v = 1;
  uint64 oracle_script_id = 2 [
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
  uint64 prepare_gas = 8;
  uint64 execute_gas = 9;
  string client_id = 10 [(gogoproto.customname) = "ClientID"];
}

message Msg%[2]vDataResponse {
}

%[1]v`
		replacementMessage := fmt.Sprintf(templateMessage, PlaceholderProtoTxMessage,
			opts.QueryName.UpperCamel,
			opts.MsgSigner.LowerCamel,
		)
		content = replacer.Replace(content, PlaceholderProtoTxMessage, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Deprecated: This function is no longer maintained.
func clientCliQueryOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/query.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `
	cmd.AddCommand(Cmd%[2]vResult())
	cmd.AddCommand(CmdLast%[2]vID())
%[1]v`
		replacement := fmt.Sprintf(template, Placeholder, opts.QueryName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Deprecated: This function is no longer maintained.
func clientCliTxOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/tx.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `cmd.AddCommand(CmdRequest%[2]vData())
%[1]v`
		replacement := fmt.Sprintf(template, Placeholder, opts.QueryName.UpperCamel)
		content := replacer.Replace(f.String(), Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Deprecated: This function is no longer maintained.
func codecOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set import if not set yet
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), Placeholder, replacement)

		// Register the module packet
		templateRegistry := `cdc.RegisterConcrete(&Msg%[3]vData{}, "%[2]v/%[3]vData", nil)
%[1]v`
		replacementRegistry := fmt.Sprintf(templateRegistry, Placeholder2, opts.ModuleName, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, Placeholder2, replacementRegistry)

		// Register the module packet interface
		templateInterface := `registry.RegisterImplementations((*sdk.Msg)(nil),
	&Msg%[2]vData{},
)
%[1]v`
		replacementInterface := fmt.Sprintf(templateInterface, Placeholder3, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, Placeholder3, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Deprecated: This function is no longer maintained.
func packetHandlerOracleModify(replacer placeholder.Replacer, opts *OracleOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "oracle.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Register the module packet
		templateRecv := `
	case types.%[3]vClientIDKey:
		var %[2]vResult types.%[3]vResult
		if err := obi.Decode(modulePacketData.Result, &%[2]vResult); err != nil {
			ack = channeltypes.NewErrorAcknowledgement(err)
			return ack, errorsmod.Wrap(sdkerrors.ErrUnknownRequest,
				"cannot decode the %[2]v received packet")
		}
		im.keeper.Set%[3]vResult(ctx, types.OracleRequestID(modulePacketData.RequestID), %[2]vResult)
	
		// TODO: %[3]v oracle data reception logic
%[1]v`
		replacementRegistry := fmt.Sprintf(templateRecv, PlaceholderOracleModuleRecv,
			opts.QueryName.LowerCamel, opts.QueryName.UpperCamel)
		content := replacer.Replace(f.String(), PlaceholderOracleModuleRecv, replacementRegistry)

		// Register the module packet interface
		templateAck := `
	case types.%[3]vClientIDKey:
		var %[2]vData types.%[3]vCallData
		if err = obi.Decode(data.GetCalldata(), &%[2]vData); err != nil {
			return nil, errorsmod.Wrap(err,
				"cannot decode the %[2]v oracle acknowledgment packet")
		}
		im.keeper.SetLast%[3]vID(ctx, requestID)
		return &sdk.Result{}, nil
%[1]v`
		replacementInterface := fmt.Sprintf(templateAck, PlaceholderOracleModuleAck,
			opts.QueryName.LowerCamel, opts.QueryName.UpperCamel)
		content = replacer.Replace(content, PlaceholderOracleModuleAck, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
