package indexedlist

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/typed"
)

var (
	//go:embed component/* component/**/*
	fsComponent embed.FS

	//go:embed messages/* messages/**/*
	fsMessages embed.FS
)

// NewStargate returns the generator to scaffold a new indexed list in a Stargate module
func NewStargate(replacer placeholder.Replacer, opts *typed.Options) (*genny.Generator, error) {
	var (
		g = genny.New()

		messagesTemplate = xgenny.NewEmbedWalker(
			fsMessages,
			"messages/",
			opts.AppPath,
		)
		componentTemplate = xgenny.NewEmbedWalker(
			fsComponent,
			"component/",
			opts.AppPath,
		)
	)

	g.RunFn(typed.ModuleGRPCGatewayModify(replacer, opts))
	g.RunFn(typed.ClientCliQueryModify(replacer, opts))
	g.RunFn(protoRPCModify(replacer, opts))

	// Genesis modifications
	genesisModify(replacer, opts, g)

	if !opts.NoMessage {
		g.RunFn(protoTxModify(replacer, opts))
		g.RunFn(typed.HandlerModify(replacer, opts))
		g.RunFn(typed.TypesCodecModify(replacer, opts))
		g.RunFn(typed.ClientCliTxModify(replacer, opts))

		// Messages template
		if err := typed.Box(messagesTemplate, opts, g); err != nil {
			return nil, err
		}
	}

	return g, typed.Box(componentTemplate, opts, g)
}

func protoRPCModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "query.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import the type
		templateImport := `import "%s/%s.proto";
%s`
		replacementImport := fmt.Sprintf(templateImport,
			opts.ModuleName,
			opts.TypeName.Snake,
			typed.Placeholder,
		)
		content := replacer.Replace(f.String(), typed.Placeholder, replacementImport)

		// Add gogo.proto
		replacementGogoImport := typed.EnsureGogoProtoImported(path, typed.Placeholder)
		content = replacer.Replace(content, typed.Placeholder, replacementGogoImport)

		var lowerCamelIndexes []string
		for _, index := range opts.Indexes {
			lowerCamelIndexes = append(lowerCamelIndexes, fmt.Sprintf("{%s}", index.Name.LowerCamel))
		}
		indexPath := strings.Join(lowerCamelIndexes, "/")

		// Add the service
		templateService := `// Queries a %[3]v by index.
	rpc %[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v/%[7]v/{id}";
	}

	// Queries a list of %[3]v items.
	rpc %[2]vAll(QueryAll%[2]vRequest) returns (QueryAll%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v/%[7]v";
	}

%[1]v`
		replacementService := fmt.Sprintf(templateService, typed.Placeholder2,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
			indexPath,
		)
		content = replacer.Replace(content, typed.Placeholder2, replacementService)

		// Add the service messages
		var queryIndexFields string
		for i, index := range opts.Indexes {
			queryIndexFields += fmt.Sprintf("  %s\n", index.ProtoType(i+1))
		}

		// Ensure custom types are imported
		protoImports := opts.Fields.ProtoImports()
		for _, f := range opts.Fields.Custom() {
			protoImports = append(protoImports,
				fmt.Sprintf("%[1]v/%[2]v.proto", opts.ModuleName, f),
			)
		}
		for _, f := range protoImports {
			importModule := fmt.Sprintf(`
import "%[1]v";`, f)
			content = strings.ReplaceAll(content, importModule, "")
			replacementImport := fmt.Sprintf("%[1]v%[2]v", typed.Placeholder, importModule)
			content = replacer.Replace(content, typed.Placeholder, replacementImport)
		}

		templateMessage := `message QueryGet%[2]vRequest {
  %[4]v  uint64 id = %[5]v;
}

message QueryGet%[2]vResponse {
  %[2]v %[3]v = 1 [(gogoproto.nullable) = false];
}

message QueryAll%[2]vRequest {
  %[4]v  cosmos.base.query.v1beta1.PageRequest pagination = %[5]v;
}

message QueryAll%[2]vResponse {
  repeated %[2]v %[3]v = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}

%[1]v`
		replacementMessage := fmt.Sprintf(templateMessage,
			typed.Placeholder3,
			opts.TypeName.UpperCamel,
			opts.TypeName.LowerCamel,
			queryIndexFields,
			len(queryIndexFields)+1,
		)
		content = replacer.Replace(content, typed.Placeholder3, replacementMessage)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoTxModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "tx.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `import "%s/%s.proto";
%s`
		replacementImport := fmt.Sprintf(templateImport,
			opts.ModuleName,
			opts.TypeName.Snake,
			typed.PlaceholderProtoTxImport,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderProtoTxImport, replacementImport)

		// RPC service
		templateRPC := `rpc Create%[2]v(MsgCreate%[2]v) returns (MsgCreate%[2]vResponse);
  rpc Update%[2]v(MsgUpdate%[2]v) returns (MsgUpdate%[2]vResponse);
  rpc Delete%[2]v(MsgDelete%[2]v) returns (MsgDelete%[2]vResponse);
  %[1]v`
		replacementRPC := fmt.Sprintf(templateRPC, typed.PlaceholderProtoTxRPC,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderProtoTxRPC, replacementRPC)

		// Messages
		var indexes string
		for i, index := range opts.Indexes {
			indexes += fmt.Sprintf("  %s\n", index.ProtoType(i+2))
		}

		var fieldsCreateMsg string
		for i, f := range opts.Fields {
			fieldsCreateMsg += fmt.Sprintf("  %s\n", f.ProtoType(i+2+len(opts.Indexes)))
		}

		var fieldsUpdateMsg string
		for i, f := range opts.Fields {
			fieldsUpdateMsg += fmt.Sprintf("  %s\n", f.ProtoType(i+3+len(opts.Indexes)))
		}

		// Ensure custom types are imported
		protoImports := append(opts.Fields.ProtoImports(), opts.Indexes.ProtoImports()...)
		customFields := append(opts.Fields.Custom(), opts.Indexes.Custom()...)
		for _, f := range customFields {
			protoImports = append(protoImports,
				fmt.Sprintf("%[1]v/%[2]v.proto", opts.ModuleName, f),
			)
		}
		for _, f := range protoImports {
			importModule := fmt.Sprintf(`
import "%[1]v";`, f)
			content = strings.ReplaceAll(content, importModule, "")

			replacementImport := fmt.Sprintf("%[1]v%[2]v", typed.PlaceholderProtoTxImport, importModule)
			content = replacer.Replace(content, typed.PlaceholderProtoTxImport, replacementImport)
		}

		templateMessages := `message MsgCreate%[2]v {
  string %[3]v = 1;
%[4]v%[5]v}
message MsgCreate%[2]vResponse {
  uint64 id = 1;
}

message MsgUpdate%[2]v {
  string %[3]v = 1;
%[4]v  uint64 id = %[6]v;
%[7]v}

message MsgUpdate%[2]vResponse {}

message MsgDelete%[2]v {
  string %[3]v = 1;
  %[4]v  uint64 id = %[6]v;
}
message MsgDelete%[2]vResponse {}

%[1]v`
		replacementMessages := fmt.Sprintf(
			templateMessages,
			typed.PlaceholderProtoTxMessage,
			opts.TypeName.UpperCamel,
			opts.MsgSigner.LowerCamel,
			indexes,
			fieldsCreateMsg,
			len(opts.Indexes)+2,
			fieldsUpdateMsg,
		)
		content = replacer.Replace(content, typed.PlaceholderProtoTxMessage, replacementMessages)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
