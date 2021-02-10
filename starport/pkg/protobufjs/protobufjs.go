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

// Generate generates static protobuf.js types for given proto where includePaths holds dependency protos.
func Generate(ctx context.Context, outPath, protoPath string, includePaths []string) error {
	if err := cacheBinary(); err != nil {
		return err
	}

	command := []string{
		CachePath,
		"-t",
		"static-module",
		"-w",
		"commonjs",
		"-o",
		outPath,
	}

	for _, includePath := range includePaths {
		command = append(
			command,
			"-p",
			includePath,
		)
	}

	command = append(command, protoanalysis.GlobPattern(protoPath))

	errb := &bytes.Buffer{}

	err := cmdrunner.
		New(
			cmdrunner.DefaultStderr(errb)).
		Run(ctx,
			step.New(step.Exec(command[0], command[1:]...)))

	if err != nil {
		return errors.Wrap(err, errb.String())
	}

	return nil
}

var cacheOnce sync.Once

func cacheBinary() (err error) {
	cacheOnce.Do(func() {
		if err = os.MkdirAll(filepath.Dir(CachePath), os.ModePerm); err != nil {
			return
		}

		var cached *os.File

		cache := func() {
			cached, err = os.OpenFile(CachePath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return
			}
			defer cached.Close()

			_, err = io.Copy(cached, bytes.NewReader(Bytes()))
		}

		cached, err = os.Open(CachePath)
		if os.IsNotExist(err) {
			cache()
			return
		}
		if err != nil {
			return
		}

		var (
			hasheroriginal = sha256.New()
			hashercached   = sha256.New()
		)

		hasheroriginal.Write(Bytes())

		if _, err = io.Copy(hashercached, cached); err != nil {
			return
		}

		hashcached := fmt.Sprintf("%x", hashercached.Sum(nil))
		hashoriginal := fmt.Sprintf("%x", hashercached.Sum(nil))

		if hashoriginal != hashcached {
			cache()
		}
	})

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
