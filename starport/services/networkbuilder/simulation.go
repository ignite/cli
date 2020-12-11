package networkbuilder

import (
	"bytes"
	"fmt"
	"context"
	"github.com/otiai10/copy"
	"github.com/pelletier/go-toml"
	"github.com/tendermint/starport/starport/pkg/availableport"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/chain"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const ErrValidatorSetNil = "validator set is nil in genesis and still empty after InitChain"

// VerifyProposals generates a genesis file from the current launch information and proposals to verify
// The function returns false if the generated genesis is invalid
func (b *Builder) VerifyProposals(ctx context.Context, chainID string, proposals []int, commandOut io.Writer) (bool, error) {
	chainInfo, err := b.ShowChain(ctx, chainID)
	if err != nil {
		return false, err
	}

	// Get the simulated launch information from these proposals
	simulatedLaunchInfo, err := b.SimulatedLaunchInformation(ctx, chainID, proposals)
	if err != nil {
		return false, err
	}

	// find out the app's name form url
	u, err := url.Parse(chainInfo.URL)
	if err != nil {
		return false, err
	}
	importPath := path.Join(u.Host, u.Path)
	path, err := gomodulepath.Parse(importPath)
	if err != nil {
		return false, err
	}
	app := chain.App{
		ChainID: chainID,
		Name:    path.Root,
		Version: cosmosver.Stargate,
	}
	chainCmd, err := chain.New(app, true, chain.LogSilent)
	if err != nil {
		return false, err
	}

	// get app dir
	homedir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	appHome := filepath.Join(homedir, app.ND())

	// create a temporary dir that holds the genesis to test
	tmpHome, err := ioutil.TempDir("", app.ND() + "*")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tmpHome)
	err = copy.Copy(appHome, tmpHome)
	if err != nil {
		return false, err
	}

	// generate the genesis to test
	if err := generateGenesis(ctx, tmpHome, chainInfo, simulatedLaunchInfo, chainCmd); err != nil {
		fmt.Fprintf(commandOut, "error generating the genesis: %s\n", err.Error())
		return false, nil
	}

	// set the config with random ports to test the start command
	if err := setSimulationConfig(tmpHome); err != nil {
		return false, err
	}

	// run validate-genesis command on the generated genesis
	err = cmdrunner.New().Run(ctx, step.New(
		step.Exec(
			app.D(),
			"validate-genesis",
			"--home",
			tmpHome,
		),
		step.Stderr(commandOut),	// This is the error of the verifying command, therefore this is the same as stdout
		step.Stdout(commandOut),
	))
	if err != nil {
		return false, nil
	}

	// verify that the chain can be started with a valid genesis
	// run validate-genesis command on the generated genesis
	errb := &bytes.Buffer{}
	err = cmdrunner.New().Run(ctx, step.New(
		step.Exec(
			app.D(),
			"start",
			"--home",
			tmpHome,
		),
		step.PostExec(func (exitErr error) error {
			// If the error is validator set is nil, it means the genesis didn't get broken after a proposal
			// The genesis was correctly generated but we don't have the necessary proposals to have a validator set
			// after the execution of gentxs
			if strings.Contains(errb.String(), ErrValidatorSetNil) {
				return nil
			}

			// We interpret any other error as if the genesis is broken
			return exitErr
		}),
		step.Stderr(io.MultiWriter(commandOut, errb)),
		step.Stdout(commandOut),
	))
	if err != nil {
		return false, nil
	}

	return true, nil
}

// setSimulationConfig sets the config for the temporary blockchain with random available port
func setSimulationConfig(appHome string) error {
	// generate random server ports and servers list.
	ports, err := availableport.Find(5)
	if err != nil {
		return err
	}
	genAddr := func(port int) string {
		return fmt.Sprintf("localhost:%d", port)
	}

	// updating app toml
	appPath := filepath.Join(appHome, "config/app.toml")
	config, err := toml.LoadFile(appPath)
	if err != nil {
		return err
	}
	config.Set("api.enable", true)
	config.Set("api.enabled-unsafe-cors", true)
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("api.address", xurl.TCP(genAddr(ports[0])))
	config.Set("grpc.address",genAddr(ports[1]))
	file, err := os.OpenFile(appPath, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	if err != nil {
		return err
	}

	// updating config toml
	configPath := filepath.Join(appHome, "config/config.toml")
	config, err = toml.LoadFile(configPath)
	if err != nil {
		return err
	}
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("consensus.timeout_commit", "1s")
	config.Set("consensus.timeout_propose", "1s")
	config.Set("rpc.laddr", xurl.TCP(genAddr(ports[2])))
	config.Set("p2p.laddr", xurl.TCP(genAddr(ports[3])))
	config.Set("rpc.pprof_laddr", genAddr(ports[4]))
	file, err = os.OpenFile(configPath, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = config.WriteTo(file)

	return err
}