package announcements_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/announcements"
)

func TestGetAnnouncements(t *testing.T) {
	fallbackData := fmt.Sprintf("\n💬 Survey: %s\n", announcements.SurveyLink)

	tests := []struct {
		name         string
		mockResponse string
		statusCode   int
		expected     string
	}{
		{
			name:         "successful retrieval",
			mockResponse: `{"announcements":["Announcement 1","Announcement 2"]}`,
			statusCode:   http.StatusOK,
			expected:     "\n🗣️ Announcements\n⋆ Announcement 1\n⋆ Announcement 2\n",
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

			originalAPI := announcements.AnnouncementAPI
			announcements.AnnouncementAPI = server.URL
			defer func() { announcements.AnnouncementAPI = originalAPI }()

			result := announcements.GetAnnouncements()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
