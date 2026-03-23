package modulecreate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

type addModuleAppConfigOptions struct {
	skipConfig    bool
	runtimeFields []string
}

type AddModuleAppConfigOption func(*addModuleAppConfigOptions)

func SkipConfigEntry() AddModuleAppConfigOption {
	return func(opts *addModuleAppConfigOptions) {
		opts.skipConfig = true
	}
}

// SpecifyModuleEntry allows to define to which field the module should be added in the app config.
// E.g. "PreBlockers", "InitGenesis", "BeginBlockers", "EndBlockers"
func SpecifyModuleEntry(fields ...string) AddModuleAppConfigOption {
	return func(opts *addModuleAppConfigOptions) {
		opts.runtimeFields = fields
	}
}

// AddModuleToAppConfig appends a given module to the chain app config.
func AddModuleToAppConfig(content, moduleName string, opts ...AddModuleAppConfigOption) (string, error) {
	options := addModuleAppConfigOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return AddModuleToAppConfigWithOptions(content, moduleName, options)
}

// AddModuleToAppConfigWithOptions appends a given module to the chain app config with options.
func AddModuleToAppConfigWithOptions(content, moduleName string, opts addModuleAppConfigOptions) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", content, parser.ParseComments)
	if err != nil {
		return "", err
	}
	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

	appConfigLit, err := findAppConfigCompositeLiteral(file)
	if err != nil {
		return "", err
	}

	modulesField, err := findKeyValueByName(appConfigLit, "Modules")
	if err != nil {
		return "", err
	}

	runtimeModuleLit, err := findRuntimeModuleCompositeLiteral(file, modulesField.Value, fileSet)
	if err != nil {
		return "", err
	}

	fields := opts.runtimeFields
	if len(fields) == 0 {
		fields = []string{"InitGenesis", "BeginBlockers", "EndBlockers"}
	}

	for _, fieldName := range fields {
		if err := appendModuleNameToRuntimeField(file, runtimeModuleLit, fieldName, moduleName, fileSet); err != nil {
			return "", err
		}
	}

	if !opts.skipConfig {
		if err := appendModuleConfigEntry(file, modulesField.Value, moduleName, fileSet); err != nil {
			return "", err
		}
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

func findAppConfigCompositeLiteral(file *ast.File) (*ast.CompositeLit, error) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if name.Name != "appConfig" && name.Name != "AppConfig" {
					continue
				}
				if len(valueSpec.Values) == 0 {
					return nil, errors.Errorf("%s has no value", name.Name)
				}

				valueIdx := i
				if valueIdx >= len(valueSpec.Values) {
					valueIdx = 0
				}

				return findCompositeLiteralByType(file, valueSpec.Values[valueIdx], "appv1alpha1", "Config")
			}
		}
	}

	return nil, errors.New("app config variable not found")
}

func findRuntimeModuleCompositeLiteral(
	file *ast.File,
	modulesExpr ast.Expr,
	fileSet *token.FileSet,
) (*ast.CompositeLit, error) {
	modulesLit, err := resolveCompositeLiteral(file, modulesExpr)
	if err != nil {
		return nil, errors.Errorf("resolve modules list: %w", err)
	}

	for _, elt := range modulesLit.Elts {
		moduleConfigLit, err := resolveCompositeLiteral(file, elt)
		if err != nil {
			continue
		}

		nameField, err := findKeyValueByName(moduleConfigLit, "Name")
		if err != nil {
			continue
		}

		nameValue, err := exprString(fileSet, nameField.Value)
		if err != nil {
			return nil, err
		}
		if nameValue != "runtime.ModuleName" {
			continue
		}

		configField, err := findKeyValueByName(moduleConfigLit, "Config")
		if err != nil {
			return nil, errors.Errorf("runtime module config field not found: %w", err)
		}

		return findCompositeLiteralByType(file, configField.Value, "runtimev1alpha1", "Module")
	}

	return nil, errors.New("runtime module not found in app config")
}

func appendModuleNameToRuntimeField(
	file *ast.File,
	runtimeModuleLit *ast.CompositeLit,
	fieldName, moduleName string,
	fileSet *token.FileSet,
) error {
	field, err := findKeyValueByName(runtimeModuleLit, fieldName)
	if err != nil {
		return errors.Errorf("%s field not found in runtime module: %w", fieldName, err)
	}

	listLit, err := resolveCompositeLiteral(file, field.Value)
	if err != nil {
		return errors.Errorf("resolve %s list: %w", fieldName, err)
	}

	moduleExprText := fmt.Sprintf("%smoduletypes.ModuleName", moduleName)
	normalizedModuleExpr := normalizedExpr(moduleExprText)

	for _, elt := range listLit.Elts {
		existing, err := exprString(fileSet, elt)
		if err != nil {
			return err
		}
		if normalizedExpr(existing) == normalizedModuleExpr {
			return nil
		}
	}

	appendCompositeLiteralElement(fileSet, listLit, moduleExprText)
	return nil
}

func appendModuleConfigEntry(
	file *ast.File,
	modulesExpr ast.Expr,
	moduleName string,
	fileSet *token.FileSet,
) error {
	modulesLit, err := resolveCompositeLiteral(file, modulesExpr)
	if err != nil {
		return errors.Errorf("resolve modules list: %w", err)
	}

	moduleNameText := fmt.Sprintf("%smoduletypes.ModuleName", moduleName)
	moduleNamePattern := normalizedExpr(fmt.Sprintf("Name:%s", moduleNameText))

	for _, elt := range modulesLit.Elts {
		existingExpr, err := exprString(fileSet, elt)
		if err == nil && strings.Contains(normalizedExpr(existingExpr), moduleNamePattern) {
			return nil
		}

		moduleConfigLit, err := resolveCompositeLiteral(file, elt)
		if err != nil {
			continue
		}

		nameField, err := findKeyValueByName(moduleConfigLit, "Name")
		if err != nil {
			continue
		}

		existingName, err := exprString(fileSet, nameField.Value)
		if err != nil {
			return err
		}
		if existingName == moduleNameText {
			return nil
		}
	}

	newEntry := ast.NewIdent(fmt.Sprintf(
		`{
	Name:   %smoduletypes.ModuleName,
	Config: appconfig.WrapAny(&%smoduletypes.Module{}),
}`,
		moduleName,
		moduleName,
	))

	appendCompositeLiteralElement(fileSet, modulesLit, newEntry.Name)
	return nil
}

func findCompositeLiteralByType(
	file *ast.File,
	expr ast.Expr,
	pkgName, typeName string,
) (*ast.CompositeLit, error) {
	lit := findCompositeLiteralByTypeExpr(file, expr, pkgName, typeName, map[string]struct{}{})
	if lit == nil {
		return nil, errors.Errorf("composite literal %s.%s not found", pkgName, typeName)
	}
	return lit, nil
}

func findCompositeLiteralByTypeExpr(
	file *ast.File,
	expr ast.Expr,
	pkgName, typeName string,
	visited map[string]struct{},
) *ast.CompositeLit {
	switch typedExpr := expr.(type) {
	case *ast.ParenExpr:
		return findCompositeLiteralByTypeExpr(file, typedExpr.X, pkgName, typeName, visited)
	case *ast.UnaryExpr:
		return findCompositeLiteralByTypeExpr(file, typedExpr.X, pkgName, typeName, visited)
	case *ast.CallExpr:
		for _, arg := range typedExpr.Args {
			if lit := findCompositeLiteralByTypeExpr(file, arg, pkgName, typeName, visited); lit != nil {
				return lit
			}
		}
	case *ast.CompositeLit:
		if isSelectorType(typedExpr.Type, pkgName, typeName) {
			return typedExpr
		}
		for _, elt := range typedExpr.Elts {
			keyValue, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if lit := findCompositeLiteralByTypeExpr(file, keyValue.Value, pkgName, typeName, visited); lit != nil {
				return lit
			}
		}
	case *ast.Ident:
		if _, ok := visited[typedExpr.Name]; ok {
			return nil
		}
		visited[typedExpr.Name] = struct{}{}

		valueExpr, err := findGlobalValueExpr(file, typedExpr.Name)
		if err != nil {
			return nil
		}

		return findCompositeLiteralByTypeExpr(file, valueExpr, pkgName, typeName, visited)
	}

	return nil
}

func resolveCompositeLiteral(file *ast.File, expr ast.Expr) (*ast.CompositeLit, error) {
	switch typedExpr := expr.(type) {
	case *ast.CompositeLit:
		return typedExpr, nil
	case *ast.ParenExpr:
		return resolveCompositeLiteral(file, typedExpr.X)
	case *ast.UnaryExpr:
		return resolveCompositeLiteral(file, typedExpr.X)
	case *ast.Ident:
		valueExpr, err := findGlobalValueExpr(file, typedExpr.Name)
		if err != nil {
			return nil, err
		}
		return resolveCompositeLiteral(file, valueExpr)
	default:
		return nil, errors.Errorf("unsupported composite literal expression %T", expr)
	}
}

func findGlobalValueExpr(file *ast.File, name string) (ast.Expr, error) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, varName := range valueSpec.Names {
				if varName.Name != name {
					continue
				}

				if len(valueSpec.Values) == 0 {
					return nil, errors.Errorf("global variable %q has no value", name)
				}

				valueIdx := i
				if valueIdx >= len(valueSpec.Values) {
					valueIdx = 0
				}

				return valueSpec.Values[valueIdx], nil
			}
		}
	}

	return nil, errors.Errorf("global variable %q not found", name)
}

func findKeyValueByName(compLit *ast.CompositeLit, name string) (*ast.KeyValueExpr, error) {
	for _, elt := range compLit.Elts {
		keyValue, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := keyValue.Key.(*ast.Ident)
		if ok && key.Name == name {
			return keyValue, nil
		}
	}

	return nil, errors.Errorf("field %q not found", name)
}

func isSelectorType(expr ast.Expr, pkgName, typeName string) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkgIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	return pkgIdent.Name == pkgName && selector.Sel.Name == typeName
}

func exprString(fileSet *token.FileSet, expr ast.Expr) (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, expr); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func normalizedExpr(expr string) string {
	expr = strings.ReplaceAll(expr, " ", "")
	expr = strings.ReplaceAll(expr, "\n", "")
	expr = strings.ReplaceAll(expr, "\t", "")
	return expr
}

func appendCompositeLiteralElement(fileSet *token.FileSet, compLit *ast.CompositeLit, code string) {
	file := fileSet.File(compLit.Pos())
	maxOffset := file.Offset(compLit.Rbrace)
	for _, elt := range compLit.Elts {
		if pos := elt.End(); pos.IsValid() {
			offset := file.Offset(pos)
			if offset > maxOffset {
				maxOffset = offset
			}
		}
	}

	insertPos := file.Pos(maxOffset)
	value := ast.NewIdent(code)
	value.NamePos = insertPos

	compLit.Elts = append(compLit.Elts, value)
	compLit.Rbrace += token.Pos(1)

	if len(compLit.Elts) > 0 {
		last := compLit.Elts[len(compLit.Elts)-1]
		if file.Line(compLit.Rbrace) == file.Line(last.End())-1 {
			file.AddLine(file.Offset(compLit.Rbrace))
			compLit.Rbrace += token.Pos(1)
		}
	}
}
