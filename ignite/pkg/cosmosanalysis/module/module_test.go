package module_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
)

var testModule = module.Module{
	Name:         "planet",
	GoModulePath: "github.com/tendermint/planet",
	Pkg: protoanalysis.Package{
		Name:         "tendermint.planet.planet",
		Path:         "testdata/planet/proto/planet",
		Files:        protoanalysis.Files{protoanalysis.File{Path: "testdata/planet/proto/planet/planet.proto", Dependencies: []string{"google/api/annotations.proto"}}},
		GoImportName: "github.com/tendermint/planet/x/planet/types",
		Messages: []protoanalysis.Message{
			{Name: "QueryMyQueryRequest", Path: "testdata/planet/proto/planet/planet.proto", HighestFieldNumber: 1},
			{Name: "QueryMyQueryResponse", Path: "testdata/planet/proto/planet/planet.proto", HighestFieldNumber: 0},
		},
		Services: []protoanalysis.Service{
			{
				Name: "Query",
				RPCFuncs: []protoanalysis.RPCFunc{
					{
						Name:        "MyQuery",
						RequestType: "QueryMyQueryRequest",
						ReturnsType: "QueryMyQueryResponse",
						HTTPRules: []protoanalysis.HTTPRule{
							{
								Params:   []string{"mytypefield"},
								HasQuery: false, HasBody: false,
							},
						},
					},
				},
			},
		},
	},
	Msgs: []module.Msg(nil),
	HTTPQueries: []module.HTTPQuery{
		{
			Name:     "MyQuery",
			FullName: "QueryMyQuery",
			Rules: []protoanalysis.HTTPRule{
				{
					Params:   []string{"mytypefield"},
					HasQuery: false,
					HasBody:  false,
				},
			},
		},
	},
	Types: []module.Type(nil),
}

func TestDiscover(t *testing.T) {
	type args struct {
		sourcePath string
		protoDir   string
	}
	tests := []struct {
		name string
		args args
		want []module.Module
	}{
		{
			name: "test valid",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "proto",
			},
			want: []module.Module{testModule},
		}, {
			name: "test no proto folder",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "",
			},
			want: []module.Module{testModule},
		}, {
			name: "test invalid proto folder",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "invalid",
			},
			want: []module.Module{},
		}, {
			name: "test invalid folder",
			args: args{
				sourcePath: "testdata/invalid",
				protoDir:   "",
			},
			want: []module.Module{},
		}, {
			name: "test invalid main and proto folder",
			args: args{
				sourcePath: "../../..",
				protoDir:   "proto",
			},
			want: []module.Module{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := module.Discover(context.Background(), "testdata/planet", tt.args.sourcePath, tt.args.protoDir)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsRootPath(t *testing.T) {
	cases := []struct {
		name, path string
		want       bool
	}{
		{
			name: "custom module import path",
			path: "github.com/chain/x/my_module",
			want: true,
		},
		{
			name: "generic import path",
			path: "github.com/username/project",
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, module.IsRootPath(tt.path))
		})
	}
}

func TestRootPath(t *testing.T) {
	cases := []struct {
		name, path, want string
	}{
		{
			name: "custom module import path",
			path: "github.com/username/chain/x/my_module/child/folder",
			want: "github.com/username/chain/x/my_module",
		},
		{
			name: "generic import path",
			path: "github.com/username/project/child/folder",
			want: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, module.RootPath(tt.path))
		})
	}
}

func TestRootGoImportPath(t *testing.T) {
	cases := []struct {
		name, path, want string
	}{
		{
			name: "go import path with version suffix",
			path: "github.com/username/chain/v2",
			want: "github.com/username/chain",
		},
		{
			name: "go import path without version suffix",
			path: "github.com/username/chain",
			want: "github.com/username/chain",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, module.RootGoImportPath(tt.path))
		})
	}
}
