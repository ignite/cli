package module_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis"
)

func newModule(relChainPath, goImportPath string) module.Module {
	return module.Module{
		Name:         "mars",
		GoModulePath: goImportPath,
		Pkg: protoanalysis.Package{
			Name: "tendermint.planet.mars",
			Path: filepath.Join(relChainPath, "proto/planet/mars"),
			Files: protoanalysis.Files{
				protoanalysis.File{
					Path: filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					Dependencies: []string{
						"cosmos/base/query/v1beta1/pagination.proto",
						"google/api/annotations.proto",
					},
				},
			},
			GoImportName: "github.com/tendermint/planet/x/mars/types",
			Messages: []protoanalysis.Message{
				{
					Name:               "MsgMyMessageRequest",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 1,
					Fields: map[string]string{
						"mytypefield": "string",
					},
				},
				{
					Name:               "MsgMyMessageResponse",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 1,
					Fields: map[string]string{
						"mytypefield": "string",
					},
				},
				{
					Name:               "MsgBarRequest",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 1,
					Fields: map[string]string{
						"mytypefield": "string",
					},
				},
				{
					Name:               "MsgBarResponse",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 1,
					Fields: map[string]string{
						"mytypefield": "string",
					},
				},
				{
					Name:               "QueryMyQueryRequest",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 2,
					Fields: map[string]string{
						"mytypefield": "string",
						"pagination":  "cosmos.base.query.v1beta1.PageRequest",
					},
				},
				{
					Name:               "QueryMyQueryResponse",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 1,
					Fields:             map[string]string{"pagination": "cosmos.base.query.v1beta1.PageResponse"},
				},
				{
					Name:               "QueryFooRequest",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 0,
					Fields:             map[string]string{},
				},
				{
					Name:               "QueryFooResponse",
					Path:               filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
					HighestFieldNumber: 1,
					Fields:             map[string]string{"bar": "string"},
				},
			},
			Services: []protoanalysis.Service{
				{
					Name: "Msg",
					RPCFuncs: []protoanalysis.RPCFunc{
						{
							Name:        "MyMessage",
							RequestType: "MsgMyMessageRequest",
							ReturnsType: "MsgMyMessageResponse",
						},
						{
							Name:        "Bar",
							RequestType: "MsgBarRequest",
							ReturnsType: "MsgBarResponse",
						},
					},
				},
				{
					Name: "Query",
					RPCFuncs: []protoanalysis.RPCFunc{
						{
							Name:        "MyQuery",
							RequestType: "QueryMyQueryRequest",
							ReturnsType: "QueryMyQueryResponse",
							HTTPRules: []protoanalysis.HTTPRule{
								{
									Endpoint: "/tendermint/mars/my_query/{mytypefield}",
									Params:   []string{"mytypefield"},
									HasQuery: true,
									QueryFields: map[string]string{
										"pagination": "cosmos.base.query.v1beta1.PageRequest",
									},
									HasBody: false,
								},
							},
						},
						{
							Name:        "Foo",
							RequestType: "QueryFooRequest",
							ReturnsType: "QueryFooResponse",
							HTTPRules: []protoanalysis.HTTPRule{
								{
									Endpoint: "/tendermint/mars/foo",
									HasQuery: false,
									HasBody:  false,
								},
							},
						},
					},
				},
			},
		},
		Msgs: []module.Msg{
			{
				Name:     "MsgMyMessageRequest",
				URI:      "tendermint.planet.mars.MsgMyMessageRequest",
				FilePath: filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
			},
			{
				Name:     "MsgBarRequest",
				URI:      "tendermint.planet.mars.MsgBarRequest",
				FilePath: filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
			},
		},
		HTTPQueries: []module.HTTPQuery{
			{
				Name:         "MyQuery",
				FullName:     "QueryMyQuery",
				RequestType:  "QueryMyQueryRequest",
				ResponseType: "QueryMyQueryResponse",
				Rules: []protoanalysis.HTTPRule{
					{
						Endpoint: "/tendermint/mars/my_query/{mytypefield}",
						Params:   []string{"mytypefield"},
						HasQuery: true,
						QueryFields: map[string]string{
							"pagination": "cosmos.base.query.v1beta1.PageRequest",
						},
						HasBody: false,
					},
				},
				Paginated: true,
				FilePath:  filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
			},
			{
				Name:         "Foo",
				FullName:     "QueryFoo",
				RequestType:  "QueryFooRequest",
				ResponseType: "QueryFooResponse",
				Rules: []protoanalysis.HTTPRule{
					{
						Endpoint: "/tendermint/mars/foo",
						HasQuery: false,
						HasBody:  false,
					},
				},
				FilePath: filepath.Join(relChainPath, "proto/planet/mars/mars.proto"),
			},
		},
		Types: []module.Type(nil),
	}
}

func TestDiscover(t *testing.T) {
	ctx := context.Background()
	sourcePath := "testdata/planet"
	testModule := newModule(sourcePath, "github.com/tendermint/planet")

	tests := []struct {
		name, sourcePath, protoDir string
		want                       []module.Module
	}{
		{
			name:       "test valid",
			sourcePath: sourcePath,
			protoDir:   "proto",
			want:       []module.Module{testModule},
		}, {
			name:       "test no proto folder",
			sourcePath: sourcePath,
			protoDir:   "",
			want:       []module.Module{testModule},
		}, {
			name:       "test invalid proto folder",
			sourcePath: sourcePath,
			protoDir:   "invalid",
			want:       []module.Module{},
		}, {
			name:       "test invalid folder",
			sourcePath: "testdata/invalid",
			protoDir:   "",
			want:       []module.Module{},
		}, {
			name:       "test invalid main and proto folder",
			sourcePath: "../../..",
			protoDir:   "proto",
			want:       []module.Module{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules, err := module.Discover(ctx, sourcePath, tt.sourcePath, module.WithProtoDir(tt.protoDir))

			require.NoError(t, err)
			require.Equal(t, tt.want, modules)
		})
	}
}

func TestDiscoverWithAppV2(t *testing.T) {
	ctx := context.Background()
	sourcePath := "testdata/earth"
	testModule := newModule(sourcePath, "github.com/tendermint/planet")

	tests := []struct {
		name, protoDir string
		want           []module.Module
	}{
		{
			name:     "test valid",
			protoDir: "proto",
			want:     []module.Module{testModule},
		}, {
			name:     "test valid with version suffix",
			protoDir: "proto",
			want:     []module.Module{testModule},
		}, {
			name:     "test no proto folder",
			protoDir: "",
			want:     []module.Module{testModule},
		}, {
			name:     "test invalid proto folder",
			protoDir: "invalid",
			want:     []module.Module{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules, err := module.Discover(ctx, sourcePath, sourcePath, module.WithProtoDir(tt.protoDir))

			require.NoError(t, err)
			require.Equal(t, tt.want, modules)
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
