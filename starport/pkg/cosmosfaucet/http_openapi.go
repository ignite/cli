package cosmosfaucet

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

const (
	fileNameOpenAPIIndex = "openapi/index.html"
	fileNameOpenAPISpec  = "openapi/openapi.yml.tmpl"
)

var (
	fileOpenAPIIndex = bytes.NewReader(MustAsset(fileNameOpenAPIIndex))
	tmplOpenAPISpec  = template.Must(template.
				New(fileNameOpenAPISpec).
				Parse(string(MustAsset(fileNameOpenAPISpec))))
)

func (f Faucet) openAPIIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, fileNameOpenAPIIndex, time.Now(), fileOpenAPIIndex)
}

type openAPIData struct {
	ChainID    string
	APIAddress string
}

func (f Faucet) openAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	tmplOpenAPISpec.Execute(w, f.openAPIData)
}
