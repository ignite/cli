package module

import (
	"context"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
	"testing"

	"github.com/stretchr/testify/require"
)

var testModule = Module{
	Name: "planet",
	Pkg: protoanalysis.Package{
		Name:         "tendermint.planet.planet",
		Path:         "testdata/planet/proto/planet",
		Files:        protoanalysis.Files{protoanalysis.File{Path: "testdata/planet/proto/planet/planet.proto", Dependencies: []string{"google/api/annotations.proto"}}},
		GoImportName: "github.com/tendermint/planet/x/planet/types",
		Messages: []protoanalysis.Message{
			protoanalysis.Message{Name: "QueryMyQueryRequest", Path: "testdata/planet/proto/planet/planet.proto", HighestFieldNumber: 1},
			protoanalysis.Message{Name: "QueryMyQueryResponse", Path: "testdata/planet/proto/planet/planet.proto", HighestFieldNumber: 0},
		},
		Services: []protoanalysis.Service{
			protoanalysis.Service{
				Name: "Query",
				RPCFuncs: []protoanalysis.RPCFunc{
					protoanalysis.RPCFunc{
						Name:        "MyQuery",
						RequestType: "QueryMyQueryRequest",
						ReturnsType: "QueryMyQueryResponse",
						HTTPRules: []protoanalysis.HTTPRule{
							protoanalysis.HTTPRule{
								Params:   []string{"mytypefield"},
								HasQuery: false, HasBody: false},
						},
					},
				},
			},
		},
	},
	Msgs: []Msg(nil),
	HTTPQueries: []HTTPQuery{
		HTTPQuery{
			Name:     "MyQuery",
			FullName: "QueryMyQuery",
			Rules: []protoanalysis.HTTPRule{
				protoanalysis.HTTPRule{
					Params:   []string{"mytypefield"},
					HasQuery: false,
					HasBody:  false},
			},
		},
	},
	Types: []Type(nil),
}

func TestDiscover(t *testing.T) {
	type args struct {
		sourcePath string
		protoDir   string
	}
	tests := []struct {
		name string
		args args
		want []Module
	}{
		{
			name: "test valid",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "proto",
			},
			want: []Module{testModule},
		}, {
			name: "test no proto folder",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "",
			},
			want: []Module{testModule},
		}, {
			name: "test invalid proto folder",
			args: args{
				sourcePath: "testdata/planet",
				protoDir:   "invalid",
			},
			want: []Module{},
		}, {
			name: "test invalid folder",
			args: args{
				sourcePath: "testdata/invalid",
				protoDir:   "",
			},
			want: []Module{},
		}, {
			name: "test invalid main and proto folder",
			args: args{
				sourcePath: "../../..",
				protoDir:   "proto",
			},
			want: []Module{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Discover(context.Background(), "testdata/planet", tt.args.sourcePath, tt.args.protoDir)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
