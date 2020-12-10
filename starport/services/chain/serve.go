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

// Serve serves an app.
func (s *Chain) Serve(ctx context.Context) error {
	// initial checks and setup.
	if err := s.setup(ctx); err != nil {
		return err
	}
	// initialize the relayer if application supports it so, secret.yml
	// can be generated and watched for changes.
	if err := s.checkIBCRelayerSupport(); err == nil {
		if _, err := s.RelayerInfo(); err != nil {
			return err
		}
	}

	// start serving components.
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.watchAppFrontend(ctx)
	})
	g.Go(func() error {
		return s.runDevServer(ctx)
	})
	g.Go(func() error {
		s.refreshServe()
		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-s.serveRefresher:
				var (
					serveCtx context.Context
					buildErr *CannotBuildAppError
				)
				serveCtx, s.serveCancel = context.WithCancel(ctx)
				err := s.serve(serveCtx)
				switch {
				case err == nil:
				case errors.Is(err, context.Canceled):
				case errors.As(err, &buildErr):
					fmt.Fprintf(s.stdLog(logStarport).err, "%s\n", errorColor(err.Error()))

					var validationErr *conf.ValidationError
					if errors.As(err, &validationErr) {
						fmt.Fprintln(s.stdLog(logStarport).out, "see: https://github.com/tendermint/starport#configure")
					}

					fmt.Fprintf(s.stdLog(logStarport).out, "%s\n", infoColor("waiting for a fix before retrying..."))
				default:
					return err
				}
			}
		}
	})
	g.Go(func() error {
		return s.watchAppBackend(ctx)
	})
	return g.Wait()
}

func (s *Chain) setup(ctx context.Context) error {
	fmt.Fprintf(s.stdLog(logStarport).out, "Cosmos' version is: %s\n", infoColor(s.plugin.Name()))

	if err := s.checkSystem(); err != nil {
		return err
	}
	if err := s.plugin.Setup(ctx); err != nil {
		return err
	}
	return nil
}

// checkSystem checks if developer's work environment comply must to have
// dependencies and pre-conditions.
func (s *Chain) checkSystem() error {
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

func (s *Chain) refreshServe() {
	if s.serveCancel != nil {
		s.serveCancel()
	}
	s.serveRefresher <- struct{}{}
}

func (s *Chain) watchAppBackend(ctx context.Context) error {
	return fswatcher.Watch(
		ctx,
		appBackendWatchPaths,
		fswatcher.Workdir(s.app.Path),
		fswatcher.OnChange(s.refreshServe),
		fswatcher.IgnoreHidden(),
	)
}

func (s *Chain) cmdOptions() []cmdrunner.Option {
	return []cmdrunner.Option{
		cmdrunner.DefaultWorkdir(s.app.Path),
	}
}

func (s *Chain) serve(ctx context.Context) error {
	conf, err := s.Config()
	if err != nil {
		return &CannotBuildAppError{err}
	}
	sconf, err := secretconf.Open(s.app.Path)
	if err != nil {
		return err
	}

	buildSteps, err := s.buildSteps(ctx, conf)
	if err != nil {
		return err
	}
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, buildSteps...); err != nil {
		return err
	}

	if err := s.Init(ctx); err != nil {
		return err
	}

	for _, account := range conf.Accounts {
		acc, err := s.CreateAccount(ctx, account.Name, "", false)
		if err != nil {
			return err
		}

		acc.Coins = strings.Join(account.Coins, ",")
		if err := s.AddGenesisAccount(ctx, acc, ""); err != nil {
			return err
		}
	}
	for _, account := range sconf.Accounts {
		acc, err := s.CreateAccount(ctx, account.Name, account.Mnemonic, false)
		if err != nil {
			return err
		}

		acc.Coins = strings.Join(account.Coins, ",")
		if err := s.AddGenesisAccount(ctx, acc, ""); err != nil {
			return err
		}
	}

	setupSteps, err := s.setupSteps(ctx, conf)
	if err != nil {
		return err
	}
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, setupSteps...); err != nil {
		return err
	}
	if _, err := s.Gentx(ctx, Validator{
		Name:          conf.Validator.Name,
		StakingAmount: conf.Validator.Staked,
	}); err != nil {
		return err
	}
	if err := s.CollectGentx(ctx, ""); err != nil {
		return err
	}

	wr := sync.WaitGroup{}
	wr.Add(1)

	go func() {
		wr.Wait()
		if err := s.initRelayer(ctx, conf); err != nil && ctx.Err() == nil {
			fmt.Fprintf(s.stdLog(logStarport).err, "could not init relayer: %s", err)
		}
	}()

	return cmdrunner.
		New(append(s.cmdOptions(), cmdrunner.RunParallel())...).
		Run(ctx, s.serverSteps(ctx, &wr, conf)...)
}

func (s *Chain) serverSteps(ctx context.Context, wr *sync.WaitGroup, conf conf.Config) (steps step.Steps) {
	var wg sync.WaitGroup
	wg.Add(len(s.plugin.StartCommands(conf)))
	go func() {
		wg.Wait()
		fmt.Fprintf(s.stdLog(logStarport).out, "🌍 Running a Cosmos '%[1]v' app with Tendermint at %s.\n", s.app.Name, xurl.HTTP(conf.Servers.RPCAddr))
		fmt.Fprintf(s.stdLog(logStarport).out, "🌍 Running a server at %s (LCD)\n", xurl.HTTP(conf.Servers.APIAddr))
		fmt.Fprintf(s.stdLog(logStarport).out, "\n🚀 Get started: %s\n\n", xurl.HTTP(conf.Servers.DevUIAddr))
		wr.Done()
	}()

	for _, exec := range s.plugin.StartCommands(conf) {
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
					Add(s.stdSteps(logAppd)...)...,
			),
		)
	}

	return
}

func (s *Chain) watchAppFrontend(ctx context.Context) error {
	conf, err := s.Config()
	if err != nil {
		return err
	}
	vueFullPath := filepath.Join(s.app.Path, vuePath)
	if _, err := os.Stat(vueFullPath); os.IsNotExist(err) {
		return nil
	}
	frontendErr := &bytes.Buffer{}
	postExec := func(err error) error {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() > 0 {
			fmt.Fprintf(s.stdLog(logStarport).err, "%s\n%s",
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

func (s *Chain) runDevServer(ctx context.Context) error {
	c, err := s.Config()
	if err != nil {
		return err
	}

	grpcconn, grpcHandler, err := newGRPCWebProxyHandler(c.Servers.GRPCAddr)
	if err != nil {
		return err
	}
	defer grpcconn.Close()

	conf := Config{
		SdkVersion:      s.plugin.Name(),
		EngineAddr:      xurl.HTTP(c.Servers.RPCAddr),
		AppBackendAddr:  xurl.HTTP(c.Servers.APIAddr),
		AppFrontendAddr: xurl.HTTP(c.Servers.FrontendAddr),
	} // TODO get vals from const
	handler, err := newDevHandler(s.app, conf, grpcHandler)
	if err != nil {
		return err
	}
	sv := &http.Server{
		Addr:    c.Servers.DevUIAddr,
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
