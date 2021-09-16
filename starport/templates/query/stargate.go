package query

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
)

// NewStargate returns the generator to scaffold a empty query in a Stargate module
func NewStargate(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fsStargate,
			"stargate/",
			opts.AppPath,
		)
	)

	g.RunFn(protoQueryModify(replacer, opts))
	g.RunFn(cliQueryModify(replacer, opts))

	return g, Box(template, opts, g)
}

func protoQueryModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "query.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// RPC service
		templateRPC := `%[1]v

	// Queries a list of %[3]v items.
	rpc %[2]v(Query%[2]vRequest) returns (Query%[2]vResponse) {
		option (google.api.http).get = "/%[4]v/%[5]v/%[6]v/%[3]v";
	}
`
		replacementRPC := fmt.Sprintf(
			templateRPC,
			Placeholder2,
			opts.QueryName.UpperCamel,
			opts.QueryName.LowerCamel,
			opts.OwnerName,
			opts.AppName,
			opts.ModuleName,
		)
		content := replacer.Replace(f.String(), Placeholder2, replacementRPC)

		// Fields for request
		var reqFields string
		for i, field := range opts.ReqFields {
			reqFields += fmt.Sprintf("  %s\n", field.ProtoDeclaration(i+1))
		}
		if opts.Paginated {
			reqFields += fmt.Sprintf("cosmos.base.query.v1beta1.PageRequest pagination = %d;\n", len(opts.ReqFields)+1)
		}

		// Fields for response
		var resFields string
		for i, field := range opts.ResFields {
			resFields += fmt.Sprintf("  %s\n", field.ProtoDeclaration(i+1))
		}
		if opts.Paginated {
			resFields += fmt.Sprintf("cosmos.base.query.v1beta1.PageResponse pagination = %d;\n", len(opts.ResFields)+1)
		}

		// Ensure custom types are imported
		protoImports := append(opts.ResFields.ProtoImports(), opts.ReqFields.ProtoImports()...)
		customFields := append(opts.ResFields.Custom(), opts.ReqFields.Custom()...)
		for _, f := range customFields {
			protoImports = append(protoImports,
				fmt.Sprintf("%[1]v/%[2]v.proto", opts.ModuleName, f),
			)
		}
		for _, f := range protoImports {
			importModule := fmt.Sprintf(`
import "%[1]v";`, f)
			content = strings.ReplaceAll(content, importModule, "")

			replacementImport := fmt.Sprintf("%[1]v%[2]v", Placeholder, importModule)
			content = replacer.Replace(content, Placeholder, replacementImport)
		}

		// Messages
		templateMessages := `%[1]v
message Query%[2]vRequest {
%[3]v}

message Query%[2]vResponse {
%[4]v}
`
		replacementMessages := fmt.Sprintf(
			templateMessages,
			Placeholder3,
			opts.QueryName.UpperCamel,
			reqFields,
			resFields,
		)
		content = replacer.Replace(content, Placeholder3, replacementMessages)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func cliQueryModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "client/cli/query.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%[1]v

	cmd.AddCommand(Cmd%[2]v())
`
		replacement := fmt.Sprintf(
			template,
			Placeholder,
			opts.QueryName.UpperCamel,
		)
		content := replacer.Replace(f.String(), Placeholder, replacement)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
