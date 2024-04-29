package clidoc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	build struct {
		Main    string   `yaml:"main,omitempty" doc:"doc of main"`
		Binary  string   `yaml:"binary,omitempty" doc:""`
		LDFlags []string `yaml:"ldflags,omitempty"`
		Proto   proto    `yaml:"proto" doc:"doc of proto"`
		Protos  []proto  `yaml:"protos" doc:"doc of protos"`
	}
	proto struct {
		Path            string   `yaml:"path" doc:"path of proto file"`
		ThirdPartyPaths []string `yaml:"third_party_paths" doc:"doc of third party paths"`
	}
)

func TestGenDoc(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want []Doc
		err  error
	}{
		{
			name: "Build struct",
			v:    build{},
			want: []Doc{
				{
					Key:     "main",
					Comment: "doc of main",
				},
				{
					Key: "binary",
				},
				{
					Key: "ldflags [array]",
				},
				{
					Key: "proto",
					Value: []Doc{
						{
							Key:     "path",
							Comment: "path of proto file",
						},
						{
							Key:     "third_party_paths [array]",
							Comment: "doc of third party paths",
						},
					},
					Comment: "doc of proto",
				},
				{
					Key: "protos [array]",
					Value: []Doc{
						{
							Key:     "path",
							Comment: "path of proto file",
						},
						{
							Key:     "third_party_paths [array]",
							Comment: "doc of third party paths",
						},
					},
					Comment: "doc of protos",
				},
			},
		},
		{
			name: "Proto struct",
			v:    proto{},
			want: []Doc{
				{
					Key:     "path",
					Comment: "path of proto file",
				},
				{
					Key:     "third_party_paths [array]",
					Comment: "doc of third party paths",
				},
			},
		},
		{
			name: "Invalid struct",
			v:    []map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDoc(tt.v)
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
