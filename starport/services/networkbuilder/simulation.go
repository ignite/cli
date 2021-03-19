package networkbuilder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/cenkalti/backoff"
	"github.com/pelletier/go-toml"
	"github.com/tendermint/starport/starport/pkg/availableport"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/chain"
)

const ValidatorSetNilErrorMessage = "validator set is nil in genesis and still empty after InitChain"

// VerifyProposals generates a genesis file from the current launch information and proposals to verify
// The function returns false if the generated genesis is invalid
func (b *Builder) VerifyProposals(ctx context.Context, chainID string, homeDir string, proposals []int, commandOut io.Writer) (bool, error) {
	// Get the simulated launch information from these proposals
	simulatedLaunchInfo, err := b.SimulatedLaunchInformation(ctx, chainID, proposals)
	if err != nil {
		return false, err
	}

	// Generate a temporary genesis from the simulated launch information
	b.ev.Send(events.New(events.StatusOngoing, "generating genesis"))
	tmpHome, err := b.GenerateTemporaryGenesis(ctx, chainID, homeDir, simulatedLaunchInfo)
	defer os.RemoveAll(tmpHome)
	if err != nil {
		return false, err
	}
	b.ev.Send(events.New(events.StatusDone, "genesis generated"))

	// set the config with random ports to test the start command
	addressAPI, err := setSimulationConfig(tmpHome)
	if err != nil {
		return false, err
	}

	// Initialize command runner
	appPath := filepath.Join(sourcePath, chainID)
	chainHandler, err := chain.New(ctx, appPath,
		chain.HomePath(tmpHome),
		chain.LogLevel(chain.LogSilent),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	)
	if err != nil {
		return false, err
	}
	commands, err := chainHandler.Commands(ctx)
	if err != nil {
		return false, err
	}
	runner := commands.
		Copy(
			chaincmdrunner.Stderr(commandOut), // This is the error of the verifying command, therefore this is the same as stdout
			chaincmdrunner.Stdout(commandOut),
		)

	// run validate-genesis command on the generated genesis
	b.ev.Send(events.New(events.StatusOngoing, "validating genesis format"))
	if runner.ValidateGenesis(ctx); err != nil {
		return false, nil
	}
	b.ev.Send(events.New(events.StatusDone, "genesis correctly formatted"))

	// verify that the chain can be started with a valid genesis
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	exit := make(chan error)

	// Go routine to check the app is listening
	go func() {
		defer cancel()
		exit <- isBlockchainListening(ctx, addressAPI)
	}()

	// Go routine to start the app
	b.ev.Send(events.New(events.StatusOngoing, "starting chain"))
	go func() {
		err := runner.Start(ctx)
		// If the error is validator set is nil, it means the genesis didn't get broken after a proposal
		// The genesis was correctly generated but we don't have the necessary proposals to have a validator set
		// after the execution of gentxs
		if strings.Contains(err.Error(), ValidatorSetNilErrorMessage) {
			err = nil
		}
		exit <- err
	}()

	if err := <-exit; err != nil {
		return false, nil
	}
	b.ev.Send(events.New(events.StatusDone, "chain can be started"))

	return true, nil
}

// setSimulationConfig sets the config for the temporary blockchain with random available port
func setSimulationConfig(appHome string) (string, error) {
	// generate random server ports and servers list.
	ports, err := availableport.Find(5)
	if err != nil {
		return "", err
	}
	genAddr := func(port int) string {
		return fmt.Sprintf("localhost:%d", port)
	}

	// updating app toml
	appPath := filepath.Join(appHome, "config/app.toml")
	config, err := toml.LoadFile(appPath)
	if err != nil {
		return "", err
	}
	config.Set("api.enable", true)
	config.Set("api.enabled-unsafe-cors", true)
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("api.address", xurl.TCP(genAddr(ports[0])))
	config.Set("grpc.address", genAddr(ports[1]))
	file, err := os.OpenFile(appPath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = config.WriteTo(file)
	if err != nil {
		return "", err
	}

	// updating config toml
	configPath := filepath.Join(appHome, "config/config.toml")
	config, err = toml.LoadFile(configPath)
	if err != nil {
		return "", err
	}
	config.Set("rpc.cors_allowed_origins", []string{"*"})
	config.Set("consensus.timeout_commit", "1s")
	config.Set("consensus.timeout_propose", "1s")
	config.Set("rpc.laddr", xurl.TCP(genAddr(ports[2])))
	config.Set("p2p.laddr", xurl.TCP(genAddr(ports[3])))
	config.Set("rpc.pprof_laddr", genAddr(ports[4]))
	file, err = os.OpenFile(configPath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = config.WriteTo(file)

	return genAddr(ports[0]), err
}

// isBlockchainListening checks if the blockchain is listening for API queries on the specified address
func isBlockchainListening(ctx context.Context, addressAPI string) error {
	checkAlive := func() error {
		ok, err := httpstatuschecker.Check(ctx, xurl.HTTP(addressAPI)+"/node_info")
		if err == nil && !ok {
			err = errors.New("app is not online")
		}
		return err
	}
	return backoff.Retry(checkAlive, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}
