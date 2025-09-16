package cosmosgen

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ettle/strcase"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
)

func TestGenerateTypeScript(t *testing.T) {
	require := require.New(t)
	testdataDir := "testdata"
	appDir := filepath.Join(testdataDir, "testchain")
	tsClientDir := filepath.Join(appDir, "ts-client")

	cacheStorage, err := cache.NewStorage(filepath.Join(t.TempDir(), "cache.db"))
	require.NoError(err)

	buf, err := cosmosbuf.New(cacheStorage, t.Name())
	require.NoError(err)

	// Use module discovery to collect test module proto.
	m, err := module.Discover(t.Context(), appDir, appDir, module.WithProtoDir("proto"))
	require.NoError(err, "failed to discover module")
	require.Len(m, 1, "expected exactly one module to be discovered")

	g := newTSGenerator(&generator{
		appPath:      appDir,
		protoDir:     "proto",
		goModPath:    "go.mod",
		cacheStorage: cacheStorage,
		buf:          buf,
		appModules:   m,
		opts: &generateOptions{
			tsClientRootPath: tsClientDir,
			useCache:         false,
			jsOut: func(m module.Module) string {
				return filepath.Join(tsClientDir, fmt.Sprintf("%s.%s.%s", "ignite", "planet", strcase.ToKebab(m.Name)))
			},
		},
	})

	err = g.generateModuleTemplate(t.Context(), appDir, m[0])
	require.NoError(err, "failed to generate TypeScript files")

	err = g.generateRootTemplates(generatePayload{
		Modules:   m,
		PackageNS: strings.ReplaceAll(appDir, "/", "-"),
	})
	require.NoError(err)

	// compare all generated files to golden files
	goldenDir := filepath.Join(testdataDir, "expected_files", "ts-client")
	_ = filepath.Walk(goldenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(goldenDir, path)
		got := filepath.Join(tsClientDir, rel)
		gold, err := os.ReadFile(path)
		require.NoError(err, "failed to read golden file: %s", path)

		gotBytes, err := os.ReadFile(got)
		require.NoError(err, "failed to read generated file: %s", got)
		require.Equal(string(gold), string(gotBytes), "file %s does not match golden file", rel)

		return nil
	})
}

func TestFetchBufToken(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		expectedToken  string
		expectError    bool
	}{
		{
			name:           "successful token fetch",
			serverResponse: `{"token":"test_token_123"}`,
			statusCode:     http.StatusOK,
			expectedToken:  "test_token_123",
			expectError:    false,
		},
		{
			name:           "server error",
			serverResponse: `{"error":"internal server error"}`,
			statusCode:     http.StatusInternalServerError,
			expectedToken:  "",
			expectError:    true,
		},
		{
			name:           "invalid json response",
			serverResponse: `invalid json`,
			statusCode:     http.StatusOK,
			expectedToken:  "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Temporarily override the endpoint
			originalEndpoint := bufTokenEndpoint
			bufTokenEndpoint = server.URL
			defer func() {
				bufTokenEndpoint = originalEndpoint
			}()

			token, err := fetchBufToken()

			if tt.expectError {
				require.Error(t, err)
				require.Empty(t, token)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestNewTSGeneratorBufTokenLogic(t *testing.T) {
	tests := []struct {
		name               string
		isLocalProto       bool
		existingEnvToken   string
		mockServerResponse string
		mockStatusCode     int
		expectEnvSet       bool
		expectLogMessage   bool
	}{
		{
			name:               "local proto available - no token fetch",
			isLocalProto:       true,
			existingEnvToken:   "",
			mockServerResponse: "",
			mockStatusCode:     0,
			expectEnvSet:       false,
			expectLogMessage:   false,
		},
		{
			name:               "env token already set - no fetch",
			isLocalProto:       false,
			existingEnvToken:   "existing_token",
			mockServerResponse: "",
			mockStatusCode:     0,
			expectEnvSet:       false,
			expectLogMessage:   false,
		},
		{
			name:               "successful token fetch",
			isLocalProto:       false,
			existingEnvToken:   "",
			mockServerResponse: `{"token":"fetched_token"}`,
			mockStatusCode:     http.StatusOK,
			expectEnvSet:       true,
			expectLogMessage:   false,
		},
		{
			name:               "failed token fetch - log message",
			isLocalProto:       false,
			existingEnvToken:   "",
			mockServerResponse: `{"error":"server error"}`,
			mockStatusCode:     http.StatusInternalServerError,
			expectEnvSet:       false,
			expectLogMessage:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			originalToken := os.Getenv(bufTokenEnvName)
			originalEndpoint := bufTokenEndpoint
			defer func() {
				if originalToken != "" {
					os.Setenv(bufTokenEnvName, originalToken)
				} else {
					os.Unsetenv(bufTokenEnvName)
				}
				bufTokenEndpoint = originalEndpoint
			}()

			// Set up test environment
			if tt.existingEnvToken != "" {
				os.Setenv(bufTokenEnvName, tt.existingEnvToken)
			} else {
				os.Unsetenv(bufTokenEnvName)
			}

			// Set up mock server if needed
			if !tt.isLocalProto && tt.existingEnvToken == "" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockStatusCode)
					w.Write([]byte(tt.mockServerResponse))
				}))
				defer server.Close()
				bufTokenEndpoint = server.URL
			}

			// Create mock generator
			g := &generator{}

			// Create TSGenerator with mocked isLocalProto
			tsg := &tsGenerator{
				g:            g,
				isLocalProto: tt.isLocalProto,
			}

			// Simulate the logic from newTSGenerator
			if !tsg.isLocalProto {
				if os.Getenv(bufTokenEnvName) == "" {
					token, err := fetchBufToken()
					if err != nil {
						// This would normally log the message
						if tt.expectLogMessage {
							require.Error(t, err)
						}
					} else {
						os.Setenv(bufTokenEnvName, token)
					}
				}
			}

			// Verify expectations
			if tt.expectEnvSet {
				token := os.Getenv(bufTokenEnvName)
				require.NotEmpty(t, token)
				require.Equal(t, "fetched_token", token)
			} else if tt.existingEnvToken != "" {
				token := os.Getenv(bufTokenEnvName)
				require.Equal(t, tt.existingEnvToken, token)
			}
		})
	}
}
