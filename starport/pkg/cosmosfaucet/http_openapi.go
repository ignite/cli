package cosmosfaucet

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

const (
	fileOpenAPIIndex = "openapi/index.html"
	fileOpenAPISpec  = "openapi/openapi.yml.tmpl"
)

func (f Faucet) openAPIIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, fileOpenAPIIndex, time.Now(), bytes.NewReader(MustAsset(fileOpenAPIIndex)))
}

type openAPIData struct {
	ChainID    string
	APIAddress string
}

func (f Faucet) openAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	t := template.
		Must(template.
			New(fileOpenAPISpec).
			Parse(string(MustAsset(fileOpenAPISpec))))

	t.Execute(w, f.openAPIData)
}
