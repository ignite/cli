package dircache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func Test_cacheKey(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	wd = filepath.Join(wd, "testdata")

	type args struct {
		src  string
		keys []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "no keys",
			args: args{
				src: wd,
			},
			want: "4cf0539ac24f8ebc9ee17b81d0ea880e55d2ba98a4e355affe3c3f8a0cdb01ee",
		},
		{
			name: "one key",
			args: args{
				src:  wd,
				keys: []string{"test"},
			},
			want: "dc7b4e68b7b9d827b3833845202818a11a1105542a3551052c012d815a64e7ae",
		},
		{
			name: "two keys",
			args: args{
				src:  wd,
				keys: []string{"test1", "test2"},
			},
			want: "a017b975dd0a30efc7fbc515af9b3c37657c20a509fd5771111d4c0e43d373b0",
		},
		{
			name: "duplicated keys",
			args: args{
				src:  wd,
				keys: []string{"test", "test"},
			},
			want: "26ce20a6c4563963fd646121948cd62137a143317c970a52a3ec8ed9979c868d",
		},
		{
			name: "many keys",
			args: args{
				src:  wd,
				keys: []string{"test1", "test2", "test3", "test4", "test5", "test6", "test6"},
			},
			want: "f9cd1468363ff902bdd5a93c9c7c43c83c9074796486306a7da046a082314121",
		},
		{
			name: "invalid source",
			args: args{
				src: "invalid_source",
			},
			err: errors.New("no file in specified paths"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cacheKey(tt.args.src, tt.args.keys...)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
