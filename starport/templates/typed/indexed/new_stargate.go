package indexed

import (
	"embed"
	"fmt"

	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/typed"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	stargateIndexedTemplate = xgenny.NewEmbedWalker(fsStargate, "stargate/")
)

// New ...
func NewStargate(opts *typed.Options) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(typesKeyModify(opts))
	g.RunFn(protoRPCModify(opts))
	g.RunFn(moduleGRPCGatewayModify(opts))
	g.RunFn(clientCliQueryModify(opts))
	g.RunFn(genesisProtoModify(opts))
	g.RunFn(genesisTypesModify(opts))
	g.RunFn(genesisModuleModify(opts))

	return g, typed.Box(stargateIndexedTemplate, opts, g)
}

func typesKeyModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/keys.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
const (
	%[1]vKey= "%[1]v-value-"
)
`, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func protoRPCModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/query.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import the type
		templateImport := `%s
import "%s/%s.proto";`
		replacementImport := fmt.Sprintf(templateImport, typed.Placeholder,
			opts.ModuleName,
			opts.TypeName,
		)
		content := strings.Replace(f.String(), typed.Placeholder, replacementImport, 1)

		// Add the service
		templateService := `%[1]v
	rpc %[2]v(QueryGet%[2]vRequest) returns (QueryGet%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v/{index}";
	}
	rpc %[2]vAll(QueryAll%[2]vRequest) returns (QueryAll%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v";
	}
`
		replacementService := fmt.Sprintf(templateService, typed.Placeholder2,
			strings.Title(opts.TypeName),
			opts.TypeName,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
		)
		content = strings.Replace(content, typed.Placeholder2, replacementService, 1)

		// Add the service messages
		templateMessage := `%[1]v
message QueryGet%[2]vRequest {
	string index = 1;
}

message QueryGet%[2]vResponse {
	%[2]v %[2]v = 1;
}

message QueryAll%[2]vRequest {
	cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAll%[2]vResponse {
	repeated %[2]v %[2]v = 1;
	cosmos.base.query.v1beta1.PageResponse pagination = 2;
}`
		replacementMessage := fmt.Sprintf(templateMessage, typed.Placeholder3,
			strings.Title(opts.TypeName),
			opts.TypeName,
		)
		content = strings.Replace(content, typed.Placeholder3, replacementMessage, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func moduleGRPCGatewayModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replacement := `"context"`
		content := strings.Replace(f.String(), typed.Placeholder, replacement, 1)
		replacement = `types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))`
		content = strings.Replace(content, typed.Placeholder2, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientCliQueryModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v

	cmd.AddCommand(CmdList%[2]v())
	cmd.AddCommand(CmdShow%[2]v())
`
		replacement := fmt.Sprintf(template, typed.Placeholder,
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), typed.Placeholder, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisProtoModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/genesis.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateProtoImport := `%[1]v
import "%[2]v/%[3]v.proto";`
		replacementProtoImport := fmt.Sprintf(templateProtoImport, typed.PlaceholderGenesisProtoImport, opts.ModuleName, opts.TypeName)
		content := strings.Replace(f.String(), typed.PlaceholderGenesisProtoImport, replacementProtoImport, 1)

		// Determine the new field number
		fieldNumber := strings.Count(content, typed.PlaceholderGenesisProtoStateField) + 1

		templateProtoState := `%[1]v
		repeated %[2]v %[3]vList = %[4]v; %[5]v`
		replacementProtoState := fmt.Sprintf(
			templateProtoState,
			typed.PlaceholderGenesisProtoState,
			strings.Title(opts.TypeName),
			opts.TypeName,
			fieldNumber,
			typed.PlaceholderGenesisProtoStateField,
		)
		content = strings.Replace(content, typed.PlaceholderGenesisProtoState, replacementProtoState, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateTypesImport := `"fmt"`
		content := strings.Replace(f.String(), typed.PlaceholderGenesisTypesImport, templateTypesImport, 1)

		templateTypesDefault := `%[1]v
%[2]vList: []*%[2]v{},`
		replacementTypesDefault := fmt.Sprintf(templateTypesDefault, typed.PlaceholderGenesisTypesDefault, strings.Title(opts.TypeName))
		content = strings.Replace(content, typed.PlaceholderGenesisTypesDefault, replacementTypesDefault, 1)

		templateTypesValidate := `%[1]v
// Check for duplicated index in %[2]v
%[2]vIndexMap := make(map[string]bool)

for _, elem := range gs.%[3]vList {
	if _, ok := %[2]vIndexMap[elem.Index]; ok {
		return fmt.Errorf("duplicated index for %[2]v")
	}
	%[2]vIndexMap[elem.Index] = true
}`
		replacementTypesValidate := fmt.Sprintf(
			templateTypesValidate,
			typed.PlaceholderGenesisTypesValidate,
			opts.TypeName,
			strings.Title(opts.TypeName),
		)
		content = strings.Replace(content, typed.PlaceholderGenesisTypesValidate, replacementTypesValidate, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisModuleModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `%[1]v
// Set all the %[2]v
for _, elem := range genState.%[3]vList {
	k.Set%[3]v(ctx, *elem)
}

`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			typed.PlaceholderGenesisModuleInit,
			opts.TypeName,
			strings.Title(opts.TypeName),
		)
		content := strings.Replace(f.String(), typed.PlaceholderGenesisModuleInit, replacementModuleInit, 1)

		templateModuleExport := `%[1]v
// Get all %[2]v
%[2]vList := k.GetAll%[3]v(ctx)
for _, elem := range %[2]vList {
	elem := elem
	genesis.%[3]vList = append(genesis.%[3]vList, &elem)
}
`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			typed.PlaceholderGenesisModuleExport,
			opts.TypeName,
			strings.Title(opts.TypeName),
		)
		content = strings.Replace(content, typed.PlaceholderGenesisModuleExport, replacementModuleExport, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
