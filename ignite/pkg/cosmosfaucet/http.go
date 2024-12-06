package cosmosfaucet

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/ignite/cli/v29/ignite/pkg/openapiconsole"
)

// ServeHTTP implements http.Handler to expose the functionality of Faucet.Transfer() via HTTP.
// request/response payloads are compatible with the previous implementation at allinbits/cosmos-faucet.
func (f Faucet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.Handle("/", cors.Default().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost || r.Method == http.MethodOptions:
			f.faucetHandler(w, r)
		case r.Method == http.MethodGet:
			openapiconsole.Handler("Faucet", "openapi.yml")(w, r)
		default:
			http.NotFound(w, r)
		}
	})))

	mux.Handle("/info", cors.Default().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodOptions {
			f.faucetInfoHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})))

	mux.HandleFunc("/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			f.openAPISpecHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	mux.ServeHTTP(w, r)
}
