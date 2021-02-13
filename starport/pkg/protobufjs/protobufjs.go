package protobufjs

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

const (
	// BinaryPath is the path where protobufjs binary is placed in the local fs.
	BinaryPath = "/tmp/protobufjs"
)

var placeOnce sync.Once

// Generate generates static protobuf.js types for given proto where includePaths holds dependency protos.
// TODO add ts generation. protobufjs supports this but by executing jsdoc command with node dynamically,
// it doesn't work with bundled node pkg. things needs to be reconstructed.
func Generate(ctx context.Context, outDir, outName, protoPath string, includePaths []string) error {
	var err error

	// places the protobufjs-cli into BinaryPath.
	placeOnce.Do(func() { err = placeBinary() })

	if err != nil {
		return err
	}

	runcmd := func(command []string) error {
		errb := &bytes.Buffer{}

		err = cmdrunner.
			New(cmdrunner.DefaultStderr(errb)).
			Run(ctx, step.New(step.Exec(command[0], command[1:]...)))

		return errors.Wrap(err, errb.String())
	}

	var (
		jsOutPath = filepath.Join(outDir, outName+".js")
	)

	// construct js gen command for the actual code generation.
	command := []string{
		BinaryPath,
		"js",
		"-t",
		"static-module",
		"-w",
		"es6",
		"-o",
		jsOutPath,
	}

	// add proto dependency paths to that.
	for _, includePath := range includePaths {
		if _, err := os.Stat(includePath); os.IsNotExist(err) {
			continue
		}

		command = append(
			command,
			"-p",
			includePath,
		)
	}

	// add target proto path to that.
	command = append(command, protoanalysis.GlobPattern(protoPath))

	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	// run the js command.
	return runcmd(command)
}

func placeBinary() error {
	if err := os.MkdirAll(filepath.Dir(BinaryPath), os.ModePerm); err != nil {
		return err
	}

	gzr, err := gzip.NewReader(bytes.NewReader(Bytes()))
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	if _, err := tr.Next(); err != nil {
		return err
	}
	f, err := os.OpenFile(BinaryPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, tr)
	return err
}

// Bytes returns the executable binary bytes of protobufjs.
func Bytes() []byte {
	names, err := AssetDir("")
	if err != nil {
		panic(err)
	}
	return MustAsset(names[0])
}
