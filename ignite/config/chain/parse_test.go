package chain_test

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/config/chain/version"
	"github.com/ignite/cli/v29/ignite/config/testdata"
	"github.com/ignite/cli/v29/ignite/pkg/availableport"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func TestReadConfigVersion(t *testing.T) {
	// Arrange
	r := strings.NewReader("version: 42")
	want := version.Version(42)

	// Act
	v, err := chainconfig.ReadConfigVersion(r)

	// Assert
	require.NoError(t, err)
	require.Equal(t, want, v)
}

func TestParse(t *testing.T) {
	// Arrange: Initialize a reader with the previous version
	ver := chainconfig.LatestVersion - 1
	r := bytes.NewReader(testdata.Versions[ver])

	// Act
	cfg, err := chainconfig.Parse(r)

	// Assert
	require.NoError(t, err)

	// Assert: Parse must return the latest version
	require.Equal(t, chainconfig.LatestVersion, cfg.Version)
	require.Equal(t, testdata.GetLatestConfig(t), cfg)
}

func TestParseWithCurrentVersion(t *testing.T) {
	// Arrange
	r := bytes.NewReader(testdata.Versions[chainconfig.LatestVersion])

	// Act
	cfg, err := chainconfig.Parse(r)

	// Assert
	require.NoError(t, err)
	require.Equal(t, chainconfig.LatestVersion, cfg.Version)
	require.Equal(t, testdata.GetLatestConfig(t), cfg)
}

func TestParseWithUnknownVersion(t *testing.T) {
	// Arrange
	version := version.Version(9999)
	r := strings.NewReader(fmt.Sprintf("version: %d", version))

	var want *chainconfig.UnsupportedVersionError

	// Act
	_, err := chainconfig.Parse(r)

	// Assert
	require.ErrorAs(t, err, &want)
	require.NotNil(t, want)
	require.Equal(t, want.Version, version)
}

func TestParseNetworkWithCurrentVersion(t *testing.T) {
	// Arrange
	r := bytes.NewReader(testdata.NetworkConfig)

	// Act
	cfg, err := chainconfig.ParseNetwork(r)

	// Assert
	require.NoError(t, err)

	// Assert: Parse must return the latest version
	require.Equal(t, chainconfig.LatestVersion, cfg.Version)
	require.Equal(t, testdata.GetLatestNetworkConfig(t).Accounts, cfg.Accounts)
	require.Equal(t, testdata.GetLatestNetworkConfig(t).Genesis, cfg.Genesis)
}

func TestParseNetworkWithInvalidData(t *testing.T) {
	// Arrange
	r := bytes.NewReader(testdata.Versions[chainconfig.LatestVersion])

	// Act
	_, err := chainconfig.ParseNetwork(r)

	// Assert error
	require.True(
		t,
		strings.Contains(
			err.Error(),
			"config is not valid: no validators can be used in config for network genesis",
		),
	)
}

func TestHandleIncludes(t *testing.T) {
	server := startTestServer(t)

	tests := []struct {
		name       string
		baseConfig string
		expected   string
		err        error
	}{
		{
			name: "Single valid include",
			baseConfig: `
version: 1
client:
  typescript:
    path: original
include:
  - "./testdata/include1.yml"
`,

			expected: `
include:
    - ./testdata/include1.yml
validation: sovereign
version: 1
build:
    proto:
        path: proto
accounts:
    - name: bob
      coins:
        - 10000token
        - 100000000stake
faucet:
    name: danilo
    coins:
        - 5token
        - 100000stake
    host: 0.0.0.0:4500
client:
    typescript:
        path: override-1
    openapi:
        path: docs/static/include1.yml
validators:
    - name: alice
      bonded: 100000000stake`,
		},
		{
			name: "Multiple valid includes",
			baseConfig: `
version: 1
client:
  typescript:
    path: original
include:
  - "./testdata/include1.yml"
  - "./testdata/include2.yml"
`,
			expected: `
include:
    - ./testdata/include1.yml
    - ./testdata/include2.yml
validation: sovereign
version: 1
build:
    proto:
        path: proto
accounts:
    - name: bob
      coins:
        - 10000token
        - 100000000stake
    - name: alice
      coins:
        - 20000token
        - 200000000stake
faucet:
    name: alice
    coins:
        - 5token
        - 100000stake
        - 5token
        - 100000stake
    host: 0.0.0.0:4500
client:
    typescript:
        path: override-2
    openapi:
        path: docs/static/include2.yml
validators:
    - name: alice
      bonded: 100000000stake
    - name: validator1
      bonded: 100000000stake`,
		},
		{
			name: "Invalid include file path",
			baseConfig: `
version: 1
client:
  typescript:
    path: original
include:
  - "./testdata/nonexistent.yml"`,
			err: errors.New("error parsing config file: failed to open included file './testdata/nonexistent.yml'"),
		},
		{
			name: "Empty include list",
			baseConfig: `
version: 1
validation: sovereign
accounts:
  - name: alice
    coins:
      - 10000token
include: []
faucet:
  name: alice
  coins:
    - 5token
client:
  typescript:
    path: override-1
  openapi:
    path: docs/static/include1.yml
validators:
  - name: alice
    bonded: 100stake`,
			expected: `
validation: sovereign
version: 1
build:
    proto:
        path: proto
accounts:
    - name: alice
      coins:
        - 10000token
faucet:
    name: alice
    coins:
        - 5token
    host: 0.0.0.0:4500
client:
    typescript:
        path: override-1
    openapi:
        path: docs/static/include1.yml
validators:
    - name: alice
      bonded: 100stake`,
		},
		{
			name: "Empty values include",
			baseConfig: `
version: 1
include:
  - "./testdata/include1.yml"
  - "./testdata/include2.yml"
`,
			expected: `
include:
    - ./testdata/include1.yml
    - ./testdata/include2.yml
validation: sovereign
version: 1
build:
    proto:
        path: proto
accounts:
    - name: bob
      coins:
        - 10000token
        - 100000000stake
    - name: alice
      coins:
        - 20000token
        - 200000000stake
faucet:
    name: alice
    coins:
        - 5token
        - 100000stake
        - 5token
        - 100000stake
    host: 0.0.0.0:4500
client:
    typescript:
        path: override-2
    openapi:
        path: docs/static/include2.yml
validators:
    - name: alice
      bonded: 100000000stake
    - name: validator1
      bonded: 100000000stake`,
		},
		{
			name: "HTTP include",
			baseConfig: fmt.Sprintf(`
version: 1
validation: sovereign
accounts:
  - name: alice
    coins:
      - 10000token
include:
    - %[1]v/include1.yml
    - %[1]v/include2.yml
faucet:
  name: alice
  coins:
    - 5token
client:
  typescript:
    path: original
  openapi:
    path: docs/static/include1.yml
validators:
  - name: alice
    bonded: 100stake
`, server),
			expected: fmt.Sprintf(`
include:
    - %[1]v/include1.yml
    - %[1]v/include2.yml
validation: sovereign
version: 1
build:
    proto:
        path: proto
accounts:
    - name: alice
      coins:
        - 10000token
    - name: bob
      coins:
        - 10000token
        - 100000000stake
    - name: alice
      coins:
        - 20000token
        - 200000000stake
faucet:
    name: alice
    coins:
        - 5token
        - 5token
        - 100000stake
        - 5token
        - 100000stake
    host: 0.0.0.0:4500
client:
    typescript:
        path: override-2
    openapi:
        path: docs/static/include2.yml
validators:
    - name: alice
      bonded: 100stake
    - name: alice
      bonded: 100000000stake
    - name: validator1
      bonded: 100000000stake`, server),
		},
		{
			name: "HTTP and local include",
			baseConfig: fmt.Sprintf(`
version: 1
validation: sovereign
accounts:
  - name: alice
    coins:
      - 10000token
include:
    - %[1]v/include1.yml
    - testdata/include2.yml
faucet:
  name: alice
  coins:
    - 5token
client:
  typescript:
    path: original
  openapi:
    path: docs/static/include1.yml
validators:
  - name: alice
    bonded: 100stake
`, server),
			expected: fmt.Sprintf(`
include:
    - %[1]v/include1.yml
    - testdata/include2.yml
validation: sovereign
version: 1
build:
    proto:
        path: proto
accounts:
    - name: alice
      coins:
        - 10000token
    - name: bob
      coins:
        - 10000token
        - 100000000stake
    - name: alice
      coins:
        - 20000token
        - 200000000stake
faucet:
    name: alice
    coins:
        - 5token
        - 5token
        - 100000stake
        - 5token
        - 100000stake
    host: 0.0.0.0:4500
client:
    typescript:
        path: override-2
    openapi:
        path: docs/static/include2.yml
validators:
    - name: alice
      bonded: 100stake
    - name: alice
      bonded: 100000000stake
    - name: validator1
      bonded: 100000000stake`, server),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseReader := bytes.NewReader([]byte(tt.baseConfig))
			baseConfig, err := chainconfig.Parse(baseReader)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)

			finalConfigYaml, err := yaml.Marshal(baseConfig)
			require.NoError(t, err)
			require.Equal(t, strings.TrimSpace(tt.expected), strings.TrimSpace(string(finalConfigYaml)))
		})
	}
}

func startTestServer(t *testing.T) string {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/include1.yml", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("testdata/include1.yml")
		require.NoError(t, err)
		_, err = w.Write(content)
		require.NoError(t, err)
	})
	mux.HandleFunc("/include2.yml", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile("testdata/include2.yml")
		require.NoError(t, err)
		_, err = w.Write(content)
		require.NoError(t, err)
	})

	ports, err := availableport.Find(1)
	require.NoError(t, err)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", ports[0]),
		Handler: mux,
	}

	listener, err := net.Listen("tcp", server.Addr)
	require.NoError(t, err)

	go server.Serve(listener)

	t.Cleanup(func() {
		server.Close()
	})

	return fmt.Sprintf("http://localhost:%d", ports[0])
}
