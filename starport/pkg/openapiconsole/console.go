package openapiconsole

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed index.tpl
var index embed.FS

// Handler returns an http handler that servers OpenAPI console for specURL.
func Handler(title, specURL string) http.HandlerFunc {
	t, err := template.ParseFS(index, "index.tpl")
	if err != nil {
		log.Fatal(err)
	}

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
