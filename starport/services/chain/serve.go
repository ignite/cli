package chain

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/confile"
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
	conf, err := s.config()
	if err != nil {
		return &CannotBuildAppError{err}
	}

	buildSteps, _ := s.buildSteps(ctx, conf)
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, buildSteps...); err != nil {
		return err
	}

	initSteps, err := s.initSteps(ctx, conf)
	if err != nil {
		return err
	}
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, initSteps...); err != nil {
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

func (s *Chain) initSteps(ctx context.Context, conf conf.Config) (steps step.Steps, err error) {
	chainID, err := s.ID()
	if err != nil {
		return nil, err
	}

	sconf, err := secretconf.Open(s.app.Path)
	if err != nil {
		return nil, err
	}

	// cleanup persistent data from previous `serve`.
	steps.Add(step.New(
		step.PreExec(func() error {
			for _, path := range s.plugin.StoragePaths() {
				if err := xos.RemoveAllUnderHome(path); err != nil {
					return err
				}
			}
			return nil
		}),
	))

	// init node.
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				s.app.d(),
				"init",
				"mynode",
				"--chain-id", chainID,
			),
			// overwrite configuration changes from Starport's config.yml to
			// over app's sdk configs.
			step.PostExec(func(err error) error {
				if err != nil {
					return err
				}

				appconfigs := []struct {
					ec      confile.EncodingCreator
					path    string
					changes map[string]interface{}
				}{
					{confile.DefaultJSONEncodingCreator, s.GenesisPath(), conf.Genesis},
					{confile.DefaultTOMLEncodingCreator, s.AppTOMLPath(), conf.Init.App},
					{confile.DefaultTOMLEncodingCreator, s.ConfigTOMLPath(), conf.Init.Config},
				}

				for _, ac := range appconfigs {
					cf := confile.New(ac.ec, ac.path)
					var conf map[string]interface{}
					if err := cf.Load(&conf); err != nil {
						return err
					}
					if err := mergo.Merge(&conf, ac.changes, mergo.WithOverride); err != nil {
						return err
					}
					if err := cf.Save(conf); err != nil {
						return err
					}
				}
				return nil
			}),
			step.PostExec(func(err error) error {
				if err != nil {
					return err
				}
				return s.plugin.PostInit(conf)
			}),
		).
		Add(s.stdSteps(logAppd)...)...,
	))

	for _, account := range conf.Accounts {
		steps.Add(s.createAccountSteps(ctx, account.Name, "", account.Coins, false)...)
	}

	for _, account := range sconf.Accounts {
		steps.Add(s.createAccountSteps(ctx, account.Name, account.Mnemonic, account.Coins, false)...)
	}

	if err := s.checkIBCRelayerSupport(); err == nil {
		steps.Add(step.New(
			step.PreExec(func() error {
				if err := xos.RemoveAllUnderHome(".relayer"); err != nil {
					return err
				}
				info, err := s.RelayerInfo()
				if err != nil {

					return err
				}
				fmt.Fprintf(s.stdLog(logStarport).out, "✨ Relayer info: %s\n", info)
				return nil
			}),
		))
	} else {
		fmt.Fprintf(s.stdLog(logStarport).out, "⚠️ Relayer error: %s\n", err)
	}

	for _, execOption := range s.plugin.ConfigCommands(chainID) {
		execOption := execOption
		steps.Add(step.New(step.NewOptions().
			Add(execOption).
			Add(s.stdSteps(logAppcli)...)...,
		))
	}

	steps.Add(step.New(step.NewOptions().
		Add(s.plugin.GentxCommand(chainID, conf)).
		Add(s.stdSteps(logAppd)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			s.app.d(),
			"collect-gentxs",
		)).
		Add(s.stdSteps(logAppd)...)...,
	))
	return steps, nil
}

func (s *Chain) createAccountSteps(ctx context.Context, name, mnemonic string, coins []string, isSilent bool) (steps step.Steps) {
	if mnemonic != "" {
		steps.Add(
			step.New(
				step.NewOptions().
					Add(s.plugin.ImportUserCommand(name, mnemonic)...)...,
			),
		)
	} else {
		generatedMnemonic := &bytes.Buffer{}
		steps.Add(
			step.New(
				step.NewOptions().
					Add(s.plugin.AddUserCommand(name)...).
					Add(
						step.PostExec(func(exitErr error) error {
							if exitErr != nil {
								return errors.Wrapf(exitErr, "cannot create %s account", name)
							}
							var user struct {
								Mnemonic string `json:"mnemonic"`
							}
							if err := json.NewDecoder(generatedMnemonic).Decode(&user); err != nil {
								return errors.Wrap(err, "cannot decode mnemonic")
							}
							if !isSilent {
								fmt.Fprintf(s.stdLog(logStarport).out, "🙂 Created an account. Password (mnemonic): %[1]v\n", user.Mnemonic)
							}
							return nil
						}),
					).
					Add(s.stdSteps(logAppcli)...).
					// Stargate pipes from stdout, Launchpad pipes from stderr.
					Add(step.Stderr(generatedMnemonic), step.Stdout(generatedMnemonic))...,
			),
		)
	}

	key := &bytes.Buffer{}
	steps.Add(
		step.New(step.NewOptions().
			Add(
				s.plugin.ShowAccountCommand(name),
				step.PostExec(func(err error) error {
					if err != nil {
						return err
					}
					coins := strings.Join(coins, ",")
					key := strings.TrimSpace(key.String())
					return cmdrunner.
						New().
						Run(ctx, step.New(step.NewOptions().
							Add(step.Exec(
								s.app.d(),
								"add-genesis-account",
								key,
								coins,
							)).
							Add(s.stdSteps(logAppd)...)...,
						))
				}),
			).
			Add(s.stdSteps(logAppcli)...).
			Add(step.Stdout(key))...,
		),
	)
	return
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
	conf, err := s.config()
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
	c, err := s.config()
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
