package chain

import (
	"io"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/lineprefixer"
	"github.com/tendermint/starport/starport/pkg/prefixgen"
)

// prefixes holds prefix configuration for logs messages.
var prefixes = map[logType]struct {
	Name  string
	Color uint8
}{
	logStarport: {"starport", 202},
	logBuild:    {"build", 203},
	logAppd:     {"%sd", 204},
	logAppcli:   {"%scli", 205},
	logRelayer:  {"relayer", 206},
}

// logType represents the different types of logs.
type logType int

const (
	logStarport logType = iota
	logBuild
	logAppd
	logAppcli
	logRelayer
)

// std returns the cmdrunner steps to configure stdout and stderr to output logs by logType.
func (s *Chain) stdSteps(logType logType) []step.Option {
	std := s.stdLog(logType)
	return []step.Option{
		step.Stdout(std.out),
		step.Stderr(std.err),
	}
}

type std struct {
	out, err io.Writer
}

// std returns the stdout and stderr to output logs by logType.
func (s *Chain) stdLog(logType logType) std {
	prefixed := func(w io.Writer) *lineprefixer.Writer {
		var (
			prefix    = prefixes[logType]
			prefixStr string
			options   = prefixgen.Common(prefixgen.Color(prefix.Color))
			gen       = prefixgen.New(prefix.Name, options...)
		)
		if strings.Count(prefix.Name, "%s") > 0 {
			prefixStr = gen.Gen(s.app.Name)
		} else {
			prefixStr = gen.Gen()
		}
		return lineprefixer.NewWriter(w, prefixStr)
	}
	var (
		stdout io.Writer = prefixed(s.stdout)
		stderr io.Writer = prefixed(s.stderr)
	)
	if logType == logStarport && !s.verbose {
		stdout = os.Stdout
		stderr = os.Stderr
	}
	return std{
		out: stdout,
		err: stderr,
	}
}
