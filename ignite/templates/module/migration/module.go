package modulemigration

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

func moduleModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		f, err := r.Disk.Find(opts.ModuleFile())
		if err != nil {
			return err
		}

		content, err := updateModule(f.String(), opts)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(opts.ModuleFile(), content))
	}
}

func updateModule(content string, opts *Options) (string, error) {
	currentVersion, err := ConsensusVersion(content)
	if err != nil {
		return "", err
	}
	if currentVersion != opts.FromVersion {
		return "", errors.Errorf("expected module consensus version %d, got %d", opts.FromVersion, currentVersion)
	}

	content, err = xast.AppendImports(
		content,
		xast.WithNamedImport(opts.MigrationImportAlias(), opts.MigrationImportPath()),
	)
	if err != nil {
		return "", err
	}

	content, err = setConsensusVersion(content, opts.ToVersion)
	if err != nil {
		return "", err
	}

	return addMigrationRegistration(content, opts)
}

// ConsensusVersion returns the current module consensus version from module.go content.
func ConsensusVersion(content string) (uint64, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return 0, err
	}

	expr, err := consensusVersionExpr(file)
	if err != nil {
		return 0, err
	}

	return parseConsensusVersionExpr(file, expr)
}

func addMigrationRegistration(content string, opts *Options) (string, error) {
	info, err := registerServicesInfoFromContent(content)
	if err != nil {
		return "", err
	}

	var functionOptions []xast.FunctionOptions

	if info.needsConfiguratorSetup {
		functionOptions = append(functionOptions, xast.AppendFuncCode(configuratorSetupCode(info)))
	}

	functionOptions = append(functionOptions, xast.AppendFuncCode(migrationRegistrationCode(info, opts)))

	return xast.ModifyFunction(content, "RegisterServices", functionOptions...)
}

func configuratorSetupCode(info registerServicesInfo) string {
	returnStmt := "return"
	if info.returnsError {
		returnStmt = "return nil"
	}

	return info.cfgVar + ", ok := " + info.parameterName + ".(module.Configurator)\n" +
		"if !ok {\n\t" + returnStmt + "\n}"
}

func migrationRegistrationCode(info registerServicesInfo, opts *Options) string {
	handleErr := "panic(err)"
	if info.returnsError {
		handleErr = "return err"
	}

	return "if err := " + info.cfgVar +
		".RegisterMigration(types.ModuleName, " +
		strconv.FormatUint(opts.FromVersion, 10) + ", " +
		opts.MigrationImportAlias() + "." + opts.MigrationFunc() +
		"); err != nil {\n\t" + handleErr + "\n}"
}

func setConsensusVersion(content string, version uint64) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return "", err
	}

	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

	expr, err := consensusVersionExpr(file)
	if err != nil {
		return "", err
	}

	switch versionExpr := expr.(type) {
	case *ast.BasicLit:
		versionExpr.Value = strconv.FormatUint(version, 10)
	case *ast.Ident:
		valueSpec, valueIndex, err := findValueSpec(file, versionExpr.Name)
		if err != nil {
			return "", err
		}
		valueSpec.Values[valueIndex] = &ast.BasicLit{
			Kind:  token.INT,
			Value: strconv.FormatUint(version, 10),
		}
	default:
		return "", errors.Errorf("unsupported consensus version expression %T", expr)
	}

	file.Comments = commentMap.Filter(file).Comments()

	return formatFile(fileSet, file)
}

type registerServicesInfo struct {
	cfgVar                 string
	needsConfiguratorSetup bool
	parameterName          string
	returnsError           bool
}

func registerServicesInfoFromContent(content string) (registerServicesInfo, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return registerServicesInfo{}, err
	}

	funcDecl := findFuncDecl(file, "RegisterServices")
	if funcDecl == nil {
		return registerServicesInfo{}, errors.New("function \"RegisterServices\" not found")
	}
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 || len(funcDecl.Type.Params.List[0].Names) == 0 {
		return registerServicesInfo{}, errors.New("RegisterServices must have a named parameter")
	}

	param := funcDecl.Type.Params.List[0]
	info := registerServicesInfo{
		parameterName: param.Names[0].Name,
		returnsError:  functionReturnsError(funcDecl),
	}

	if isModuleConfiguratorType(param.Type) {
		info.cfgVar = info.parameterName
		return info, nil
	}

	cfgVar := findConfiguratorVar(funcDecl, info.parameterName)
	if cfgVar != "" {
		info.cfgVar = cfgVar
		return info, nil
	}

	info.cfgVar = "cfg"
	info.needsConfiguratorSetup = true

	return info, nil
}

func functionReturnsError(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) != 1 {
		return false
	}

	ident, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident)
	return ok && ident.Name == "error"
}

func findConfiguratorVar(funcDecl *ast.FuncDecl, parameterName string) string {
	for _, stmt := range funcDecl.Body.List {
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok || len(assignStmt.Lhs) < 1 || len(assignStmt.Rhs) != 1 {
			continue
		}

		typeAssert, ok := assignStmt.Rhs[0].(*ast.TypeAssertExpr)
		if !ok || !isModuleConfiguratorType(typeAssert.Type) {
			continue
		}

		ident, ok := typeAssert.X.(*ast.Ident)
		if !ok || ident.Name != parameterName {
			continue
		}

		cfgVar, ok := assignStmt.Lhs[0].(*ast.Ident)
		if !ok {
			continue
		}

		return cfgVar.Name
	}

	return ""
}

func isModuleConfiguratorType(expr ast.Expr) bool {
	switch typedExpr := expr.(type) {
	case *ast.Ident:
		return typedExpr.Name == "Configurator"
	case *ast.SelectorExpr:
		return typedExpr.Sel.Name == "Configurator"
	default:
		return false
	}
}

func consensusVersionExpr(file *ast.File) (ast.Expr, error) {
	funcDecl := findFuncDecl(file, "ConsensusVersion")
	if funcDecl == nil {
		return nil, errors.New("function \"ConsensusVersion\" not found")
	}
	if funcDecl.Body == nil || len(funcDecl.Body.List) == 0 {
		return nil, errors.New("ConsensusVersion has an empty body")
	}

	lastStmt, ok := funcDecl.Body.List[len(funcDecl.Body.List)-1].(*ast.ReturnStmt)
	if !ok || len(lastStmt.Results) != 1 {
		return nil, errors.New("ConsensusVersion must return exactly one value")
	}

	return lastStmt.Results[0], nil
}

func parseConsensusVersionExpr(file *ast.File, expr ast.Expr) (uint64, error) {
	switch typedExpr := expr.(type) {
	case *ast.BasicLit:
		return parseConsensusVersionLiteral(typedExpr)
	case *ast.Ident:
		valueSpec, valueIndex, err := findValueSpec(file, typedExpr.Name)
		if err != nil {
			return 0, err
		}
		return parseConsensusVersionExpr(file, valueSpec.Values[valueIndex])
	default:
		return 0, errors.Errorf("unsupported consensus version expression %T", expr)
	}
}

func parseConsensusVersionLiteral(lit *ast.BasicLit) (uint64, error) {
	if lit.Kind != token.INT {
		return 0, errors.Errorf("unsupported consensus version literal kind %v", lit.Kind)
	}

	version, err := strconv.ParseUint(lit.Value, 10, 64)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func findValueSpec(file *ast.File, name string) (*ast.ValueSpec, int, error) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || (genDecl.Tok != token.CONST && genDecl.Tok != token.VAR) {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, specName := range valueSpec.Names {
				if specName.Name != name {
					continue
				}
				if len(valueSpec.Values) == 0 {
					return nil, 0, errors.Errorf("%s has no value", name)
				}

				valueIndex := i
				if valueIndex >= len(valueSpec.Values) {
					valueIndex = len(valueSpec.Values) - 1
				}

				return valueSpec, valueIndex, nil
			}
		}
	}

	return nil, 0, errors.Errorf("%s value not found", name)
}

func findFuncDecl(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == name {
			return funcDecl
		}
	}

	return nil
}

func formatFile(fileSet *token.FileSet, file *ast.File) (string, error) {
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
