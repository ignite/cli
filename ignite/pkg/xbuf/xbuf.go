package xbuf

import (
	"bytes"
	"context"
	"path/filepath"

	"github.com/bufbuild/buf/private/buf/bufcli"
	"github.com/bufbuild/buf/private/buf/buffetch"
	"github.com/bufbuild/buf/private/buf/bufgen"
	"github.com/bufbuild/buf/private/bufpkg/bufanalysis"
	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/bufpkg/bufwasm"
	"github.com/bufbuild/buf/private/pkg/app/applog"
	"github.com/bufbuild/buf/private/pkg/app/appverbose"
	"github.com/bufbuild/buf/private/pkg/command"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
)

const (
	gogoTemplate    = "buf.gen.gogo.yaml"
	swaggerTemplate = "buf.gen.swagger.yaml"
	errorFormat     = "json"
	logLevel        = "error"
	appName         = "ignite"
)

func Generate(
	ctx context.Context,
	template,
	input,
	output string,
) error {
	var (
		stdOut            = new(bytes.Buffer)
		storageOSProvider = bufcli.NewStorageosProvider(false)
		runner            = command.NewRunner()
	)
	logger, err := applog.NewLogger(stdOut, logLevel, errorFormat)
	if err != nil {
		return err
	}

	verbosePrinter := appverbose.NewVerbosePrinter(stdOut, appName, false)
	ctn, err := newContainer(appName, logger, verbosePrinter)
	if err != nil {
		return err
	}

	ref, err := buffetch.NewRefParser(logger).GetRef(ctx, input)
	if err != nil {
		return err
	}

	readWriteBucket, err := storageOSProvider.NewReadWriteBucket(
		".",
		storageos.ReadWriteBucketWithSymlinksIfSupported(),
	)
	if err != nil {
		return err
	}

	genConfig, err := bufgen.ReadConfig(
		ctx,
		logger,
		bufgen.NewProvider(logger),
		readWriteBucket,
		bufgen.ReadConfigWithOverride(template),
	)
	if err != nil {
		return err
	}

	clientConfig, err := bufcli.NewConnectClientConfig(ctn)
	if err != nil {
		return err
	}

	imageConfigReader, err := bufcli.NewWireImageConfigReader(
		ctn,
		storageOSProvider,
		runner,
		clientConfig,
	)
	if err != nil {
		return err
	}

	imageConfigs, fileAnnotations, err := imageConfigReader.GetImageConfigs(
		ctx,
		ctn,
		ref,
		"",
		nil,
		nil,
		false,
		false,
	)
	if err != nil {
		return err
	}

	if len(fileAnnotations) > 0 {
		err := bufanalysis.PrintFileAnnotations(ctn.Stderr(), fileAnnotations, errorFormat)
		if err != nil {
			return err
		}
		return bufcli.ErrFileAnnotation
	}
	images := make([]bufimage.Image, 0, len(imageConfigs))
	for _, imageConfig := range imageConfigs {
		images = append(images, imageConfig.Image())
	}
	image, err := bufimage.MergeImages(images...)
	if err != nil {
		return err
	}
	generateOptions := []bufgen.GenerateOption{
		bufgen.GenerateWithBaseOutDirPath(output),
	}

	wasmEnabled, err := bufcli.IsAlphaWASMEnabled(ctn)
	if err != nil {
		return err
	}
	if wasmEnabled {
		generateOptions = append(
			generateOptions,
			bufgen.GenerateWithWASMEnabled(),
		)
	}

	wasmPluginExecutor, err := bufwasm.NewPluginExecutor(
		filepath.Join(ctn.CacheDirPath(), bufcli.WASMCompilationCacheDir))
	if err != nil {
		return err
	}
	return bufgen.NewGenerator(
		logger,
		storageOSProvider,
		runner,
		wasmPluginExecutor,
		clientConfig,
	).Generate(
		ctx,
		ctn,
		genConfig,
		image,
		generateOptions...,
	)
}
