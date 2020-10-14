// TODO change pkg name to chain.
package starportserve

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/fswatcher"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/pkg/xurl"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
	"github.com/tendermint/starport/starport/services/serve/rly"
	starportsecretconf "github.com/tendermint/starport/starport/services/serve/secretconf"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

var (
	appBackendWatchPaths = append([]string{
		"app",
		"cmd",
		"x",
		starportsecretconf.SecretFile,
	}, starportconf.FileNames...)

	vuePath = "vue"

	errorColor = color.Red.Render
	infoColor  = color.Yellow.Render
)

type version struct {
	tag  string
	hash string
}

// TODO change name to Chain.
type Serve struct {
	app            App
	plugin         Plugin
	version        version
	verbose        bool
	serveCancel    context.CancelFunc
	serveRefresher chan struct{}
	stdout, stderr io.Writer
}

func New(app App, verbose bool) (*Serve, error) {
	s := &Serve{
		app:            app,
		verbose:        verbose,
		serveRefresher: make(chan struct{}, 1),
		stdout:         ioutil.Discard,
		stderr:         ioutil.Discard,
	}

	if verbose {
		s.stdout = os.Stdout
		s.stderr = os.Stderr
	}

	var err error

	s.version, err = s.appVersion()
	if err != nil && err != git.ErrRepositoryNotExists {
		return nil, err
	}

	s.plugin, err = s.pickPlugin()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Build builds an app.
func (s *Serve) Build(ctx context.Context) error {
	if err := s.setup(ctx); err != nil {
		return err
	}
	conf, err := s.config()
	if err != nil {
		return &CannotBuildAppError{err}
	}
	steps, binaries := s.buildSteps(ctx, conf)
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, steps...); err != nil {
		return err
	}
	fmt.Fprintf(s.stdLog(logStarport).out, "🗃  Installed. Use with: %s\n", infoColor(strings.Join(binaries, ", ")))
	return nil
}

type relayerInfo struct {
	ChainID    string
	Mnemonic   string
	RPCAddress string
}

// RelayerInfo initializes or updates relayer setup for the chain itself and returns
// a meta info to share with other chains so they can connect.
// TODO only stargate
func (s *Serve) RelayerInfo() (base64Info string, err error) {
	sconf, err := starportsecretconf.Open(s.app.Path)
	if err != nil {
		return "", err
	}
	relayerAcc, found := sconf.SelfRelayerAccount(s.app.n())
	if !found {
		if err := sconf.SetSelfRelayerAccount(s.app.n()); err != nil {
			return "", err
		}
		relayerAcc, _ = sconf.SelfRelayerAccount(s.app.n())
		if err := starportsecretconf.Save(s.app.Path, sconf); err != nil {
			return "", err
		}
	}
	rpcAddress, err := s.rpcAddress()
	if err != nil {
		return "", err
	}
	info := relayerInfo{
		ChainID:    s.app.n(),
		Mnemonic:   relayerAcc.Mnemonic,
		RPCAddress: rpcAddress,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

func (s *Serve) RelayerAdd(base64Info string) error {
	data, err := base64.RawStdEncoding.DecodeString(base64Info)
	if err != nil {
		return err
	}
	var info relayerInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}
	sconf, err := starportsecretconf.Open(s.app.Path)
	if err != nil {
		return err
	}
	sconf.UpsertRelayerAccount(starportsecretconf.RelayerAccount{
		ID:         info.ChainID,
		Mnemonic:   info.Mnemonic,
		RPCAddress: info.RPCAddress,
	})
	if err := starportsecretconf.Save(s.app.Path, sconf); err != nil {
		return err
	}
	fmt.Fprint(s.stdLog(logStarport).out, "\n💫  Chain added\n")
	return nil
}

func (s *Serve) rpcAddress() (string, error) {
	rpcAddress := os.Getenv("RPC_ADDRESS")
	if rpcAddress == "" {
		conf, err := s.config()
		if err != nil {
			return "", err
		}
		rpcAddress = conf.Servers.RPCAddr
	}
	return rpcAddress, nil
}

// Serve serves an app.
func (s *Serve) Serve(ctx context.Context) error {
	if err := s.setup(ctx); err != nil {
		return err
	}
	_, err := s.RelayerInfo()
	if err != nil {
		return err
	}
	xos.RemoveAllUnderHome(".relayer")

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

					var validationErr *starportconf.ValidationError
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

func (s *Serve) setup(ctx context.Context) error {
	fmt.Fprintf(s.stdLog(logStarport).out, "Cosmos' version is: %s\n", infoColor(s.plugin.Name()))

	if err := s.checkSystem(); err != nil {
		return err
	}
	if err := s.plugin.Migrate(ctx); err != nil {
		return err
	}
	return nil
}

// checkSystem checks if developer's work environment comply must to have
// dependencies and pre-conditions.
func (s *Serve) checkSystem() error {
	// check if Go has installed.
	if !xexec.IsCommandAvailable("go") {
		return errors.New("Please, check that Go language is installed correctly in $PATH. See https://golang.org/doc/install")
	}
	if s.plugin.Version() == cosmosver.Stargate && !xexec.IsCommandAvailable("rly") {
		return errors.New("Please, check that Relayer is installed.")
	}
	// check if Go's bin added to System's path.
	gobinpath := path.Join(build.Default.GOPATH, "bin")
	if err := xos.IsInPath(gobinpath); err != nil {
		return errors.New("$(go env GOPATH)/bin must be added to your $PATH. See https://golang.org/doc/gopath_code.html#GOPATH")
	}
	return nil
}

func (s *Serve) refreshServe() {
	if s.serveCancel != nil {
		s.serveCancel()
	}
	s.serveRefresher <- struct{}{}
}

func (s *Serve) watchAppBackend(ctx context.Context) error {
	return fswatcher.Watch(
		ctx,
		appBackendWatchPaths,
		fswatcher.Workdir(s.app.Path),
		fswatcher.OnChange(s.refreshServe),
		fswatcher.IgnoreHidden(),
	)
}

func (s *Serve) cmdOptions() []cmdrunner.Option {
	return []cmdrunner.Option{
		cmdrunner.DefaultWorkdir(s.app.Path),
	}
}

func (s *Serve) serve(ctx context.Context) error {
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
	if err := cmdrunner.
		New(s.cmdOptions()...).
		Run(ctx, s.initSteps(ctx, conf)...); err != nil {
		return err
	}

	return cmdrunner.
		New(append(s.cmdOptions(), cmdrunner.RunParallel())...).
		Run(ctx, s.serverSteps(ctx, conf)...)
}

func (s *Serve) initSteps(ctx context.Context, conf starportconf.Config) (
	steps step.Steps) {
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
				"--chain-id", s.app.n(),
			),
			step.PostExec(func(err error) error {
				// overwrite Genesis with user configs.
				if err != nil {
					return err
				}
				if conf.Genesis == nil {
					return nil
				}
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				path := filepath.Join(home, s.plugin.GenesisPath())
				file, err := os.OpenFile(path, os.O_RDWR, 644)
				if err != nil {
					return err
				}
				defer file.Close()
				var genesis map[string]interface{}
				if err := json.NewDecoder(file).Decode(&genesis); err != nil {
					return err
				}
				if err := mergo.Merge(&genesis, conf.Genesis, mergo.WithOverride); err != nil {
					return err
				}
				if err := file.Truncate(0); err != nil {
					return err
				}
				if _, err := file.Seek(0, 0); err != nil {
					return err
				}
				return json.NewEncoder(file).Encode(&genesis)
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
		account := account
		steps.Add(s.createAccountSteps(ctx, account.Name, account.Coins, false)...)
	}
	info, err := s.RelayerInfo()
	if err != nil {
		panic(err)
	}
	steps.Add(step.New(
		step.PreExec(func() error {
			fmt.Fprintf(s.stdLog(logStarport).out, "✨ Relayer info: %s\n", info)
			return nil
		}),
	))

	if s.plugin.Version() == cosmosver.Stargate {
		sconf, err := starportsecretconf.Open(s.app.Path)
		if err != nil {
			panic(err)
		}
		var (
			key = &bytes.Buffer{}
		)
		for _, acc := range sconf.Accounts {
			acc := acc
			steps.Add(step.New(
				step.Exec(
					s.app.d(),
					"keys",
					"add",
					acc.Name,
					"--recover",
					"--keyring-backend", "test",
				),
				step.Write([]byte(acc.Mnemonic+"\n")),
			))
			steps.Add(step.New(step.NewOptions().
				Add(
					s.plugin.ShowAccountCommand(acc.Name),
					step.PostExec(func(err error) error {
						if err != nil {
							return err
						}
						coins := strings.Join(acc.Coins, ",")
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
				Add(step.Stdout(key))...,
			))
		}
	}

	for _, execOption := range s.plugin.ConfigCommands() {
		execOption := execOption
		steps.Add(step.New(step.NewOptions().
			Add(execOption).
			Add(s.stdSteps(logAppcli)...)...,
		))
	}

	steps.Add(step.New(step.NewOptions().
		Add(s.plugin.GentxCommand(conf)).
		Add(s.stdSteps(logAppd)...)...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(step.Exec(
			s.app.d(),
			"collect-gentxs",
		)).
		Add(s.stdSteps(logAppd)...)...,
	))
	return
}

func (s *Serve) relayerSteps(ctx context.Context, sconf *starportsecretconf.Config) (steps step.Steps) {
	// prep relayer config
	rlyConf := rly.Config{
		Global: rly.GlobalConfig{
			Timeout:       "10s",
			LiteCacheSize: 20,
		},
		Paths: rly.Paths{},
	}
	rpcAddress, err := s.rpcAddress()
	if err != nil {
		// TODO no panic
		panic(err)
	}

	selfAcc, _ := sconf.SelfRelayerAccount(s.app.n())
	rlyConf.Chains = append(rlyConf.Chains, rly.NewChain(selfAcc.Name, xurl.HTTP(rpcAddress)))
	for _, acc := range sconf.Relayer.Accounts {
		rlyConf.Chains = append(rlyConf.Chains, rly.NewChain(acc.ID, xurl.HTTP(acc.RPCAddress)))
		rlyConf.Paths[fmt.Sprintf("%s-%s", selfAcc.Name, acc.ID)] = rly.NewPath(
			rly.NewPathEnd(selfAcc.Name, acc.ID),
			rly.NewPathEnd(acc.ID, selfAcc.Name),
		)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		// TODO no panic
		panic(err)
	}
	relayerHome := filepath.Join(home, s.app.nd(), "relayer")
	if err := os.MkdirAll(filepath.Join(relayerHome, "config"), os.ModePerm); err != nil {
		// TODO no panic
		panic(err)
	}
	if os.Getenv("GITPOD_WORKSPACE_ID") != "" {
		relayerHome = filepath.Join(home, ".relayer")
	}

	configPath := filepath.Join(relayerHome, "config/config.yaml")
	os.MkdirAll(filepath.Join(relayerHome, "config"), os.ModePerm)

	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		// TODO no panic
		panic(err)
	}
	defer file.Close()
	if err := yaml.NewEncoder(file).Encode(rlyConf); err != nil {
		// TODO no panic
		panic(err)
	}

	if len(sconf.Relayer.Accounts) == 0 {
		return steps
	}
	conf, err := s.config()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(s.stdLog(logStarport).out, "⌛ detected chains, linking them...\n")
	time.Sleep(time.Second * 10)
	steps.Add(step.New(
		step.Exec("rly", "--home", relayerHome, "keys", "delete", s.app.n()),
		step.PostExec(func(error) error {
			return nil
		}),
	))
	steps.Add(step.New(
		step.Exec("rly", "--home", relayerHome, "keys", "restore", s.app.n(), "testkey", selfAcc.Mnemonic),
		step.PreExec(func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			for {
				available, err := httpstatuschecker.Check(ctx, xurl.HTTP(conf.Servers.RPCAddr))
				if err == context.Canceled {
					return errors.New("tendermint cannot get online")
				}
				if err != nil || !available {
					time.Sleep(time.Millisecond * 300)
					continue
				}
				return nil
			}
		}),
		step.Stderr(os.Stderr),
	))
	for _, acc := range sconf.Relayer.Accounts {
		acc := acc
		steps.Add(step.New(
			step.Exec("rly", "--home", relayerHome, "keys", "delete", acc.ID),
			step.PostExec(func(error) error {
				return nil
			}),
		))
		steps.Add(step.New(
			step.PreExec(func() error {
				ctx, cancel := context.WithTimeout(ctx, time.Second*20)
				defer cancel()
				for {
					available, err := httpstatuschecker.Check(ctx, xurl.HTTP(acc.RPCAddress))
					if err == context.Canceled {
						fmt.Fprintf(s.stdLog(logStarport).err, "\n❌ relayer: cannot link with %s\n", acc.ID)
						return nil
					}
					if err != nil || !available {
						time.Sleep(time.Millisecond * 300)
						continue
					}
					return nil
				}
			}),
			step.Exec("rly", "--home", relayerHome, "keys", "restore", acc.ID, "testkey", acc.Mnemonic),
			step.Stderr(os.Stderr),
		))
	}
	steps.Add(step.New(
		step.Exec("rly", "--home", relayerHome, "light", "init", s.app.n(), "-f"),
		step.Stderr(os.Stderr),
	))
	for _, acc := range sconf.Relayer.Accounts {
		acc := acc
		steps.Add(step.New(
			step.Exec("rly", "--home", relayerHome, "light", "init", acc.ID, "-f"),
			step.Stderr(os.Stderr),
		))
	}
	for name := range rlyConf.Paths {
		name := name
		steps.Add(step.New(
			step.Exec("rly", "--home", relayerHome, "tx", "link", name, "-d", "-o", "3s"),
			step.PostExec(func(execErr error) error {
				c := strings.Split(name, "-")
				if execErr != nil {
					fmt.Fprintf(s.stdLog(logStarport).err, "❌ couldn't link %s <-/-> %s\n", c[0], c[1])
				} else {
					fmt.Fprintf(s.stdLog(logStarport).out, "⛓️  linked %s <--> %s\n", c[0], c[1])
				}
				return nil
			}),
			step.Stderr(os.Stderr),
		))
	}
	return steps
}

func (s *Serve) createAccountSteps(ctx context.Context, name string, coins []string, isSilent bool) []*step.Step {
	var (
		key      = &bytes.Buffer{}
		mnemonic = &bytes.Buffer{}
	)
	return []*step.Step{
		step.New(step.NewOptions().
			Add(
				s.plugin.AddUserCommand(name),
				step.PostExec(func(exitErr error) error {
					if exitErr != nil {
						return errors.Wrapf(exitErr, "cannot create %s account", name)
					}
					var user struct {
						Mnemonic string `json:"mnemonic"`
					}
					if err := json.NewDecoder(mnemonic).Decode(&user); err != nil {
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
			Add(step.Stderr(mnemonic), step.Stdout(mnemonic))...,
		),
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
	}
}

func (s *Serve) buildSteps(ctx context.Context, conf starportconf.Config) (
	steps step.Steps, binaries []string) {
	ldflags := fmt.Sprintf(`'-X github.com/cosmos/cosmos-sdk/version.Name=NewApp 
	-X github.com/cosmos/cosmos-sdk/version.ServerName=%sd 
	-X github.com/cosmos/cosmos-sdk/version.ClientName=%scli 
	-X github.com/cosmos/cosmos-sdk/version.Version=%s 
	-X github.com/cosmos/cosmos-sdk/version.Commit=%s'`, s.app.Name, s.app.Name, s.version.tag, s.version.hash)
	var (
		buildErr = &bytes.Buffer{}
	)
	captureBuildErr := func(err error) error {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &CannotBuildAppError{errors.New(buildErr.String())}
		}
		return err
	}
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				"go",
				"mod",
				"tidy",
			),
			step.PreExec(func() error {
				fmt.Fprintln(s.stdLog(logStarport).out, "\n📦 Installing dependencies...")
				return nil
			}),
			step.PostExec(captureBuildErr),
		).
		Add(s.stdSteps(logStarport)...).
		Add(step.Stderr(buildErr))...,
	))
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				"go",
				"mod",
				"verify",
			),
			step.PostExec(captureBuildErr),
		).
		Add(s.stdSteps(logBuild)...).
		Add(step.Stderr(buildErr))...,
	))

	// install the app.
	steps.Add(step.New(
		step.PreExec(func() error {
			fmt.Fprintln(s.stdLog(logStarport).out, "🛠️  Building the app...")
			return nil
		}),
	))
	installOptions, binaries := s.plugin.InstallCommands(ldflags)
	for _, execOption := range installOptions {
		execOption := execOption
		steps.Add(step.New(step.NewOptions().
			Add(
				execOption,
				step.PostExec(captureBuildErr),
			).
			Add(s.stdSteps(logStarport)...).
			Add(step.Stderr(buildErr))...,
		))
	}
	return steps, binaries
}

func (s *Serve) serverSteps(ctx context.Context, conf starportconf.Config) (steps step.Steps) {
	var wg sync.WaitGroup
	wg.Add(len(s.plugin.StartCommands(conf)))
	go func() {
		wg.Wait()
		fmt.Fprintf(s.stdLog(logStarport).out, "🌍 Running a Cosmos '%[1]v' app with Tendermint at %s.\n", s.app.Name, xurl.HTTP(conf.Servers.RPCAddr))
		fmt.Fprintf(s.stdLog(logStarport).out, "🌍 Running a server at %s (LCD)\n", xurl.HTTP(conf.Servers.APIAddr))
		fmt.Fprintf(s.stdLog(logStarport).out, "\n🚀 Get started: %s\n\n", xurl.HTTP(conf.Servers.DevUIAddr))
		if s.plugin.Version() == cosmosver.Stargate {
			sconf, err := starportsecretconf.Open(s.app.Path)
			if err != nil {
				if err != context.Canceled {
					panic(err)
				}
			}
			relayerSteps := s.relayerSteps(ctx, sconf)
			if len(relayerSteps) > 0 {
				if err := cmdrunner.
					New(s.cmdOptions()...).
					Run(ctx, relayerSteps...); err != nil {
					if err != context.Canceled {
						panic(err)
					}
				}
			}
		}

	}()
	for _, execOption := range s.plugin.StartCommands(conf) {
		execOption := execOption
		steps.Add(step.New(step.NewOptions().
			Add(execOption...).
			Add(step.InExec(func() error {
				wg.Done()
				return nil
			})).
			Add(s.stdSteps(logAppd)...)...,
		))
	}
	return
}

func (s *Serve) watchAppFrontend(ctx context.Context) error {
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

func (s *Serve) runDevServer(ctx context.Context) error {
	c, err := s.config()
	if err != nil {
		return err
	}
	conf := Config{
		SdkVersion:      s.plugin.Name(),
		EngineAddr:      xurl.HTTP(c.Servers.RPCAddr),
		AppBackendAddr:  xurl.HTTP(c.Servers.APIAddr),
		AppFrontendAddr: xurl.HTTP(c.Servers.FrontendAddr),
	} // TODO get vals from const
	handler, err := newDevHandler(s.app, conf)
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

func (s *Serve) appVersion() (v version, err error) {
	repo, err := git.PlainOpen(s.app.Path)
	if err != nil {
		return version{}, err
	}
	iter, err := repo.Tags()
	if err != nil {
		return version{}, err
	}
	ref, err := iter.Next()
	if err != nil {
		return version{}, nil
	}
	v.tag = strings.TrimPrefix(ref.Name().Short(), "v")
	v.hash = ref.Hash().String()
	return v, nil
}

func (s *Serve) config() (starportconf.Config, error) {
	var paths []string
	for _, name := range starportconf.FileNames {
		paths = append(paths, filepath.Join(s.app.Path, name))
	}
	confFile, err := xos.OpenFirst(paths...)
	if err != nil {
		return starportconf.Config{}, errors.Wrap(err, "config file cannot be found")
	}
	defer confFile.Close()
	return starportconf.Parse(confFile)
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
