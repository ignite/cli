package clidoc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	build struct {
		Main     string   `yaml:"main,omitempty" doc:"doc of main"`
		Binary   string   `yaml:"binary,omitempty" doc:""`
		LDFlags  []string `yaml:"ldflags,omitempty"`
		Proto    proto    `yaml:"proto" doc:"doc of proto"`
		PtrProto *proto   `yaml:"ptr_proto" doc:"doc of pointer proto"`
		Protos   []proto  `yaml:"protos" doc:"doc of protos"`
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
		want Docs
		err  error
	}{
		{
			name: "build struct",
			v:    build{},
			want: Docs{
				{
					Key:     "main",
					Comment: "doc of main",
					Type:    "string",
				},
				{
					Key:  "binary",
					Type: "string",
				},
				{
					Key:  "ldflags",
					Type: listName("string"),
				},
				{
					Key: "proto",
					Value: Docs{
						{
							Key:     "path",
							Comment: "path of proto file",
							Type:    "string",
						},
						{
							Key:     "third_party_paths",
							Comment: "doc of third party paths",
							Type:    listName("string"),
						},
					},
					Comment: "doc of proto",
				},
				{
					Key: "ptr_proto",
					Value: Docs{
						{
							Key:     "path",
							Comment: "path of proto file",
							Type:    "string",
						},
						{
							Key:     "third_party_paths",
							Comment: "doc of third party paths",
							Type:    listName("string"),
						},
					},
					Comment: "doc of pointer proto",
				},
				{
					Key:  "protos",
					Type: "list",
					Value: Docs{
						{
							Key:     "path",
							Comment: "path of proto file",
							Type:    "string",
						},
						{
							Key:     "third_party_paths",
							Comment: "doc of third party paths",
							Type:    listName("string"),
						},
					},
					Comment: "doc of protos",
				},
			},
		},
		{
			name: "proto struct",
			v:    proto{},
			want: Docs{
				{
					Key:     "path",
					Comment: "path of proto file",
					Type:    "string",
				},
				{
					Key:     "third_party_paths",
					Comment: "doc of third party paths",
					Type:    listName("string"),
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

func TestDocs_String(t *testing.T) {
	tests := []struct {
		name string
		d    Docs
		want string
	}{
		{
			name: "many entries",
			d: Docs{
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
					Value: Docs{
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
					Value: Docs{
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
			want: `
main: # doc of main
binary: # 
ldflags [array]: # 
proto: # doc of proto
  path: # path of proto file
  third_party_paths [array]: # doc of third party paths
protos [array]: # doc of protos
  path: # path of proto file
  third_party_paths [array]: # doc of third party paths`,
		},
		{
			name: "no entries",
			d:    Docs{},
		},
		{
			name: "two entries",
			d: Docs{
				{
					Key:     "path",
					Comment: "path of proto file",
				},
				{
					Key:     "third_party_paths [array]",
					Comment: "doc of third party paths",
				},
			},
			want: `
path: # path of proto file
third_party_paths [array]: # doc of third party paths`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.String()
			require.Equal(t, strings.TrimSpace(tt.want), strings.TrimSpace(got))
		})
	}
}
