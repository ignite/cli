package modulecreate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// NewIBC returns the generator to scaffold the implementation of the IBCModule interface inside a module.
func NewIBC(opts *CreateOptions) (*genny.Generator, error) {
	subFs, err := fs.Sub(fsIBC, "files/ibc")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	g.RunFn(genesisModify(opts))
	g.RunFn(genesisTypesModify(opts))
	g.RunFn(genesisProtoModify(opts))

	if err := g.OnlyFS(subFs, nil, nil); err != nil {
		return g, errors.Errorf("generator fs: %w", err)
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("protoVer", opts.ProtoVer)
	ctx.Set("ibcOrdering", opts.IBCOrdering)
	ctx.Set("dependencies", opts.Dependencies)
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName, opts.ProtoVer))

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))

	return g, nil
}

func genesisModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "keeper/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Genesis init
		replacementModuleInit := `if err := k.Port.Set(ctx, genState.PortId); err != nil {
		return err
	}`
		content, err := xast.ModifyFunction(
			f.String(),
			"InitGenesis",
			xast.AppendFuncCode(replacementModuleInit),
		)
		if err != nil {
			return err
		}

		// Genesis export
		replacementModuleExport := `genesis.PortId, err = k.Port.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}`
		content, err = xast.ModifyFunction(
			content,
			"ExportGenesis",
			xast.AppendFuncCode(replacementModuleExport),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("host", "github.com/cosmos/ibc-go/v10/modules/core/24-host"),
		)
		if err != nil {
			return err
		}

		// Default genesis
		content, err = xast.ModifyFunction(
			content,
			"DefaultGenesis",
			xast.AppendFuncStruct("GenesisState", "PortId", "PortID"),
		)
		if err != nil {
			return err
		}

		// Validate genesis
		// PlaceholderIBCGenesisTypeValidate
		replacementTypesValidate := `if err := host.PortIdentifierValidator(gs.PortId); err != nil {
	return err
}`
		content, err = xast.ModifyFunction(
			content,
			"Validate",
			xast.AppendFuncCode(replacementTypesValidate),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Modifies genesis.proto to add a new field.
//
// What it depends on:
//   - Existence of a message named 'GenesisState' in genesis.proto.
func genesisProtoModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("genesis.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		// Grab GenesisState and add next (always 2, I gather) available field.
		// TODO: typed.ProtoGenesisStateMessage exists but in subfolder, so we can't use it here, refactor?
		genesisState, err := protoutil.GetMessageByName(protoFile, "GenesisState")
		if err != nil {
			return errors.Errorf("couldn't find message 'GenesisState' in %s: %w", path, err)
		}
		field := protoutil.NewField("port_id", "string", protoutil.NextUniqueID(genesisState))
		protoutil.Append(genesisState, field)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func ensureGlobalValue(content string, globalType xast.GlobalType, name, value string) (string, error) {
	exists, err := globalExists(content, name)
	if err != nil {
		return "", err
	}
	if exists {
		return content, nil
	}
	return xast.InsertGlobal(content, globalType, xast.WithGlobal(name, "", value))
}

func globalExists(content, globalName string) (bool, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return false, err
	}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range valueSpec.Names {
				if name.Name == globalName {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func appIBCModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathIBCConfigGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport(
				fmt.Sprintf("%[1]vmodule", opts.ModuleName),
				fmt.Sprintf("%[1]v/x/%[2]v/module", opts.ModulePath, opts.ModuleName),
			),
			xast.WithNamedImport(
				fmt.Sprintf("%[1]vmoduletypes", opts.ModuleName),
				fmt.Sprintf("%[1]v/x/%[2]v/types", opts.ModulePath, opts.ModuleName),
			),
		)
		if err != nil {
			return err
		}

		content, err = addIBCModuleRoute(content, opts.ModuleName)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func addIBCModuleRoute(content, moduleName string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return "", err
	}
	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

	registerIBCModules := findFunctionByName(file, "registerIBCModules")
	if registerIBCModules == nil {
		return "", errors.New(`function "registerIBCModules" not found`)
	}

	moduleNameExprText := fmt.Sprintf("%smoduletypes.ModuleName", moduleName)
	moduleConstructorExprText := fmt.Sprintf(
		"%smodule.NewIBCModule(app.appCodec, app.%sKeeper)",
		moduleName,
		xstrings.Title(moduleName),
	)

	hasRoute, err := hasAddRoute(registerIBCModules, moduleNameExprText, fileSet)
	if err != nil {
		return "", err
	}
	if !hasRoute {
		insertionIndex := -1
		for i, stmt := range registerIBCModules.Body.List {
			assignStmt, ok := stmt.(*ast.AssignStmt)
			if !ok || len(assignStmt.Lhs) != 1 {
				continue
			}

			lhs, ok := assignStmt.Lhs[0].(*ast.Ident)
			if !ok || lhs.Name != "ibcv2Router" {
				continue
			}

			insertionIndex = i
			break
		}
		if insertionIndex == -1 {
			return "", errors.New(`assignment to "ibcv2Router" not found`)
		}

		stmtCode := fmt.Sprintf(
			"ibcRouter = ibcRouter.AddRoute(%s, %s)",
			moduleNameExprText,
			moduleConstructorExprText,
		)
		statements, err := parseStatements(stmtCode)
		if err != nil {
			return "", err
		}
		if len(statements) != 1 {
			return "", errors.New("unexpected number of statements while creating ibc route assignment")
		}

		registerIBCModules.Body.List = append(
			registerIBCModules.Body.List[:insertionIndex],
			append([]ast.Stmt{statements[0]}, registerIBCModules.Body.List[insertionIndex:]...)...,
		)
	}

	file.Comments = commentMap.Filter(file).Comments()

	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, file); err != nil {
		return "", err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func hasAddRoute(funcDecl *ast.FuncDecl, moduleNameExpr string, fileSet *token.FileSet) (bool, error) {
	var (
		found bool
		err   error
	)

	ast.Inspect(funcDecl, func(n ast.Node) bool {
		if found || err != nil {
			return false
		}

		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "AddRoute" || len(callExpr.Args) == 0 {
			return true
		}

		argText, argErr := exprString(fileSet, callExpr.Args[0])
		if argErr != nil {
			err = argErr
			return false
		}

		if normalizedExpr(argText) == normalizedExpr(moduleNameExpr) {
			found = true
			return false
		}

		return true
	})

	return found, err
}

func parseStatements(code string) ([]ast.Stmt, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", fmt.Sprintf("package p\nfunc _(){\n%s\n}", code), 0)
	if err != nil {
		return nil, err
	}

	funcDecl, ok := file.Decls[0].(*ast.FuncDecl)
	if !ok || funcDecl.Body == nil {
		return nil, errors.New("failed to parse statements")
	}

	return funcDecl.Body.List, nil
}

func findFunctionByName(file *ast.File, funcName string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name == funcName {
			return funcDecl
		}
	}
	return nil
}
