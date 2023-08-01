package clictx_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/clictx"
)

func TestDo(t *testing.T) {
	ctxCanceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name        string
		ctx         context.Context
		f           func() error
		expectedErr string
	}{
		{
			name: "f returns nil",
			ctx:  context.Background(),
			f:    func() error { return nil },
		},
		{
			name:        "f returns an error",
			ctx:         context.Background(),
			f:           func() error { return errors.New("oups") },
			expectedErr: "oups",
		},
		{
			name: "ctx is canceled",
			ctx:  ctxCanceled,
			f: func() error {
				time.Sleep(time.Second)
				return nil
			},
			expectedErr: context.Canceled.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := clictx.Do(tt.ctx, tt.f)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
