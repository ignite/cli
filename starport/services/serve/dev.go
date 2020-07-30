package starportserve

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xhttp"
)

const (
	appNodeInfoEndpoint = "/node_info"
)

// serviceStatus keeps the state of development environment and services needed
// for development.
type serviceStatus struct {
	Env                    env  `json:"env"`
	IsConsensusEngineAlive bool `json:"is_consensus_engine__alive"`
	IsMyAppBackendAlive    bool `json:"is_my_app_backend_alive"`
	IsMyAppFrontendAlive   bool `json:"is_my_app_frontend_alive"`
}

// env holds info about development environment.
type env struct {
	ChainID string `json:"chain_id"`
	NodeJS  bool   `json:"node_js"`
}

// serviceStatusResponse is the status response message returned to client.
type serviceStatusResponse struct {
	Status serviceStatus `json:"status"`
}

// development handler builder.
type development struct {
	app  App
	conf Config
}

// Config used to configure development handler.
type Config struct {
	EngineAddr            string
	AppBackendAddr        string
	AppFrontendAddr       string
	DevFrontendAssetsPath string
}

// newDevHandler creates a new development server handler for app by given conf.
func newDevHandler(app App, conf Config) http.Handler {
	dev := &development{app, conf}
	router := mux.NewRouter()
	router.Handle("/status", dev.statusHandler()).Methods(http.MethodGet)
	router.PathPrefix("/").Handler(dev.devAssetsHandler()).Methods(http.MethodGet)
	return router
}

func (d *development) devAssetsHandler() http.Handler {
	return http.FileServer(packr.New("ui/dist", d.conf.DevFrontendAssetsPath))
}

func (d *development) statusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := d.env()
		engineStatus, err := httpstatuschecker.Check(d.conf.EngineAddr)
		if err != nil {
			xhttp.ResponseJSON(w, http.StatusInternalServerError, xhttp.NewErrorResponse(err))
			return
		}
		appBackendStatus, err := httpstatuschecker.Check(d.conf.AppBackendAddr + appNodeInfoEndpoint)
		if err != nil {
			xhttp.ResponseJSON(w, http.StatusInternalServerError, xhttp.NewErrorResponse(err))
			return
		}
		appFrontendStatus, err := httpstatuschecker.Check(d.conf.AppFrontendAddr)
		if err != nil {
			xhttp.ResponseJSON(w, http.StatusInternalServerError, xhttp.NewErrorResponse(err))
			return
		}

		status := serviceStatus{
			Env:                    env,
			IsConsensusEngineAlive: engineStatus,
			IsMyAppBackendAlive:    appBackendStatus,
			IsMyAppFrontendAlive:   appFrontendStatus,
		}
		xhttp.ResponseJSON(w, http.StatusOK, serviceStatusResponse{status})
	})
}

func (d *development) env() env {
	return env{
		d.app.Name,
		isCommandAvailable("node"),
	}
}
