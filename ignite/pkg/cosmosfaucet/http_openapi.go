package cosmosfaucet

import (
	_ "embed" // used for embedding openapi assets.
	"html/template"
	"net/http"
)

const (
	fileNameOpenAPISpec = "openapi/openapi.yml.tmpl"
)

//go:embed openapi/openapi.yml.tmpl
var bytesOpenAPISpec []byte

var tmplOpenAPISpec = template.Must(template.New(fileNameOpenAPISpec).Parse(string(bytesOpenAPISpec)))

type openAPIData struct {
	ChainID    string
	APIAddress string
}

func (f Faucet) openAPISpecHandler(w http.ResponseWriter, _ *http.Request) {
	tmplOpenAPISpec.Execute(w, f.openAPIData)
}
