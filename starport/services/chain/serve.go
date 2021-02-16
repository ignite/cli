package chain

import (
	"bytes"
	"context"
	"fmt"
	"go/build"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tendermint/starport/starport/pkg/dirchange"

	"github.com/tendermint/starport/starport/services"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	conf "github.com/tendermint/starport/starport/chainconf"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/fswatcher"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/pkg/xhttp"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"golang.org/x/sync/errgroup"
)

var (
	// ignoredExts holds a list of ignored files from watching.
	ignoredExts = []string{"pb.go", "pb.gw.go"}

	// chainSavePath is the place where chain exported genesis are saved
	chainSavePath = filepath.Join(services.StarportConfDir, "local-chains")

	// exportedGenesis is the name of the exported genesis file for a chain
	exportedGenesis = "exported_genesis.json"

	// sourceChecksum is the file containing the checksum to detect source modification
	sourceChecksum = "source_checksum.txt"

	// binaryChecksum is the file containing the checksum to detect binary modification
	binaryChecksum = "binary_checksum.txt"

	// configChecksum is the file containing the checksum to detect config modification
	configChecksum = "config_checksum.txt"
)

type serveOptions struct {
	forceReset bool
	resetOnce  bool
}

func newServeOption() serveOptions {
	return serveOptions{
		forceReset: false,
		resetOnce:  false,
	}
}

// ServeOption provides options for the serve command
type ServeOption func(*serveOptions)

// ServeForceReset allows to force reset of the state when the chain is served and on every source change
func ServeForceReset() ServeOption {
	return func(c *serveOptions) {
		c.forceReset = true
	}
}

// ServeResetOnce allows to reset of the state when the chain is served once
func ServeResetOnce() ServeOption {
	return func(c *serveOptions) {
		c.resetOnce = true
	}
}

// Serve serves an app.
func (c *Chain) Serve(ctx context.Context, options ...ServeOption) error {
	serveOptions := newServeOption()

	// apply the options
	for _, apply := range options {
		apply(&serveOptions)
	}

	// initial checks and setup.
	if err := c.setup(ctx); err != nil {
		return err
	}

	// make sure that config.yml exists.
	if c.options.ConfigName != "" {
		if _, err := os.Stat(filepath.Join(c.app.Path, c.options.ConfigName)); err != nil {
			return err
		}
	} else if _, err := conf.LocateDefault(c.app.Path); err != nil {
		return err
	}

	// start serving components.
	g, ctx := errgroup.WithContext(ctx)

	// routine to watch front-end
	g.Go(func() error {
		return c.watchAppFrontend(ctx)
	})

	// development server routine
	g.Go(func() error {
		return c.runDevServer(ctx)
	})

	// blockchain node routine
	g.Go(func() error {
		c.refreshServe()

		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-c.serveRefresher:
				commands, err := c.Commands(ctx)
				if err != nil {
					return err
				}

				var (
					serveCtx context.Context
					buildErr *CannotBuildAppError
					startErr *CannotStartAppError
				)
				serveCtx, c.serveCancel = context.WithCancel(ctx)

				// determine if the chain should reset the state
				shouldReset := serveOptions.forceReset || serveOptions.resetOnce

				// serve the app.
				err = c.serve(serveCtx, shouldReset)
				serveOptions.resetOnce = false

				switch {
				case err == nil:
				case errors.Is(err, context.Canceled):
					// If the app has been served, we save the genesis state
					if c.served {
						c.served = false

						fmt.Fprintln(c.stdLog(logStarport).out, "💿 Saving genesis state...")

						// If serve has been stopped, save the genesis state
						if err := c.saveChainState(context.TODO(), commands); err != nil {
							fmt.Fprint(c.stdLog(logStarport).err, err.Error())
							return err
						}

						genesisPath, err := c.exportedGenesisPath()
						if err != nil {
							fmt.Fprintln(c.stdLog(logStarport).err, err.Error())
							return err
						}
						fmt.Fprintf(c.stdLog(logStarport).out, "💿 Genesis state saved in %s\n", genesisPath)
					}
				case errors.As(err, &buildErr):
					fmt.Fprintf(c.stdLog(logStarport).err, "%s\n", errorColor(err.Error()))

					var validationErr *conf.ValidationError
					if errors.As(err, &validationErr) {
						fmt.Fprintln(c.stdLog(logStarport).out, "see: https://github.com/tendermint/starport#configure")
					}

					fmt.Fprintf(c.stdLog(logStarport).out, "%s\n", infoColor("Waiting for a fix before retrying..."))

				case errors.As(err, &startErr):

					// Parse returned error logs
					parsedErr := startErr.ParseStartError()

					// If empty, we cannot recognized the error
					// Therefore, the error may be caused by a new logic that is not compatible with the old app state
					// We suggest the user to eventually reset the app state
					if parsedErr == "" {
						fmt.Fprintf(c.stdLog(logStarport).out, "%s %s\n", infoColor(`Blockchain failed to start.
If the new code is no longer compatible with the saved state, you can reset the database by launching:`), "starport serve --reset-once")

						return fmt.Errorf("cannot run %s", startErr.AppName)
					}

					// return the clear parsed error
					return errors.New(parsedErr)
				default:
					return err
				}
			}
		}
	})

	// routine to watch back-end
	g.Go(func() error {
		return c.watchAppBackend(ctx)
	})

	return g.Wait()
}

func (c *Chain) setup(ctx context.Context) error {
	fmt.Fprintf(c.stdLog(logStarport).out, "Cosmos SDK's version is: %s\n\n", infoColor(c.Version))

	if err := c.checkSystem(); err != nil {
		return err
	}
	if err := c.plugin.Setup(ctx); err != nil {
		return err
	}
	return nil
}

// checkSystem checks if developer's work environment comply must to have
// dependencies and pre-conditions.
func (c *Chain) checkSystem() error {
	// check if Go has installed.
	if !xexec.IsCommandAvailable("go") {
		return errors.New("Please, check that Go language is installed correctly in $PATH. See https://golang.org/doc/install")
	}
	// check if Go's bin added to System's path.
	gobinpath := path.Join(build.Default.GOPATH, "bin")
	if err := xos.IsInPath(gobinpath); err != nil {
		return errors.New("$(go env GOPATH)/bin must be added to your $PATH. See https://golang.org/doc/gopath_code.html#GOPATH")
	}
	return nil
}

func (c *Chain) refreshServe() {
	if c.serveCancel != nil {
		c.serveCancel()
	}
	c.serveRefresher <- struct{}{}
}

func (c *Chain) watchAppBackend(ctx context.Context) error {
	return fswatcher.Watch(
		ctx,
		append(appBackendSourceWatchPaths, c.AppBackendConfigWatchPaths()...),
		fswatcher.Workdir(c.app.Path),
		fswatcher.OnChange(c.refreshServe),
		fswatcher.IgnoreHidden(),
		fswatcher.IgnoreExt(ignoredExts...),
	)
}

func (c *Chain) cmdOptions() []cmdrunner.Option {
	return []cmdrunner.Option{
		cmdrunner.DefaultWorkdir(c.app.Path),
	}
}

// serve performs the operations to serve the blockchain: build, init and start
// if the chain is already initialized and the file didn't changed, the app is directly started
// if the files changed, the state is imported
func (c *Chain) serve(ctx context.Context, forceReset bool) error {
	conf, err := c.Config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	saveDir, err := c.chainSavePath()
	if err != nil {
		return err
	}

	// isInit determines if the app is initialized
	var isInit bool

	// determine if the app must reset the state
	// if the state must be reset, then we consider the chain as being not initialized
	isInit, err = c.IsInitialized()
	if err != nil {
		return err
	}
	if isInit {
		configModified, err := dirchange.HasDirChecksumChanged(c.app.Path, c.AppBackendConfigWatchPaths(), saveDir, configChecksum)
		if err != nil {
			return err
		}

		if forceReset || configModified {
			// if forceReset is set, we consider the app as being not initialized
			fmt.Fprintln(c.stdLog(logStarport).out, "🔄 Resetting the app state...")
			isInit = false
		}
	}

	// check if source has been modified since last serve
	// if the state must not be reset but the source has changed, we rebuild the chain and import the exported state
	sourceModified, err := dirchange.HasDirChecksumChanged(c.app.Path, appBackendSourceWatchPaths, saveDir, sourceChecksum)
	if err != nil {
		return err
	}

	// we also consider the binary in the checksum to ensure the binary has not been changed by a third party
	var binaryModified bool
	binaryName, err := c.Binary()
	if err != nil {
		return err
	}
	binaryPath, err := exec.LookPath(binaryName)
	if err != nil {
		if !errors.Is(err, exec.ErrNotFound) {
			return err
		}
		binaryModified = true
	} else {
		binaryModified, err = dirchange.HasDirChecksumChanged("", []string{binaryPath}, saveDir, binaryChecksum)
		if err != nil {
			return err
		}
	}

	appModified := sourceModified || binaryModified

	// check if exported genesis exists
	exportGenesisExists := true
	exportedGenesisPath, err := c.exportedGenesisPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(exportedGenesisPath); os.IsNotExist(err) {
		exportGenesisExists = false
	} else if err != nil {
		return err
	}

	// build phase
	if !isInit || appModified {
		// build proto
		if err := c.buildProto(ctx); err != nil {
			return err
		}

		// build the blockchain app
		buildSteps, err := c.buildSteps()
		if err != nil {
			return err
		}
		if err := cmdrunner.
			New(c.cmdOptions()...).
			Run(ctx, buildSteps...); err != nil {
			return err
		}
	}

	// init phase
	// nolint:gocritic
	if !isInit || (appModified && !exportGenesisExists) {
		fmt.Fprintln(c.stdLog(logStarport).out, "💿 Initializing the app...")

		// initialize the blockchain
		if err := c.Init(ctx); err != nil {
			return err
		}

		// initialize the blockchain accounts
		if err := c.InitAccounts(ctx, conf); err != nil {
			return err
		}
	} else if appModified {
		// if the chain is already initialized but the source has been modified
		// we reset the chain database and import the genesis state
		fmt.Fprintln(c.stdLog(logStarport).out, "💿 Existent genesis detected, restoring the database...")

		if err := commands.UnsafeReset(ctx); err != nil {
			return err
		}

		if err := c.importChainState(); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(c.stdLog(logStarport).out, "▶️  Restarting existing app...")
	}

	// save checksums
	if err := dirchange.SaveDirChecksum(c.app.Path, c.AppBackendConfigWatchPaths(), saveDir, configChecksum); err != nil {
		return err
	}
	if err := dirchange.SaveDirChecksum(c.app.Path, appBackendSourceWatchPaths, saveDir, sourceChecksum); err != nil {
		return err
	}
	binaryPath, err = exec.LookPath(binaryName)
	if err != nil {
		return err
	}
	if err := dirchange.SaveDirChecksum("", []string{binaryPath}, saveDir, binaryChecksum); err != nil {
		return err
	}

	// start the blockchain
	return c.start(ctx, conf)
}

func (c *Chain) start(ctx context.Context, conf conf.Config) error {
	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	// start the blockchain.
	g.Go(func() error { return c.plugin.Start(ctx, commands, conf) })

	// start the faucet if enabled.
	faucet, err := c.Faucet(ctx)
	isFaucetEnabled := err != ErrFaucetIsNotEnabled

	if isFaucetEnabled {
		if err == ErrFaucetAccountDoesNotExist {
			return &CannotBuildAppError{errors.Wrap(err, "faucet account doesn't exist")}
		}
		if err != nil {
			return err
		}

		g.Go(func() (err error) {
			if err := c.runFaucetServer(ctx, faucet); err != nil {
				return &CannotBuildAppError{err}
			}
			return nil
		})
	}

	// set the app as being served
	c.served = true

	// print the server addresses.
	fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a Cosmos '%[1]v' app with Tendermint at %s.\n", c.app.Name, xurl.HTTP(conf.Servers.RPCAddr))
	fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a server at %s (LCD)\n", xurl.HTTP(conf.Servers.APIAddr))

	if isFaucetEnabled {
		fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a faucet at http://0.0.0.0:%d\n", conf.Faucet.Port)
	}

	fmt.Fprintf(c.stdLog(logStarport).out, "\n🚀 Get started: %s\n\n", xurl.HTTP(conf.Servers.DevUIAddr))

	return g.Wait()
}

func (c *Chain) watchAppFrontend(ctx context.Context) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}
	vueFullPath := filepath.Join(c.app.Path, vuePath)
	if _, err := os.Stat(vueFullPath); os.IsNotExist(err) {
		return nil
	}
	frontendErr := &bytes.Buffer{}
	postExec := func(err error) error {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() > 0 {
			fmt.Fprintf(c.stdLog(logStarport).err, "%s\n%s",
				infoColor("skipping serving Vue frontend due to following errors:"), errorColor(frontendErr.String()))
		}
		return nil // ignore errors.
	}
	host, port, err := net.SplitHostPort(conf.Servers.FrontendAddr)
	if err != nil {
		return err
	}
	return cmdrunner.
		New(
			cmdrunner.DefaultWorkdir(vueFullPath),
			cmdrunner.DefaultStderr(frontendErr),
		).
		Run(ctx,
			step.New(
				step.Exec("npm", "i"),
				step.PostExec(postExec),
			),
			step.New(
				step.Exec("npm", "run", "serve"),
				step.Env(
					fmt.Sprintf("HOST=%s", host),
					fmt.Sprintf("PORT=%s", port),
					fmt.Sprintf("VUE_APP_API_COSMOS=%s", xurl.HTTP(conf.Servers.APIAddr)),
					fmt.Sprintf("VUE_APP_API_TENDERMINT=%s", xurl.HTTP(conf.Servers.RPCAddr)),
					fmt.Sprintf("VUE_APP_WS_TENDERMINT=%s/websocket", xurl.WS(conf.Servers.RPCAddr)),
				),
				step.PostExec(postExec),
			),
		)
}

func (c *Chain) runDevServer(ctx context.Context) error {
	config, err := c.Config()
	if err != nil {
		return err
	}

	grpcconn, grpcHandler, err := newGRPCWebProxyHandler(config.Servers.GRPCAddr)
	if err != nil {
		return err
	}
	defer grpcconn.Close()

	conf := Config{
		SdkVersion:      c.plugin.Name(),
		EngineAddr:      xurl.HTTP(config.Servers.RPCAddr),
		AppBackendAddr:  xurl.HTTP(config.Servers.APIAddr),
		AppFrontendAddr: xurl.HTTP(config.Servers.FrontendAddr),
	} // TODO get vals from const
	handler, err := newDevHandler(c.app, conf, grpcHandler)
	if err != nil {
		return err
	}

	return xhttp.Serve(ctx, &http.Server{
		Addr:    config.Servers.DevUIAddr,
		Handler: handler,
	})
}

func (c *Chain) runFaucetServer(ctx context.Context, faucet cosmosfaucet.Faucet) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	return xhttp.Serve(ctx, &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", conf.Faucet.Port),
		Handler: faucet,
	})
}

// saveChainState runs the export command of the chain and store the exported genesis in the chain saved config
func (c *Chain) saveChainState(ctx context.Context, commands chaincmdrunner.Runner) error {
	genesisPath, err := c.exportedGenesisPath()
	if err != nil {
		return err
	}

	return commands.Export(ctx, genesisPath)
}

// importChainState imports the saved genesis in chain config to use it as the genesis
func (c *Chain) importChainState() error {
	exportGenesisPath, err := c.exportedGenesisPath()
	if err != nil {
		return err
	}
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	return copy.Copy(exportGenesisPath, genesisPath)
}

// chainSavePath returns the path where the chain state is saved
// create the path if it doesn't exist
func (c *Chain) chainSavePath() (string, error) {
	chainID, err := c.ID()
	if err != nil {
		return "", err
	}
	savePath := filepath.Join(chainSavePath, chainID)

	// ensure the path exists
	if err := os.MkdirAll(savePath, 0700); err != nil && !os.IsExist(err) {
		return "", err
	}

	return savePath, nil
}

// exportedGenesisPath returns the path of the exported genesis file
func (c *Chain) exportedGenesisPath() (string, error) {
	savePath, err := c.chainSavePath()
	if err != nil {
		return "", err
	}

	return filepath.Join(savePath, exportedGenesis), nil
}

type CannotBuildAppError struct {
	Err error
}

func (e *CannotBuildAppError) Error() string {
	return fmt.Sprintf("cannot build app:\n\n\t%s", e.Err)
}

func (e *CannotBuildAppError) Unwrap() error {
	return e.Err
}

type CannotStartAppError struct {
	AppName string
	Err     error
}

func (e *CannotStartAppError) Error() string {
	return fmt.Sprintf("cannot run %sd start:\n%s", e.AppName, errors.Unwrap(e.Err))
}

func (e *CannotStartAppError) Unwrap() error {
	return e.Err
}

// ParseStartError parses the error into a clear error string
// The error logs from Cosmos SDK application are too extensive to be directly printed
// If the error is not recognized, returns an empty string
func (e *CannotStartAppError) ParseStartError() string {
	errorLogs := errors.Unwrap(e.Err).Error()
	switch {
	case strings.Contains(errorLogs, "bind: address already in use"):
		r := regexp.MustCompile(`listen .* bind: address already in use`)
		return r.FindString(errorLogs)
	case strings.Contains(errorLogs, "validator set is nil in genesis"):
		return "Error: error during handshake: error on replay: validator set is nil in genesis and still empty after InitChain"
	default:
		return ""
	}
}
