package httpstatuschecker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestClient(fn roundTripperFunc) *http.Client {
	return &http.Client{Transport: fn}
}

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
			client := newTestClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode:    tt.returnedStatus,
					Body:          io.NopCloser(bytes.NewReader(nil)),
					Header:        make(http.Header),
					ContentLength: 0,
					Request:       req,
				}, nil
			})

			isAvailable, err := Check(context.Background(), "http://example.com", Client(client))
			require.NoError(t, err)
			require.Equal(t, tt.isAvaiable, isAvailable)
		})
	}
}

func TestCheckServerUnreachable(t *testing.T) {
	client := newTestClient(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("dial tcp: connection refused")
	})
	isAvailable, err := Check(context.Background(), "http://example.com", Client(client))
	require.NoError(t, err)
	require.False(t, isAvailable)
}
