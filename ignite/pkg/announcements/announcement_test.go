package announcements_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ignite/cli/v29/ignite/pkg/announcements"
)

func TestFetchAnnouncements(t *testing.T) {
	fallbackData := fmt.Sprintf("\n💬 Survey: %s\n", announcements.SurveyLink)

	tests := []struct {
		name         string
		mockResponse string
		statusCode   int
		expected     string
	}{
		{
			name:         "successful retrieval",
			mockResponse: `{"version":1,"announcements":[{"id":"1744230503810","text":"New Ignite announcement: v1.0.0 released!","timestamp":"2025-04-09T20:28:23.810Z","user":"announcement-bot"}]}`,
			statusCode:   http.StatusOK,
			expected:     "\n🗣️ Announcements\n⋆ New Ignite announcement: v1.0.0 released!\n",
		},
		{
			name:         "empty announcements",
			mockResponse: `{"announcements":[]}`,
			statusCode:   http.StatusOK,
			expected:     fallbackData,
		},
		{
			name:         "invalid JSON response",
			mockResponse: `invalid json`,
			statusCode:   http.StatusOK,
			expected:     fallbackData,
		},
		{
			name:         "non-200 HTTP response",
			mockResponse: ``,
			statusCode:   http.StatusInternalServerError,
			expected:     fallbackData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			originalAPI := announcements.APIURL
			announcements.APIURL = server.URL
			defer func() { announcements.APIURL = originalAPI }()

			result := announcements.Fetch()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
