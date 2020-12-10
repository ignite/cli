package networkbuilder

import (
	"context"
	"github.com/otiai10/copy"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/services/chain"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

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

	return true, nil
}