package dircache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
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
			want: "78f544d2184b8076ac527ba4728822de1a7fc77bf2d6a77e44d0193cb63ed26e",
		},
		{
			name: "one key",
			args: args{
				src:  wd,
				keys: []string{"test"},
			},
			want: "5701099a1fcc67cd8b694295fbdecf537edcc8733bcc3adae0bdd7e65e28c8e5",
		},
		{
			name: "two keys",
			args: args{
				src:  wd,
				keys: []string{"test1", "test2"},
			},
			want: "6299c9bd405a1c073fa711006f8aadf6420cf522ef446e36fc01586354726095",
		},
		{
			name: "duplicated keys",
			args: args{
				src:  wd,
				keys: []string{"test", "test"},
			},
			want: "b9eb1b01931deccc44a354ab5aeb52337a465e5559069eb35b71ea0cbfe3c87f",
		},
		{
			name: "many keys",
			args: args{
				src:  wd,
				keys: []string{"test1", "test2", "test3", "test4", "test5", "test6", "test6"},
			},
			want: "bbe74cfd33ba4d1244e8d0ea3e430081d06ed55be12c7772d345d3117a4dfc90",
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
