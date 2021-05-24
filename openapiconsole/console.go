package openapiconsole

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed index.tpl
var index embed.FS

// Handler returns an http handler that servers OpenAPI console for an OpenAPI spec at specURL.
func Handler(title, specURL string) http.HandlerFunc {
	t, _ := template.ParseFS(index, "index.tpl")

	return func(w http.ResponseWriter, req *http.Request) {
		t.Execute(w, struct {
			Title string
			URL   string
		}{
			title,
			specURL,
		})
	}
}
