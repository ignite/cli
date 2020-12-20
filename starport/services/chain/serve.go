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
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/fswatcher"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/chain/conf"
	secretconf "github.com/tendermint/starport/starport/services/chain/conf/secret"
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
	fmt.Fprintf(c.stdLog(logStarport).out, "Cosmos' version is: %s\n\n", infoColor(c.plugin.Name()))

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
	sconf, err := secretconf.Open(c.app.Path)
	if err != nil {
		return err
	}

	if err := c.buildProto(ctx); err != nil {
		return err
	}

	buildSteps, err := c.buildSteps()
	if err != nil {
		return err
	}
	if err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, buildSteps...); err != nil {
		return err
	}

	if err := c.Init(ctx); err != nil {
		return err
	}

	for _, account := range conf.Accounts {
		acc, err := c.CreateAccount(ctx, account.Name, "", false)
		if err != nil {
			return err
		}

		acc.Coins = strings.Join(account.Coins, ",")
		if err := c.AddGenesisAccount(ctx, acc); err != nil {
			return err
		}
	}
	for _, account := range sconf.Accounts {
		acc, err := c.CreateAccount(ctx, account.Name, account.Mnemonic, false)
		if err != nil {
			return err
		}

		acc.Coins = strings.Join(account.Coins, ",")
		if err := c.AddGenesisAccount(ctx, acc); err != nil {
			return err
		}
	}

	setupSteps, err := c.setupSteps()
	if err != nil {
		return err
	}
	if err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, setupSteps...); err != nil {
		return err
	}
	if _, err := c.Gentx(ctx, Validator{
		Name:          conf.Validator.Name,
		StakingAmount: conf.Validator.Staked,
	}); err != nil {
		return err
	}
	if err := c.CollectGentx(ctx); err != nil {
		return err
	}

	wr := sync.WaitGroup{}
	wr.Add(1)

	go func() {
		wr.Wait()
		if err := c.initRelayer(ctx, conf); err != nil && ctx.Err() == nil {
			fmt.Fprintf(c.stdLog(logStarport).err, "could not init relayer: %s", err)
		}
	}()

	return cmdrunner.
		New(append(c.cmdOptions(), cmdrunner.RunParallel())...).
		Run(ctx, c.serverSteps(ctx, &wr, conf)...)
}

func (c *Chain) serverSteps(_ context.Context, wr *sync.WaitGroup, conf conf.Config) (steps step.Steps) {
	var wg sync.WaitGroup
	wg.Add(len(c.plugin.StartCommands(conf)))
	go func() {
		wg.Wait()
		fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a Cosmos '%[1]v' app with Tendermint at %s.\n", c.app.Name, xurl.HTTP(conf.Servers.RPCAddr))
		fmt.Fprintf(c.stdLog(logStarport).out, "🌍 Running a server at %s (LCD)\n", xurl.HTTP(conf.Servers.APIAddr))
		fmt.Fprintf(c.stdLog(logStarport).out, "\n🚀 Get started: %s\n\n", xurl.HTTP(conf.Servers.DevUIAddr))
		wr.Done()
	}()

	for _, exec := range c.plugin.StartCommands(conf) {
		steps.Add(
			step.New(
				step.NewOptions().
					Add(exec...).
					Add(
						step.InExec(func() error {
							wg.Done()
							return nil
						}),
					).
					Add(c.stdSteps(logAppd)...)...,
			),
		)
	}

	return
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
	sv := &http.Server{
		Addr:    config.Servers.DevUIAddr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		sv.Shutdown(shutdownCtx)
	}()

	err = sv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
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
