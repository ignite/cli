package chain

import (
	"io"
	"os"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/prefixgen"
)

const (
	TagFaucetAddressNotify         = "faucetAddress"
	TagBlockchainApiAddressNotify  = "blockchainAPIAddress"
	TagTendermintNodeAddressNotify = "tendermintNodeAddress"
	TagCosmosSDKVersionNotify      = "cosmosSDKVersion"
	TagBuildErrorNotify            = "buildError"
	TagValidationErrorNotify       = "buildError"
	TagWaitForChangesNotify        = "waitForChanges"
	TagRebuildTriggeredNotify      = "rebuildTriggered"
	TagUnrecognisedErrorNotify     = "unrecognisedError"
)

// prefixes holds prefix configuration for logs messages.
var prefixes = map[logType]struct {
	Name  string
	Color uint8
}{
	logStarport: {"starport", 202},
	logBuild:    {"build", 203},
	logAppd:     {"%s daemon", 204},
}

// logType represents the different types of logs.
type logType int

const (
	logStarport logType = iota
	logBuild
	logAppd
)

type std struct {
	out, err io.Writer
}

// std returns the stdout and stderr to output logs by logType.
func (c *Chain) stdLog() std {

	stdout := os.Stdout
	stderr := os.Stderr
	return std{
		out: stdout,
		err: stderr,
	}
}

func (c *Chain) genPrefix(logType logType) string {
	prefix := prefixes[logType]

	return prefixgen.New(prefix.Name, prefixgen.Common(prefixgen.Color(prefix.Color))...).
		Gen(c.app.Name)
}
