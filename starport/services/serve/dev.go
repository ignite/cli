package starportserve

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/pkg/xhttp"
	"golang.org/x/sync/errgroup"

	// register dev ui for statik filesystem.
	_ "github.com/tendermint/starport/starport/ui/dist-go/statik"
)

const (
	appNodeInfoEndpoint = "/node_info"
)

// serviceStatusResponse holds the status of development environment and http services
// needed for development.
type statusResponse struct {
	Status serviceStatus `json:"status"`
	Env    env           `json:"env"`
}

// serviceStatus holds the availibity status of http services.
type serviceStatus struct {
	IsConsensusEngineAlive bool `json:"is_consensus_engine_alive"`
	IsMyAppBackendAlive    bool `json:"is_my_app_backend_alive"`
	IsMyAppFrontendAlive   bool `json:"is_my_app_frontend_alive"`
}

// env holds info about development environment.
type env struct {
	ChainID   string `json:"chain_id"`
	NodeJS    bool   `json:"node_js"`
	CustomURL string `json:"custom_url"`
}

// development handler builder.
type development struct {
	app  App
	conf Config
	uifs http.FileSystem
}

// Config used to configure development handler.
type Config struct {
	EngineAddr      string
	AppBackendAddr  string
	AppFrontendAddr string
}

// newDevHandler creates a new development server handler for app by given conf.
func newDevHandler(app App, conf Config) (http.Handler, error) {
	uifs, err := fs.New()
	if err != nil {
		return nil, err
	}
	dev := &development{
		app:  app,
		conf: conf,
		uifs: uifs,
	}
	router := mux.NewRouter()
	router.Handle("/status", dev.statusHandler()).Methods(http.MethodGet)
	router.PathPrefix("/").Handler(dev.devAssetsHandler()).Methods(http.MethodGet)
	return router, nil
}

func (d *development) devAssetsHandler() http.Handler {
	return http.FileServer(d.uifs)
}

func (d *development) statusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			engineStatus,
			appBackendStatus,
			appFrontendStatus bool
		)
		g := &errgroup.Group{}
		g.Go(func() (err error) {
			engineStatus, err = httpstatuschecker.Check(d.conf.EngineAddr)
			return
		})
		g.Go(func() (err error) {
			appBackendStatus, err = httpstatuschecker.Check(d.conf.AppBackendAddr + appNodeInfoEndpoint)
			return
		})
		g.Go(func() (err error) {
			appFrontendStatus, err = httpstatuschecker.Check(d.conf.AppFrontendAddr)
			return
		})
		if err := g.Wait(); err != nil {
			xhttp.ResponseJSON(w, http.StatusInternalServerError, xhttp.NewErrorResponse(err))
			return
		}

		resp := statusResponse{
			Env: d.env(),
			Status: serviceStatus{
				IsConsensusEngineAlive: engineStatus,
				IsMyAppBackendAlive:    appBackendStatus,
				IsMyAppFrontendAlive:   appFrontendStatus,
			},
		}
		xhttp.ResponseJSON(w, http.StatusOK, resp)
	})
}

func (d *development) env() env {
	return env{
		d.app.Name,
		xexec.IsCommandAvailable("node"),
		os.Getenv("CUSTOM_URL"),
	}
}
