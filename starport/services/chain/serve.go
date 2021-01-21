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

	"github.com/pkg/errors"
	conf "github.com/tendermint/starport/starport/chainconf"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmoscoin"
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
)

// Serve serves an app.
func (c *Chain) Serve(ctx context.Context) error {
	// initial checks and setup.
	if err := c.setup(ctx); err != nil {
		return err
	}

	// make sure that config.yml exists.
	if _, err := conf.Locate(c.app.Path); err != nil {
		return err
	}

	// initialize the relayer if application supports it so, secret.yml
	// can be generated and watched for changes.
	if err := c.checkIBCRelayerSupport(); err == nil {
		if _, err := c.RelayerInfo(); err != nil {
			return err
		}
	}

	// start serving components.
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return c.watchAppFrontend(ctx)
	})
	g.Go(func() error {
		return c.runDevServer(ctx)
	})

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
				var (
					serveCtx context.Context
					buildErr *CannotBuildAppError
				)
				serveCtx, c.serveCancel = context.WithCancel(ctx)

				// serve the app.
				err := c.serve(serveCtx)
				switch {
				case err == nil:
				case errors.Is(err, context.Canceled):
				case errors.As(err, &buildErr):
					fmt.Fprintf(c.stdLog(logStarport).err, "%s\n", errorColor(err.Error()))

					var validationErr *conf.ValidationError
					if errors.As(err, &validationErr) {
						fmt.Fprintln(c.stdLog(logStarport).out, "see: https://github.com/tendermint/starport#configure")
					}

					fmt.Fprintf(c.stdLog(logStarport).out, "%s\n", infoColor("waiting for a fix before retrying..."))
				default:
					return err
				}
			}
		}
	})
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
		appBackendWatchPaths,
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

func (c *Chain) serve(ctx context.Context) error {
	conf, err := c.Config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

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

	// initialize the blockchain
	if err := c.Init(ctx); err != nil {
		return err
	}

	// initialize the blockchain accounts
	if err := c.InitAccounts(ctx, conf); err != nil {
		return err
	}

	// start the blockchain
	return c.start(ctx, conf)
}

func (c *Chain) start(ctx context.Context, conf conf.Config) error {
	g, ctx := errgroup.WithContext(ctx)

	// start the blockchain.
	g.Go(func() error { return c.plugin.Start(ctx, c.Commands(), conf) })

	// run relayer.
	go func() {
		if err := c.initRelayer(ctx, conf); err != nil && ctx.Err() == nil {
			fmt.Fprintf(c.stdLog(logStarport).err, "could not init relayer: %s", err)
		}
	}()

	// start the faucet if enabled.
	var isFaucedEnabled bool

	if conf.Faucet.Name != nil {
		if _, err := c.cmd.ShowAccount(ctx, *conf.Faucet.Name); err != nil {
			if err == chaincmdrunner.ErrAccountDoesNotExist {
				return &CannotBuildAppError{errors.Wrap(err, "faucet account doesn't exist")}
			}
			return err
		}

		isFaucedEnabled = true

		g.Go(func() (err error) {
			if err := c.runFaucetServer(ctx); err != nil {
				return &CannotBuildAppError{err}
			}
			return nil
		})
	}

	// print the server addresses.
	fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a Cosmos '%[1]v' app with Tendermint at %s.\n", c.app.Name, xurl.HTTP(conf.Servers.RPCAddr))
	fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a server at %s (LCD)\n", xurl.HTTP(conf.Servers.APIAddr))

	if isFaucedEnabled {
		fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a faucet at http://localhost:%d\n", conf.Faucet.Port)
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

func (c *Chain) runFaucetServer(ctx context.Context) error {
	config, err := c.Config()
	if err != nil {
		return err
	}

	faucetOptions := []cosmosfaucet.Option{
		cosmosfaucet.Account(*config.Faucet.Name, ""),
	}

	// parse coins to pass to the faucet as coins.
	for _, coin := range config.Faucet.Coins {
		amount, denom, err := cosmoscoin.Parse(coin)
		if err != nil {
			return fmt.Errorf("%s: %s", err, coin)
		}

		var amountMax uint64

		// find out the max amount for this coin.
		for _, coinMax := range config.Faucet.CoinsMax {
			amount, denomMax, err := cosmoscoin.Parse(coinMax)
			if err != nil {
				return fmt.Errorf("%s: %s", err, coin)
			}
			if denomMax == denom {
				amountMax = amount
				break
			}
		}

		faucetOptions = append(faucetOptions, cosmosfaucet.Coin(amount, amountMax, denom))
	}

	faucet, err := cosmosfaucet.New(ctx, c.cmd, faucetOptions...)
	if err != nil {
		return err
	}

	return xhttp.Serve(ctx, &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.Faucet.Port),
		Handler: faucet,
	})
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
