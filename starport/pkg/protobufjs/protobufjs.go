package protobufjs

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
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
	// CachePath is the path where protobufjs binary is cached in the local fs.
	CachePath = "/tmp/protobufjs"
)

var cacheOnce sync.Once

// Generate generates static protobuf.js types for given proto where includePaths holds dependency protos.
// TODO add ts generation. protobufjs supports this but by executing jsdoc command with node dynamically,
// it doesn't work with bundled node pkg. things needs to be reconstructed.
func Generate(ctx context.Context, outDir, outName, protoPath string, includePaths []string) error {
	var err error

	// caches the protobufjs-cli into CachePath if it isn't there already.
	cacheOnce.Do(func() { err = cacheBinary() })

	if err != nil {
		return err
	}

	runcmd := func(command []string) error {
		errb := &bytes.Buffer{}

		err = cmdrunner.
			New(
				cmdrunner.DefaultStderr(errb)).
			Run(ctx,
				step.New(step.Exec(command[0], command[1:]...)))

		return errors.Wrap(err, errb.String())
	}

	var (
		jsOutPath = filepath.Join(outDir, outName+".js")
	)

	// construct js gen command for the actual code generation.
	command := []string{
		CachePath,
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
		command = append(
			command,
			"-p",
			includePath,
		)
	}

	// add target proto path to that.
	command = append(command, protoanalysis.GlobPattern(protoPath))

	// run the js command.
	return runcmd(command)
}

func cacheBinary() (err error) {
	// make sure the parent dir of CachePath exists.
	if err = os.MkdirAll(filepath.Dir(CachePath), os.ModePerm); err != nil {
		return err
	}

	// save saves the cli at CachePath.
	save := func() error {
		cachedFile, err := os.OpenFile(CachePath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer cachedFile.Close()

		_, err = io.Copy(cachedFile, bytes.NewReader(Bytes()))
		return err
	}

	// cache the cli if it doesn't exists.
	cachedFile, err := os.Open(CachePath)
	if os.IsNotExist(err) {
		return save()
	}
	if err != nil {
		return err
	}

	// compare hashes of the existent cli and original one.
	// if they're not the same, cache the original again.
	var (
		hasherOriginal = sha256.New()
		hasherCached   = sha256.New()
	)

	hasherOriginal.Write(Bytes())

	if _, err = io.Copy(hasherCached, cachedFile); err != nil {
		return err
	}

	hashCached := fmt.Sprintf("%x", hasherCached.Sum(nil))
	hashOriginal := fmt.Sprintf("%x", hasherOriginal.Sum(nil))

	if hashOriginal != hashCached {
		return save()
	}

	return nil
}

// Bytes returns the executable binary bytes of protobufjs.
func Bytes() []byte {
	names, err := AssetDir("")
	if err != nil {
		panic(err)
	}
	return MustAsset(names[0])
}
