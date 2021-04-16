// Package scaffolder initializes Starport apps and modifies existing ones
// to add more features in a later time.
package scaffolder

import (
	"context"
	"errors"
	"os"
	"strings"

	sperrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gocmd"
	"github.com/tendermint/starport/starport/pkg/gomodule"
)

// Scaffolder is Starport app scaffolder.
type Scaffolder struct {
	// path is app's path on the filesystem.
	path string

	// options to configure scaffolding.
	options *scaffoldingOptions
}

// New initializes a new Scaffolder for app at path.
func New(path string, options ...Option) (*Scaffolder, error) {
	s := &Scaffolder{
		path:    path,
		options: newOptions(options...),
	}

	version, err := s.version()
	if err != nil && !errors.Is(err, gomodule.ErrGoModNotFound) {
		return nil, err
	}

	if err == nil && !version.Major().Is(cosmosver.Stargate) {
		return nil, sperrors.ErrOnlyStargateSupported
	}

	return s, nil
}

func (s *Scaffolder) version() (cosmosver.Version, error) {
	v, err := cosmosver.Detect(s.path)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func owner(modulePath string) string {
	return strings.Split(modulePath, "/")[1]
}

func fmtProject(path string) error {
	return cmdrunner.
		New(
			cmdrunner.DefaultStderr(os.Stderr),
			cmdrunner.DefaultWorkdir(path),
		).
		Run(context.Background(),
			step.New(
				step.Exec(
					gocmd.Name(),
					"fmt",
					"./...",
				),
			),
		)
}
