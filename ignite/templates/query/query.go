package query

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

//go:embed files/* files/**/*
var files embed.FS

func Box(box fs.FS, opts *Options, g *genny.Generator) error {
	if err := g.OnlyFS(box, nil, nil); err != nil {
		return err
	}
	ctx := plush.NewContext()
	ctx.Set("ModuleName", opts.ModuleName)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("ProtoVer", opts.ProtoVer)
	ctx.Set("QueryName", opts.QueryName)
	ctx.Set("Description", opts.Description)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("ReqFields", opts.ReqFields)
	ctx.Set("ResFields", opts.ResFields)
	ctx.Set("Paginated", opts.Paginated)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))
	g.Transformer(genny.Replace("{{queryName}}", opts.QueryName.Snake))
	return nil
}

// NewGenerator returns the generator to scaffold a empty query in a module.
func NewGenerator(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	subFs, err := fs.Sub(files, "files")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	g.RunFn(protoQueryModify(opts))
	g.RunFn(cliQueryModify(replacer, opts))

	return g, Box(subFs, opts, g)
}

// Modifies query.proto to add the required RPCs and Messages.
//
// What it depends on:
//   - Existence of a service with name "Query" since that is where the RPCs will be added.
func protoQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("query.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		// if the query has request fields, they are appended to the rpc query
		var requestPath string
		for _, field := range opts.ReqFields {
			requestPath += "/"
			requestPath = filepath.Join(requestPath, fmt.Sprintf("{%s}", field.ProtoFieldName()))
		}
		serviceQuery, err := protoutil.GetServiceByName(protoFile, "Query")
		if err != nil {
			return errors.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}

		typenamePascal, appModulePath := opts.QueryName.PascalCase, gomodulepath.ExtractAppPath(opts.ModulePath)
		rpcSingle := protoutil.NewRPC(
			typenamePascal,
			fmt.Sprintf("Query%sRequest", typenamePascal),
			fmt.Sprintf("Query%sResponse", typenamePascal),
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s/%s%s",
						appModulePath, opts.ModuleName, opts.ProtoVer, opts.QueryName.Snake, requestPath,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.AttachComment(rpcSingle, fmt.Sprintf("%[1]v Queries a list of %[1]v items.", typenamePascal))
		protoutil.Append(serviceQuery, rpcSingle)

		// Fields for request
		paginationType, paginationName := "cosmos.base.query.v1beta1.Page", "pagination"
		var reqFields []*proto.NormalField
		for i, field := range opts.ReqFields {
			reqFields = append(reqFields, field.ToProtoField(i+1))
		}
		if opts.Paginated {
			reqFields = append(reqFields, protoutil.NewField(paginationName, paginationType+"Request", len(opts.ReqFields)+1))
		}
		requestMessage := protoutil.NewMessage("Query"+typenamePascal+"Request", protoutil.WithFields(reqFields...))

		// Fields for response
		var resFields []*proto.NormalField
		for i, field := range opts.ResFields {
			resFields = append(resFields, field.ToProtoField(i+1))
		}
		if opts.Paginated {
			resFields = append(resFields, protoutil.NewField(paginationName, paginationType+"Response", len(opts.ResFields)+1))
		}
		responseMessage := protoutil.NewMessage("Query"+typenamePascal+"Response", protoutil.WithFields(resFields...))
		protoutil.Append(protoFile, requestMessage, responseMessage)

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range append(opts.ResFields.ProtoImports(), opts.ReqFields.ProtoImports()...) {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range append(opts.ResFields.Custom(), opts.ReqFields.Custom()...) {
			protoPath := fmt.Sprintf("%[1]v/%[2]v/%[3]v/%[4]v.proto", opts.AppName, opts.ModuleName, opts.ProtoVer, f)
			protoImports = append(protoImports, protoutil.NewImport(protoPath))
		}
		if err = protoutil.AddImports(protoFile, true, protoImports...); err != nil {
			return errors.Errorf("failed to add imports to %s: %w", path, err)
		}

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func cliQueryModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "module/autocli.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `{
					RpcMethod: "%[2]v",
					Use: "%[3]v",
					Short: "%[4]v",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{%[5]s},
				},

				%[1]v`
		replacement := fmt.Sprintf(
			template,
			PlaceholderAutoCLIQuery,
			opts.QueryName.PascalCase,
			fmt.Sprintf("%s %s", opts.QueryName.Kebab, opts.ReqFields.CLIUsage()),
			opts.Description,
			opts.ReqFields.ProtoFieldNameAutoCLI(),
		)
		content := replacer.Replace(f.String(), PlaceholderAutoCLIQuery, replacement)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
