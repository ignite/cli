package chain

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/pkg/cache"
	chaincmdrunner "github.com/ignite-hq/cli/ignite/pkg/chaincmd/runner"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosfaucet"
	"github.com/ignite-hq/cli/ignite/pkg/dirchange"
	"github.com/ignite-hq/cli/ignite/pkg/localfs"
	"github.com/ignite-hq/cli/ignite/pkg/xexec"
	"github.com/ignite-hq/cli/ignite/pkg/xfilepath"
	"github.com/ignite-hq/cli/ignite/pkg/xhttp"
	"github.com/ignite-hq/cli/ignite/pkg/xurl"
)

const (
	// exportedGenesis is the name of the exported genesis file for a chain
	exportedGenesis = "exported_genesis.json"

	// sourceChecksumKey is the cache key for the checksum to detect source modification
	sourceChecksumKey = "source_checksum"

	// binaryChecksumKey is the cache key for the checksum to detect binary modification
	binaryChecksumKey = "binary_checksum"

	// configChecksumKey is the cache key for containing the checksum to detect config modification
	configChecksumKey = "config_checksum"

	// serveDirchangeCacheNamespace is the name of the cache namespace for detecting changes in directories
	serveDirchangeCacheNamespace = "serve.dirchange"
)

var (
	// ignoredExts holds a list of ignored files from watching.
	ignoredExts = []string{"pb.go", "pb.gw.go"}

	// starportSavePath is the place where chain exported genesis are saved
	starportSavePath = xfilepath.Join(
		chainconfig.ConfigDirPath,
		xfilepath.Path("local-chains"),
	)
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
func (c *Chain) Serve(ctx context.Context, cacheStorage cache.Storage, options ...ServeOption) error {
	serveOptions := newServeOption()

	// apply the options
	for _, apply := range options {
		apply(&serveOptions)
	}

	// initial checks and setup.
	if err := c.setup(); err != nil {
		return err
	}

	// make sure that config.yml exists
	if c.options.ConfigFile != "" {
		if _, err := os.Stat(c.options.ConfigFile); err != nil {
			return err
		}
	} else if _, err := chainconfig.LocateDefault(c.app.Path); err != nil {
		return err
	}

	// start serving components.
	g, ctx := errgroup.WithContext(ctx)

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
				err = c.serve(serveCtx, cacheStorage, shouldReset)
				serveOptions.resetOnce = false

				switch {
				case err == nil:
				case errors.Is(err, context.Canceled):
					// If the app has been served, we save the genesis state
					if c.served {
						c.served = false

						fmt.Fprintln(c.stdLog().out, "💿 Saving genesis state...")

						// If serve has been stopped, save the genesis state
						if err := c.saveChainState(context.TODO(), commands); err != nil {
							fmt.Fprint(c.stdLog().err, err.Error())
							return err
						}

						genesisPath, err := c.exportedGenesisPath()
						if err != nil {
							fmt.Fprintln(c.stdLog().err, err.Error())
							return err
						}
						fmt.Fprintf(c.stdLog().out, "💿 Genesis state saved in %s\n", genesisPath)
					}
				case errors.As(err, &buildErr):
					fmt.Fprintf(c.stdLog().err, "%s\n", errorColor(err.Error()))

					var validationErr *chainconfig.ValidationError
					if errors.As(err, &validationErr) {
						fmt.Fprintln(c.stdLog().out, "see: https://github.com/ignite-hq/cli#configure")
					}

					fmt.Fprintf(c.stdLog().out, "%s\n", infoColor("Waiting for a fix before retrying..."))

				case errors.As(err, &startErr):
					// Parse returned error logs
					parsedErr := startErr.ParseStartError()

					// If empty, we cannot recognized the error
					// Therefore, the error may be caused by a new logic that is not compatible with the old app state
					// We suggest the user to eventually reset the app state
					if parsedErr == "" {
						fmt.Fprintf(c.stdLog().out, "%s %s\n", infoColor(`Blockchain failed to start.
If the new code is no longer compatible with the saved state, you can reset the database by launching:`), "ignite chain serve --reset-once")

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

func (c *Chain) setup() error {
	fmt.Fprintf(c.stdLog().out, "Cosmos SDK's version is: %s\n\n", infoColor(c.Version))

	return c.checkSystem()
}

// checkSystem checks if developer's work environment comply must to have
// dependencies and pre-conditions.
func (c *Chain) checkSystem() error {
	// check if Go has installed.
	if !xexec.IsCommandAvailable("go") {
		return errors.New("Please, check that Go language is installed correctly in $PATH. See https://golang.org/doc/install")
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
	watchPaths := appBackendSourceWatchPaths
	if c.ConfigPath() != "" {
		watchPaths = append(watchPaths, c.ConfigPath())
	}

	return localfs.Watch(
		ctx,
		watchPaths,
		localfs.WatcherWorkdir(c.app.Path),
		localfs.WatcherOnChange(c.refreshServe),
		localfs.WatcherIgnoreHidden(),
		localfs.WatcherIgnoreFolders(),
		localfs.WatcherIgnoreExt(ignoredExts...),
	)
}

// serve performs the operations to serve the blockchain: build, init and start
// if the chain is already initialized and the file didn't changed, the app is directly started
// if the files changed, the state is imported
func (c *Chain) serve(ctx context.Context, cacheStorage cache.Storage, forceReset bool) error {
	conf, err := c.Config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	// isInit determines if the app is initialized
	var isInit bool

	dirCache := cache.New[[]byte](cacheStorage, serveDirchangeCacheNamespace)

	// determine if the app must reset the state
	// if the state must be reset, then we consider the chain as being not initialized
	isInit, err = c.IsInitialized()
	if err != nil {
		return err
	}
	if isInit {
		configModified := false
		if c.ConfigPath() != "" {
			configModified, err = dirchange.HasDirChecksumChanged(dirCache, configChecksumKey, c.app.Path, c.ConfigPath())
			if err != nil {
				return err
			}
		}

		if forceReset || configModified {
			// if forceReset is set, we consider the app as being not initialized
			fmt.Fprintln(c.stdLog().out, "🔄 Resetting the app state...")
			isInit = false
		}
	}

	// check if source has been modified since last serve
	// if the state must not be reset but the source has changed, we rebuild the chain and import the exported state
	sourceModified, err := dirchange.HasDirChecksumChanged(dirCache, sourceChecksumKey, c.app.Path, appBackendSourceWatchPaths...)
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
		binaryModified, err = dirchange.HasDirChecksumChanged(dirCache, binaryChecksumKey, "", binaryPath)
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
		// build the blockchain app
		if err := c.build(ctx, cacheStorage, ""); err != nil {
			return err
		}
	}

	// init phase
	// nolint:gocritic
	if !isInit || (appModified && !exportGenesisExists) {
		fmt.Fprintln(c.stdLog().out, "💿 Initializing the app...")

		if err := c.Init(ctx, true); err != nil {
			return err
		}
	} else if appModified {
		// if the chain is already initialized but the source has been modified
		// we reset the chain database and import the genesis state
		fmt.Fprintln(c.stdLog().out, "💿 Existent genesis detected, restoring the database...")

		if err := commands.UnsafeReset(ctx); err != nil {
			return err
		}

		if err := c.importChainState(); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(c.stdLog().out, "▶️  Restarting existing app...")
	}

	// save checksums
	if c.ConfigPath() != "" {
		if err := dirchange.SaveDirChecksum(dirCache, configChecksumKey, c.app.Path, c.ConfigPath()); err != nil {
			return err
		}
	}
	if err := dirchange.SaveDirChecksum(dirCache, sourceChecksumKey, c.app.Path, appBackendSourceWatchPaths...); err != nil {
		return err
	}
	binaryPath, err = exec.LookPath(binaryName)
	if err != nil {
		return err
	}
	if err := dirchange.SaveDirChecksum(dirCache, binaryChecksumKey, "", binaryPath); err != nil {
		return err
	}

	// start the blockchain
	return c.start(ctx, conf)
}

func (c *Chain) start(ctx context.Context, config chainconfig.Config) error {
	commands, err := c.Commands(ctx)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	// start the blockchain.
	g.Go(func() error { return c.plugin.Start(ctx, commands, config) })

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

	// note: address format errors are handled by the
	// error group, so they can be safely ignored here
	rpcAddr, _ := xurl.HTTP(config.Host.RPC)
	apiAddr, _ := xurl.HTTP(config.Host.API)

	// print the server addresses.
	fmt.Fprintf(c.stdLog().out, "🌍 Tendermint node: %s\n", rpcAddr)
	fmt.Fprintf(c.stdLog().out, "🌍 Blockchain API: %s\n", apiAddr)

	if isFaucetEnabled {
		faucetAddr, _ := xurl.HTTP(chainconfig.FaucetHost(config))
		fmt.Fprintf(c.stdLog().out, "🌍 Token faucet: %s\n", faucetAddr)
	}

	return g.Wait()
}

func (c *Chain) runFaucetServer(ctx context.Context, faucet cosmosfaucet.Faucet) error {
	config, err := c.Config()
	if err != nil {
		return err
	}

	return xhttp.Serve(ctx, &http.Server{
		Addr:    chainconfig.FaucetHost(config),
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
	savePath, err := starportSavePath()
	if err != nil {
		return "", err
	}

	chainID, err := c.ID()
	if err != nil {
		return "", err
	}
	chainSavePath := filepath.Join(savePath, chainID)

	// ensure the path exists
	if err := os.MkdirAll(savePath, 0700); err != nil && !os.IsExist(err) {
		return "", err
	}

	return chainSavePath, nil
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
