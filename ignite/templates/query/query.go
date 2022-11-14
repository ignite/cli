package query

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/emicklei/proto"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
)

//go:embed files/* files/**/*
var fs embed.FS

func Box(box packd.Walker, opts *Options, g *genny.Generator) error {
	if err := g.Box(box); err != nil {
		return err
	}
	ctx := plush.NewContext()
	ctx.Set("ModuleName", opts.ModuleName)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("QueryName", opts.QueryName)
	ctx.Set("Description", opts.Description)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("ReqFields", opts.ReqFields)
	ctx.Set("ResFields", opts.ResFields)
	ctx.Set("Paginated", opts.Paginated)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{queryName}}", opts.QueryName.Snake))
	return nil
}

// NewGenerator returns the generator to scaffold a empty query in a module
func NewGenerator(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(
			fs,
			"files/",
			opts.AppPath,
		)
	)

	g.RunFn(protoQueryModify(opts))
	g.RunFn(cliQueryModify(replacer, opts))

	return g, Box(template, opts, g)
}

// Modifies query.proto to add the required RPCs and Messages.
//
// What it depends on:
//   - Existence of a service with name "Query" since that is where the RPCs will be added.
func protoQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.AppName, opts.ModuleName, "query.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		pf, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		// if the query has request fields, they are appended to the rpc query
		var reqPath string
		for _, field := range opts.ReqFields {
			reqPath += "/"
			reqPath = filepath.Join(reqPath, fmt.Sprintf("{%s}", field.ProtoFieldName()))
		}
		srv, err := protoutil.GetServiceByName(pf, "Query")
		if err != nil {
			return fmt.Errorf("failed while looking up service 'Query' in %s: %w", path, err)
		}

		typU, appModulePath := opts.QueryName.UpperCamel, gomodulepath.ExtractAppPath(opts.ModulePath)
		single := protoutil.NewRPC(typU, "Query"+typU+"Request", "Query"+typU+"Response",
			protoutil.WithRPCOptions(
				protoutil.NewOption(
					"google.api.http",
					fmt.Sprintf(
						"/%s/%s/%s%s",
						appModulePath, opts.ModuleName, opts.QueryName.Snake, reqPath,
					),
					protoutil.Custom(),
					protoutil.SetField("get"),
				),
			),
		)
		protoutil.Append(srv, single)

		// Fields for request
		pagT, pagN := "cosmos.base.query.v1beta1.Page", "pagination"
		var reqFields []*proto.NormalField
		for i, field := range opts.ReqFields {
			reqFields = append(reqFields, field.ToProtoField(i+1))
		}
		if opts.Paginated {
			reqFields = append(reqFields, protoutil.NewField(pagT+"Request", pagN, len(opts.ReqFields)+1))
		}
		msgReq := protoutil.NewMessage("Query"+typU+"Request", protoutil.WithFields(reqFields...))

		// Fields for response
		var resFields []*proto.NormalField
		for i, field := range opts.ResFields {
			resFields = append(resFields, field.ToProtoField(i+1))
		}
		if opts.Paginated {
			resFields = append(resFields, protoutil.NewField(pagT+"Response", pagN, len(opts.ResFields)+1))
		}
		msgResp := protoutil.NewMessage("Query"+typU+"Response", protoutil.WithFields(resFields...))
		protoutil.Append(pf, msgReq, msgResp)

		// Ensure custom types are imported
		var protoImports []*proto.Import
		for _, imp := range append(opts.ResFields.ProtoImports(), opts.ReqFields.ProtoImports()...) {
			protoImports = append(protoImports, protoutil.NewImport(imp))
		}
		for _, f := range append(opts.ResFields.Custom(), opts.ReqFields.Custom()...) {
			protopath := fmt.Sprintf("%[1]v/%[2]v/%[3]v.proto", opts.AppName, opts.ModuleName, f)
			protoImports = append(protoImports, protoutil.NewImport(protopath))
		}
		if err = protoutil.AddImports(pf, true, protoImports...); err != nil {
			// shouldn't really occur.
			return fmt.Errorf("failed to add imports to %s: %w", path, err)
		}

		newFile := genny.NewFileS(path, protoutil.Printer(pf))
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

		template := `cmd.AddCommand(Cmd%[2]v())

%[1]v`
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
