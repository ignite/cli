package cosmosgen

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/takuoki/gocase"
)

var (
	//go:embed templates/*
	templates embed.FS

	templateJSClient  = newTemplateWriter("js")         // js wrapper client.
	templateVuexRoot  = newTemplateWriter("vuex/root")  // vuex store loader.
	templateVuexStore = newTemplateWriter("vuex/store") // vuex store.

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
		"camelCaseSta": func(word string) string {
			return gocase.Revert(strcase.ToLowerCamel(word))
		},
		"resolveFile": func(fullPath string) string {
			rel, _ := filepath.Rel(protoPath, fullPath)
			rel = strings.TrimSuffix(rel, ".proto")
			return rel
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

		f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
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
