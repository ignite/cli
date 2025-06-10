package cosmosgen

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed templates/*
	templates embed.FS

	templateTSClientRoot           = newTemplateWriter("root")
	templateTSClientModule         = newTemplateWriter("module")
	templateTSClientVue            = newTemplateWriter("vue")
	templateTSClientVueRoot        = newTemplateWriter("vue-root")
	templateTSClientRest           = newTemplateWriter("rest")
	templateTSClientComposable     = newTemplateWriter("composable")
	templateTSClientComposableRoot = newTemplateWriter("composable-root")
)

type templateWriter struct {
	templateDir string
}

// tpl returns a func for template residing at templatePath to initialize a text template
// with given protoPath.
func newTemplateWriter(templateDir string) templateWriter {
	return templateWriter{
		templateDir,
	}
}

func (t templateWriter) Write(destDir, protoPath string, data interface{}) error {
	base := filepath.Join("templates", t.templateDir)

	// find out templates inside the dir.
	files, err := templates.ReadDir(base)
	if err != nil {
		return err
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(base, file.Name()))
	}

	funcs := template.FuncMap{
		"camelCase": strcase.ToLowerCamel,
		"capitalCase": func(word string) string {
			replacer := strings.NewReplacer("-", "_", ".", "_")
			word = strcase.ToCamel(replacer.Replace(word))

			return cases.Title(language.English).String(word)
		},
		"camelCaseLowerSta": func(word string) string {
			replacer := strings.NewReplacer("-", "_", ".", "_")

			return strcase.ToLowerCamel(replacer.Replace(word))
		},
		"camelCaseUpperSta": func(word string) string {
			replacer := strings.NewReplacer("-", "_", ".", "_")

			return strcase.ToCamel(replacer.Replace(word))
		},
		"resolveFile": func(fullPath string) string {
			_ = protoPath // eventually, we should use the proto folder name of this, for the application (but not for the other modules)

			res := strings.Split(fullPath, "proto/")
			rel := res[len(res)-1] // get path after proto/
			rel = strings.TrimSuffix(rel, ".proto")

			return "./types/" + rel
		},
		"transformPath": func(path string) string {
			// transformPath converts a endpoint path to a valid JS substring path.
			// e.g. /cosmos/bank/v1beta1/spendable_balances/{address}/by_denom -> /cosmos/bank/v1beta1/spendable_balances/${address}/by_denom
			path = strings.ReplaceAll(path, "{", "${")
			path = strings.ReplaceAll(path, "=**}", "}")
			return path
		},
		"mapToTypeScriptObject": func(m map[string]string) string {
			// mapToTypeScriptObject converts a map to a TypeScript object string.
			// e.g. {"key1"?: value1; "key2"?: value2}

			sortedKeys := make([]string, 0, len(m))
			for k := range m {
				sortedKeys = append(sortedKeys, k)
			}
			sort.Strings(sortedKeys)

			var sb strings.Builder
			sb.WriteString("{")
			sb.WriteString("\n")
			for _, k := range sortedKeys {
				typeStr := m[k]

				if strings.Contains(typeStr, ".") {
					// TODO(@julienrbrt): parse proto types to deepest inner type and remove hardcoded pagination types.
					if strings.Contains(typeStr, ".") {
						if strings.EqualFold(typeStr, "cosmos.base.query.v1beta1.PageRequest") {
							sb.WriteString(`      "pagination.key"?: string;`)
							sb.WriteString("\n")
							sb.WriteString(`      "pagination.offset"?: string;`)
							sb.WriteString("\n")
							sb.WriteString(`      "pagination.limit"?: string;`)
							sb.WriteString("\n")
							sb.WriteString(`      "pagination.count_total"?: boolean;`)
							sb.WriteString("\n")
							sb.WriteString(`      "pagination.reverse"?: boolean;`)
							sb.WriteString("\n")
							continue
						}

						sb.WriteString(fmt.Sprintf(`      "%s"?: any /* TODO */;`, k))
						sb.WriteString("\n")
						continue
					}
				}

				sb.WriteString(fmt.Sprintf(`      "%s"?: %s;`, k, protoTypeToTypeScriptType(typeStr)))
				sb.WriteString("\n")
			}
			sb.WriteString("    }")
			return sb.String()
		},
		"inc": func(i int) int {
			return i + 1
		},
		"replace": strings.ReplaceAll,
	}

	// render and write the template.
	write := func(path string) error {
		tpl := template.
			Must(
				template.
					New(filepath.Base(path)).
					Funcs(funcs).
					ParseFS(templates, paths...),
			)

		out := filepath.Join(destDir, strings.TrimSuffix(filepath.Base(path), ".tpl"))

		f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o766)
		if err != nil {
			return err
		}
		defer f.Close()

		return tpl.Execute(f, data)
	}

	for _, path := range paths {
		if err := write(path); err != nil {
			return err
		}
	}

	return nil
}

// protoTypeToTypeScriptType converts a proto type string to a TypeScript type string.
// e.g. "string" -> "string", "int32" -> "number", "bool" -> "boolean", "bytes" -> "Uint8Array", etc.
func protoTypeToTypeScriptType(pt string) string {
	isRepeated := strings.HasPrefix(pt, "repeated ")
	if isRepeated {
		pt = strings.TrimPrefix(pt, "repeated ")
	}

	var tsBaseType string
	switch pt {
	case "string":
		tsBaseType = "string"
	case "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "float", "double":
		tsBaseType = "number"
	case "bool":
		tsBaseType = "boolean"
	case "bytes":
		tsBaseType = "Uint8Array"
	default:
		tsBaseType = pt
	}

	if isRepeated {
		return tsBaseType + "[]"
	}
	return tsBaseType
}
