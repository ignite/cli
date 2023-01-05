// Package nodetime provides a single, and standalone NodeJS runtime executable that contains
// several NodeJS CLI programs bundled inside where those are reachable via subcommands.
// the CLI bundled programs are the ones that needed by Ignite CLI and more can be added as needed.
package nodetime

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"sync"

	"github.com/ignite/cli/ignite/pkg/localfs"
	"github.com/ignite/cli/ignite/pkg/nodetime/data"
)

// the list of CLIs included.
const (
	// CommandTSProto is https://github.com/stephenh/ts-proto.
	CommandTSProto CommandName = "ts-proto"

	// CommandSTA is https://github.com/acacode/swagger-typescript-api.
	CommandSTA CommandName = "sta"

	// CommandSwaggerCombine is https://www.npmjs.com/package/swagger-combine.
	CommandSwaggerCombine CommandName = "swagger-combine"

	// CommandIBCSetup is https://github.com/confio/ts-relayer/blob/main/spec/ibc-setup.md.
	CommandIBCSetup = "ibc-setup"

	// CommandIBCRelayer is https://github.com/confio/ts-relayer/blob/main/spec/ibc-relayer.md.
	CommandIBCRelayer = "ibc-relayer"

	// CommandXRelayer is a relayer wrapper for Ignite CLI made using the confio relayer.
	CommandXRelayer = "xrelayer"
)

// CommandName represents a high level command under nodetime.
type CommandName string

var (
	onceBinary sync.Once
	binary     []byte
)

// Binary returns the binary bytes of the executable.
func Binary() []byte {
	onceBinary.Do(func() {
		// untar the binary.
		gzr, err := gzip.NewReader(bytes.NewReader(data.Binary()))
		if err != nil {
			panic(err)
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)

		if _, err := tr.Next(); err != nil {
			panic(err)
		}

		if binary, err = io.ReadAll(tr); err != nil {
			panic(err)
		}
	})

	return binary
}

// Command setups the nodetime binary and returns the command needed to execute c.
func Command(c CommandName) (command []string, cleanup func(), err error) {
	cs := string(c)
	path, cleanup, err := localfs.SaveBytesTemp(Binary(), cs, 0o755)
	if err != nil {
		return nil, nil, err
	}
	command = []string{
		path,
		cs,
	}
	return command, cleanup, nil
}
