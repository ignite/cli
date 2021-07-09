package modulecreate

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xstrings"
	"github.com/tendermint/starport/starport/templates/module"
)

// NewOracle returns the generator to scaffold the implementation of the Oracle interface inside a module
func NewOracle(replacer placeholder.Replacer, opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(moduleOracleModify(replacer, opts))
	g.RunFn(eventOracleModify(replacer, opts))
	g.RunFn(protoQueryOracleModify(replacer, opts))
	g.RunFn(protoTxOracleModify(replacer, opts))
	g.RunFn(handlerTxOracleModify(replacer, opts))
	g.RunFn(clientCliQueryOracleModify(replacer, opts))
	g.RunFn(clientCliTxOracleModify(replacer, opts))
	g.RunFn(codecOracleModify(replacer, opts))

	if err := g.Box(oracleTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("title", strings.Title)

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func moduleOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module_ibc.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Recv packet dispatch
		templateRecv := `	
	ack, oracleResult, err := am.handleOraclePacket(ctx, modulePacket)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error())
	} else if oracleResult.Size() > 0 {
		ctx.Logger().Debug("Receive oracle packet", "result", oracleResult)
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, ack.GetBytes(), nil
	}`
		content := replacer.Replace(f.String(), module.PlaceholderOraclePacketModuleRecv, templateRecv)

		// Ack packet dispatch
		templateAck := `
	var requestID types.RequestID
	ctx, requestID = am.handleOracleAcknowledgement(ctx, ack)
	if requestID > 0 {
		ctx.Logger().Debug("Receive oracle ack", "request_id", requestID)
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, nil
	}`
		content = replacer.Replace(content, module.PlaceholderOraclePacketModuleAck, templateAck)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func eventOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/events_ibc.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%[1]v
EventTypeOraclePacket       = "oracle_packet"
`
		replacement := fmt.Sprintf(template, module.PlaceholderIBCPacketEvent)
		content := replacer.Replace(f.String(), module.PlaceholderIBCPacketEvent, replacement)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoQueryOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import the type
		templateImport := `%[1]v
import "%[2]v/oracle.proto";`
		replacementImport := fmt.Sprintf(templateImport, module.Placeholder, opts.ModuleName)
		content := replacer.Replace(f.String(), module.Placeholder, replacementImport)

		// Add the service
		templateService := `%[1]v

  	// Request defines a rpc handler method for MsgOracleData.
  	rpc OracleResult(QueryOracleRequest) returns (QueryOracleResponse) {
		option (google.api.http).get = "/%[2]v/%[3]v/result/{request_id}";
  	}

  	// LastOracleId
  	rpc LastOracleId(QueryLastOracleIdRequest) returns (QueryLastOracleIdResponse) {
		option (google.api.http).get = "/%[2]v/%[3]v/last_request_id";
  	}
`
		replacementService := fmt.Sprintf(templateService, module.Placeholder2,
			opts.AppName,
			opts.ModuleName,
		)
		content = replacer.Replace(content, module.Placeholder2, replacementService)

		// Add the service messages
		templateMessage := `%[1]v
message QueryOracleRequest {int64 request_id = 1;}

message QueryOracleResponse {
  OracleResult result = 1;
}

message QueryLastOracleIdRequest {}

message QueryLastOracleIdResponse {int64 request_id = 1;}
`
		replacementMessage := fmt.Sprintf(templateMessage, module.Placeholder3)
		content = replacer.Replace(content, module.Placeholder3, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/tx.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `%[1]v
import "gogoproto/gogo.proto";
import "cosmos/base/v1beta1/coin.proto";
import "%[2]v/oracle.proto";`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderProtoTxImport, opts.ModuleName)
		content := replacer.Replace(f.String(), module.PlaceholderProtoTxImport, replacementImport)

		// RPC
		templateRPC := `%[1]v
  rpc OracleData(MsgOracleData) returns (MsgOracleDataResponse);`
		replacementRPC := fmt.Sprintf(templateRPC, module.PlaceholderProtoTxRPC)
		content = replacer.Replace(content, module.PlaceholderProtoTxRPC, replacementRPC)

		templateMessage := `%[1]v
message MsgOracleData {
  string creator = 1;
  int64 oracle_script_id = 2 [
    (gogoproto.customname) = "OracleScriptID",
    (gogoproto.moretags) = "yaml:\"oracle_script_id\""
  ];
  string source_channel = 3;
  CallData calldata = 4;
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

message MsgOracleDataResponse {
}
`
		replacementMessage := fmt.Sprintf(templateMessage, module.PlaceholderProtoTxMessage)
		content = replacer.Replace(content, module.PlaceholderProtoTxMessage, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func handlerTxOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set once the MsgServer definition if it is not defined yet
		replacementMsgServer := `msgServer := keeper.NewMsgServerImpl(k)`
		content := replacer.ReplaceOnce(f.String(), module.PlaceholderHandlerMsgServer, replacementMsgServer)

		templateHandlers := `%[1]v
		case *types.MsgOracleData:
					res, err := msgServer.OracleData(sdk.WrapSDKContext(ctx), msg)
					return sdk.WrapServiceResult(ctx, res, err)
`
		replacementHandlers := fmt.Sprintf(templateHandlers, module.Placeholder)
		content = replacer.Replace(content, module.Placeholder, replacementHandlers)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliQueryOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v

	cmd.AddCommand(CmdOracleResult())
	cmd.AddCommand(CmdLastRequest())
`
		replacement := fmt.Sprintf(template, module.Placeholder)
		content := replacer.Replace(f.String(), module.Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliTxOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	cmd.AddCommand(CmdRequestOracleData())
`
		replacement := fmt.Sprintf(template, module.Placeholder)
		content := replacer.Replace(f.String(), module.Placeholder, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Set import if not set yet
		replacement := `sdk "github.com/cosmos/cosmos-sdk/types"`
		content := replacer.ReplaceOnce(f.String(), module.Placeholder, replacement)

		// Register the module packet
		templateRegistry := `%[1]v
cdc.RegisterConcrete(&MsgOracleData{}, "%[2]v/OracleData", nil)
`
		replacementRegistry := fmt.Sprintf(templateRegistry, module.Placeholder2, opts.ModuleName)
		content = replacer.Replace(content, module.Placeholder2, replacementRegistry)

		// Register the module packet interface
		templateInterface := `%[1]v
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgOracleData{},
)`
		replacementInterface := fmt.Sprintf(templateInterface, module.Placeholder3)
		content = replacer.Replace(content, module.Placeholder3, replacementInterface)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
