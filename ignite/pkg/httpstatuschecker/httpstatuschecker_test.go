package httpstatuschecker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckStatus(t *testing.T) {
	cases := []struct {
		name           string
		returnedStatus int
		isAvaiable     bool
	}{
		{"200 OK", 200, true},
		{"202 Accepted ", 202, true},
		{"404 Not Found", 404, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.returnedStatus)
			}))
			defer ts.Close()

			isAvailable, err := Check(context.Background(), ts.URL)
			require.NoError(t, err)
			require.Equal(t, tt.isAvaiable, isAvailable)
		})
	}
}

func TestCheckServerUnreachable(t *testing.T) {
	isAvailable, err := Check(context.Background(), "http://localhost:63257")
	require.NoError(t, err)
	require.False(t, isAvailable)
}
