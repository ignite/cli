package localspn

import (
	"context"
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/httpstatuschecker"
	"github.com/tendermint/starport/starport/pkg/xurl"
)

const (
	repoURL        = "https://github.com/tendermint/spn"
	defaultBranch  = "master"
	defaultAPIPort = "1317"
)

type spnOptions struct {
	ref plumbing.ReferenceName
}

// SPNOption is an option to configure local SPN setup
type SPNOption func(*spnOptions)

// WithBranch configures the branch to fetch SPN
func WithBranch(branch string) SPNOption {
	return func(o *spnOptions) {
		o.ref = plumbing.NewBranchReferenceName(branch)
	}
}

// newSPNOptions initializes spnOptions
func newSPNOptions(options ...SPNOption) *spnOptions {
	spnOptions := &spnOptions{
		ref: plumbing.NewBranchReferenceName(defaultBranch),
	}

	for _, option := range options {
		option(spnOptions)
	}

	return spnOptions
}

// SetupSPN installs a starts a local SPN
func SetupSPN(ctx context.Context, options ...SPNOption) (cleanup func(), err error) {
	// Create temporary directory for SPN repo
	spnPath, err := os.MkdirTemp("", "spn")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(spnPath)
		}
	}()

	// Clone SPN
	spnOptions := newSPNOptions(options...)
	gitoptions := &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: spnOptions.ref,
		SingleBranch:  true,
	}
	_, err = git.PlainCloneContext(ctx, spnPath, false, gitoptions)
	if err != nil {
		return nil, err
	}

	// Create a temporary directory for SPN home
	spnHome, err := os.MkdirTemp("", "spn-home")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(spnHome)
		}
	}()

	// Start SPN
	if err := startSPN(ctx, spnPath, spnHome); err != nil {
		return nil, err
	}

	return func() {
		os.RemoveAll(spnPath)
		os.RemoveAll(spnHome)
	}, nil
}

// startSPN starts a local instance of SPN
func startSPN(ctx context.Context, spnPath, spnHome string) error {
	runStep := step.NewSteps(step.New(
		step.Exec("starport", "serve", "--home", spnHome),
		step.Workdir(spnPath),
	))
	cmdrunner.New().Run(ctx, runStep...)

	checkReadiness := func() error {
		ok, err := httpstatuschecker.Check(ctx, xurl.HTTP(defaultAPIPort)+"/node_info")
		if err == nil && !ok {
			err = errors.New("app is not online")
		}
		return err
	}

	return nil
}
