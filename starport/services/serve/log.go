package starportserve

import (
	"io"
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
	logAppd:     {"%sd", 203},
	logAppcli:   {"%scli", 205},
}

// logType represents the different types of logs.
type logType int

const (
	logStarport logType = iota
	logAppd
	logAppcli
)

// defaultStd configures default stdout, stderr for cmdrunner steps by logTypes.
func (s *starportServe) defaultStd(logType logType) []step.Option {
	std := func(w io.Writer) *lineprefixer.Writer {
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
	return []step.Option{
		step.Stdout(std(s.stdout)),
		step.Stderr(std(s.stderr)),
	}
}
