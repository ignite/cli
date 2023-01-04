package cosmoserror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
)

func TestUnwrap(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "should return internal error",
			err:  status.Error(codes.Internal, "test error 1"),
			want: cosmoserror.ErrInternal,
		},
		{
			name: "should return invalid request",
			err:  status.Error(codes.InvalidArgument, "test error 2"),
			want: cosmoserror.ErrInvalidRequest,
		},
		{
			name: "should return not found",
			err:  status.Error(codes.NotFound, "test error 3"),
			want: cosmoserror.ErrNotFound,
		},
		{
			name: "should return not found with wrapped error",
			err:  fmt.Errorf("oups: %w", status.Error(codes.NotFound, "test error 4")),
			want: cosmoserror.ErrNotFound,
		},
		{
			name: "should return same error",
			err:  errors.New("test error 5"),
			want: errors.New("test error 5"),
		},
		{
			name: "should unwrap error",
			err:  fmt.Errorf("test error 4: %w", errors.New("test error 6")),
			want: errors.New("test error 6"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			unwrapped := cosmoserror.Unwrap(tc.err)
			require.Equal(t, tc.want, unwrapped)
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "should return false from invalid code",
			err:  status.Error(codes.Internal, "test error 1"),
			want: false,
		},
		{
			name: "should return false from invalid error",
			err:  errors.New("test error 4"),
			want: false,
		},
		{
			name: "should return true",
			err:  status.Error(codes.NotFound, "test error 3"),
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := cosmoserror.IsNotFound(tc.err)
			require.Equal(t, tc.want, got)
		})
	}
}
