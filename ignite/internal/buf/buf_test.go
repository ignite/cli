package buf_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/internal/buf"
)

func TestFetchToken(t *testing.T) {
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
			originalEndpoint := buf.BufTokenURL
			buf.BufTokenURL = server.URL
			defer func() {
				buf.BufTokenURL = originalEndpoint
			}()

			token, err := buf.FetchToken()
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
