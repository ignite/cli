package starportserve

import (
	"encoding/json"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

// newDevHandler creates a new development server handler.
func newDevHandler(app App) http.Handler {
	router := mux.NewRouter()
	devUI := packr.New("ui/dist", "../../ui/dist")
	router.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		env := Env{app.Name, isCommandAvailable("node")}
		js, err := json.Marshal(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:1317/node_info")
		if err != nil || res.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 error"))
		} else if res.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 ok"))
		}
	})
	router.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:26657")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if res.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	router.HandleFunc("/frontend", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://localhost:8080")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if res.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	router.PathPrefix("/").Handler(http.FileServer(devUI))
	return router
}
