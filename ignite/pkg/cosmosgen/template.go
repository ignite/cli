package cosmosgen

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ignite/cli/v29/ignite/pkg/xstrcase"
)

var (
	//go:embed templates/*
	templates embed.FS

	templateTSClientRoot           = newTemplateWriter("root")
	templateTSClientModule         = newTemplateWriter("module")
	templateTSClientComposable     = newTemplateWriter("composable")
	templateTSClientComposableRoot = newTemplateWriter("composable-root")
)

type templateWriter struct {
	templateDir string
}

// newTemplateWriter returns a func for template residing at templatePath to initialize a text template
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
			word = xstrcase.UpperCamel(replacer.Replace(word))

			return cases.Title(language.English).String(word)
		},
		"camelCaseLowerSta": func(word string) string {
			replacer := strings.NewReplacer("-", "_", ".", "_")

			return strcase.ToLowerCamel(replacer.Replace(word))
		},
		"camelCaseUpperSta": func(word string) string {
			replacer := strings.NewReplacer("-", "_", ".", "_")

			return xstrcase.UpperCamel(replacer.Replace(word))
		},
		"resolveFile": func(fullPath string) string { // TODO to check
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
