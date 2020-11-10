package chain

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xexec"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/chain/conf"
	secretconf "github.com/tendermint/starport/starport/services/chain/conf/secret"
	"github.com/tendermint/starport/starport/services/chain/rly"
	"gopkg.in/yaml.v2"
)

const (
	relayerVersion = "3080dab10e37e6db9e691c343b1c0cdb6d845753"
)

// relayerInfo holds relayer info that is shared between chains to make a connection.
type relayerInfo struct {
	ChainID    string
	Mnemonic   string
	RPCAddress string
}

// RelayerInfo initializes or updates relayer setup for the chain itself and returns
// a meta info to share with other chains so they can connect.
func (s *Chain) RelayerInfo() (base64Info string, err error) {
	if err := s.checkIBCRelayerSupport(); err != nil {
		return "", err
	}
	sconf, err := secretconf.Open(s.app.Path)
	if err != nil {
		return "", err
	}
	relayerAcc, found := sconf.SelfRelayerAccount(s.app.n())
	if !found {
		if err := sconf.SetSelfRelayerAccount(s.app.n()); err != nil {
			return "", err
		}
		relayerAcc, _ = sconf.SelfRelayerAccount(s.app.n())
		if err := secretconf.Save(s.app.Path, sconf); err != nil {
			return "", err
		}
	}
	rpcAddress, err := s.rpcPublicAddress()
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

// RelayerAdd adds another chain by its relayer info to establish a connnection
// in between.
func (s *Chain) RelayerAdd(base64Info string) error {
	if err := s.checkIBCRelayerSupport(); err != nil {
		return err
	}
	data, err := base64.RawStdEncoding.DecodeString(base64Info)
	if err != nil {
		return err
	}
	var info relayerInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}
	sconf, err := secretconf.Open(s.app.Path)
	if err != nil {
		return err
	}
	sconf.UpsertRelayerAccount(conf.Account{
		Name:       info.ChainID,
		Mnemonic:   info.Mnemonic,
		RPCAddress: info.RPCAddress,
	})
	if err := secretconf.Save(s.app.Path, sconf); err != nil {
		return err
	}
	fmt.Fprint(s.stdLog(logStarport).out, "\n💫  Chain added\n")
	return nil
}

func (s *Chain) initRelayer(ctx context.Context, c conf.Config) error {
	sconf, err := secretconf.Open(s.app.Path)
	if err != nil {
		return err
	}
	if err := s.checkIBCRelayerSupport(); err != nil {
		return nil
	}
	if len(sconf.Relayer.Accounts) > 0 {
		fmt.Fprintf(s.stdLog(logStarport).out, "⌛ detected chains, linking them...\n")
	}

	// init path for the relayer.
	relayerHome, err := s.initRelayerHome()
	if err != nil {
		return err
	}
	configPath := filepath.Join(relayerHome, "config/config.yaml")

	rpcAddress, err := s.rpcPublicAddress()
	if err != nil {
		return err
	}

	selfacc, _ := sconf.SelfRelayerAccount(s.app.n())
	selfacc.RPCAddress = rpcAddress

	// prep and save relayer config.
	if _, err := s.initRelayerConfig(configPath, selfacc, sconf.Relayer.Accounts); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*15)
	defer cancel()

	// add all keys to relayer.
	for _, account := range append(
		[]conf.Account{selfacc},
		sconf.Relayer.Accounts...,
	) {
		account := account

		if err := cmdrunner.
			New().
			Run(ctx, step.New(
				step.Exec(
					"rly",
					"--home",
					relayerHome,
					"keys",
					"delete",
					account.Name,
				),
				// ignore errors related to key is not being exists anyway.
				step.PostExec(func(error) error { return nil }),
			)); err != nil {
			return err
		}

		if err := cmdrunner.
			New().
			Run(ctx, step.New(
				step.Exec(
					"rly",
					"--home",
					relayerHome,
					"keys",
					"restore",
					account.Name,
					"testkey",
					account.Mnemonic,
				),
				// check if RPC is available before adding key for this account.
				step.PreExec(func() error {
					for {
						available, err := httpstatuschecker.Check(ctx, xurl.HTTP(account.RPCAddress))
						if err == context.Canceled {
							return fmt.Errorf("Tendermint RPC not online for %q", account.Name)
						}
						if err != nil || !available {
							time.Sleep(time.Millisecond * 300)
							continue
						}
						return nil
					}
				}),
				step.Stderr(s.stdLog(logRelayer).err),
			)); err != nil {
			return err
		}
	}

	initLightClient := func(name string) error {
		return cmdrunner.
			New().
			Run(ctx, step.New(
				step.Exec(
					"rly",
					"--home",
					relayerHome,
					"light",
					"init",
					name,
					"-f",
				),
				step.Stderr(s.stdLog(logRelayer).err),
			))
	}

	// link chains.
	var wg sync.WaitGroup
	for _, account := range sconf.Relayer.Accounts {
		wg.Add(1)
		go func(account conf.Account) {
			defer wg.Done()
			err := backoff.Retry(func() error {
				if err := initLightClient(selfacc.Name); err != nil {
					return err
				}
				if err := initLightClient(account.Name); err != nil {
					return err
				}
				return cmdrunner.
					New().
					Run(ctx, step.New(
						step.Exec(
							"rly",
							"--home",
							relayerHome,
							"tx",
							"link",
							fmt.Sprintf("%s-%s", selfacc.Name, account.Name),
							"-d",
							"-o",
							"3s",
						),
						step.Stderr(s.stdLog(logRelayer).err),
					))
			}, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
			if err != nil {
				fmt.Fprintf(s.stdLog(logStarport).err, "❌ couldn't link %s <-/-> %s\n", selfacc.Name, account.Name)
			} else {
				fmt.Fprintf(s.stdLog(logStarport).out, "⛓️  linked %s <--> %s\n", selfacc.Name, account.Name)
			}
		}(account)
	}
	wg.Wait()

	return nil
}

// relayerHome initializes and returns the path to a home folder for relayer.
func (s *Chain) initRelayerHome() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	relayerHome := filepath.Join(home, s.app.nd(), "relayer")
	if os.Getenv("GITPOD_WORKSPACE_ID") != "" {
		relayerHome = filepath.Join(home, ".relayer")
	}
	if err := os.MkdirAll(filepath.Join(relayerHome, "config"), os.ModePerm); err != nil {
		return "", err
	}
	return relayerHome, nil
}

// initRelayerConfig initializes the config file of relayer and returns it.
func (s *Chain) initRelayerConfig(path string, selfacc conf.Account, accounts []conf.Account) (rly.Config, error) {
	c := rly.Config{
		Global: rly.GlobalConfig{
			Timeout:       "10s",
			LiteCacheSize: 20,
		},
		Paths: rly.Paths{},
	}

	for _, account := range append([]conf.Account{selfacc}, accounts...) {
		c.Chains = append(c.Chains, rly.NewChain(account.Name, xurl.HTTP(account.RPCAddress)))
	}

	for _, acc := range accounts {
		c.Paths[fmt.Sprintf("%s-%s", selfacc.Name, acc.Name)] = rly.NewPath(
			rly.NewPathEnd(selfacc.Name, acc.Name),
			rly.NewPathEnd(acc.Name, selfacc.Name),
		)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return rly.Config{}, err
	}
	defer file.Close()

	err = yaml.NewEncoder(file).Encode(c)
	return c, err
}

func (s *Chain) checkIBCRelayerSupport() error {
	if !s.plugin.SupportsIBC() {
		return errors.New("IBC is not available for your app.")
	}
	if !xexec.IsCommandAvailable("rly") {
		return errors.New("Relayer is not available.")
	}
	version := &bytes.Buffer{}
	return cmdrunner.
		New().
		Run(context.Background(), step.New(
			step.Exec("rly", "version"),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if !strings.Contains(version.String(), relayerVersion) {
					return fmt.Errorf("relayer is not at the required version %q", relayerVersion)
				}
				return nil
			}),
			step.Stdout(version),
		))
}
