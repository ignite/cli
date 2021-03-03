package cosmosgen

import (
	"embed"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var (
	//go:embed templates/*
	templates embed.FS

	templateJSClient   = tpl("js/client.ts.tpl")   // js wrapper client.
	templateVuexStore  = tpl("vuex/store.ts.tpl")  // vuex store.
	templateVuexLoader = tpl("vuex/loader.ts.tpl") // vuex store loader.
)

// tpl returns a func for template residing at templatePath to initialize a text template
// with given protoPath.
func tpl(templatePath string) func(protoPath string) *template.Template {
	return func(protoPath string) *template.Template {
		path := filepath.Join("templates", templatePath)

		funcs := template.FuncMap{
			"camelCase": strcase.ToLowerCamel,
			"resolveFile": func(fullPath string) string {
				rel, _ := filepath.Rel(protoPath, fullPath)
				rel = strings.TrimSuffix(rel, ".proto")
				return rel
			},
		}

		return template.
			Must(
				template.
					New(filepath.Base(path)).
					Funcs(funcs).
					ParseFS(templates, path),
			)
	}
}
