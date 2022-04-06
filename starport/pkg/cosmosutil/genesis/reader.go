package genesis

import (
	"context"
	"github.com/karlseguin/jsonwriter"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/tarball"
	"io"
	"net/http"
	"os"
)

type (
	// GenReader represents the genesis reader
	GenReader struct {
		writer      *jsonwriter.Writer
		TarballPath string
	}
	// UpdateGenesisOption configures genesis update.
	UpdateGenesisOption func(*ChainGenesis)
)

// FromPath parse GenReader object from path
func FromPath(genesisPath string) (*GenReader, error) {
	file, err := os.OpenFile(genesisPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open genesis file")
	}
	return New(file), nil
}

// FromURL fetches the genesis from the given URL and returns its content.
func FromURL(ctx context.Context, url string) (*GenReader, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out io.Writer
	tarballPath, err := tarball.ExtractFile(resp.Body, out, genesisFilename)
	if err != nil {
		return nil, err
	}
	return &GenReader{
		writer:      jsonwriter.New(out),
		TarballPath: tarballPath,
	}, nil
}

func New(genesis io.Writer) *GenReader {
	return &GenReader{
		writer: jsonwriter.New(genesis),
	}
}

func (g *GenReader) UpdateJSONWriter(key, value string) {
	g.writer.RootObject(func() { g.writer.KeyValue(key, value) })
}

// SaveGenesis save the genesis file
func (g *GenReader) SaveGenesis(filePath string) error {
	var (
		genesis ChainGenesis
	)
	g.Genesis()

	file, err := os.Create(filePath)
	file.
	//file := os.WriteFile(filePath, genesisBytes, 0644)
}
