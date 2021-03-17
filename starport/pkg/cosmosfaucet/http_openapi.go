package cosmosfaucet

import (
	"bytes"
	_ "embed" // used for embedding openapi assets.
	"html/template"
	"net/http"
	"time"
)

const (
	fileNameOpenAPIIndex = "openapi/index.html"
	fileNameOpenAPISpec  = "openapi/openapi.yml.tmpl"
)

var (
	//go:embed openapi/index.html
	bytesOpenAPIIndex []byte

	//go:embed openapi/openapi.yml.tmpl
	bytesOpenAPISpec []byte
)

var (
	fileOpenAPIIndex = bytes.NewReader(bytesOpenAPIIndex)
	tmplOpenAPISpec  = template.Must(template.New(fileNameOpenAPISpec).Parse(string(bytesOpenAPISpec)))
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
