package localspn

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/cenkalti/backoff"
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
	defaultAPIHost = "127.0.0.1:1317"
	defaultTimeout = 30 * time.Second
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
	gitOptions := &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: spnOptions.ref,
		SingleBranch:  true,
	}
	_, err = git.PlainCloneContext(ctx, spnPath, false, gitOptions)
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
	serveStep := step.NewSteps(step.New(
		step.Exec("starport", "serve", "--home", spnHome),
		step.Workdir(spnPath),
	))
	spnChan := make(chan error, 1)

	// SPN must be served before the timeout
	go func () {
		// SPN execution routine
		err := cmdrunner.New().Run(ctx, serveStep...)
		if err == nil {
			err = errors.New("spn server stopped")
		}
		spnChan <- err
	}()
	go func () {
		// Timeout routine
		time.Sleep(defaultTimeout)
		spnChan <- errors.New("spn server failed to start")
	}()
	go func () {
		// Check SPN readiness
		spnChan <- spnServed(ctx)
	}()

	// Wait for the first routine to complete
	spnError := <- spnChan
	if spnError != nil {
		return spnError
	}

	return nil
}

// spnServed returns once spn server is served
func spnServed(ctx context.Context) error {
	checkReadiness := func() error {
		ok, err := httpstatuschecker.Check(ctx, xurl.HTTP(defaultAPIHost)+"/node_info")
		if err == nil && !ok {
			err = errors.New("spn is not online")
		}
		return err
	}

	return backoff.Retry(checkReadiness, backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx))
}
