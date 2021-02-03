package cosmosfaucet

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// ServeHTTP implements http.Handler to expose the functionality of Faucet.Transfer() via HTTP.
// request/response payloads are compatible with the previous implementation at allinbits/cosmos-faucet.
func (f Faucet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.Handle("/", cors.Default().Handler(http.HandlerFunc(f.faucetHandler))).
		Methods(http.MethodPost)

	router.HandleFunc("/", f.openAPIIndexHandler).
		Methods(http.MethodGet)

	router.HandleFunc("/openapi.yml", f.openAPISpecHandler).
		Methods(http.MethodGet)

	router.ServeHTTP(w, r)
}
