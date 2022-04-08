package yaml

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	type byteSliceParser struct {
		Field1 string `json:"field1"`
		Field2 struct {
			Field1 []byte `json:"field1"`
			Field2 struct {
				Field1 []byte `json:"field1"`
				Field2 string `json:"field2"`
			} `json:"field2"`
			Field3 string `json:"field3"`
		} `json:"field2"`
		Field3 string `json:"field3"`
	}
	bParser := &byteSliceParser{
		Field1: "field1",
		Field3: "field3",
	}
	bParser.Field2.Field1 = []byte("field1")
	bParser.Field2.Field3 = "field3"
	bParser.Field2.Field2.Field1 = []byte("field1")
	bParser.Field2.Field2.Field2 = "field2"

	type simpleParser struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"`
	}
	sParser := &simpleParser{
		Field1: "field1",
		Field2: "field2",
	}

	type args struct {
		obj   interface{}
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "parse nil obj",
			want: "null",
		},
		{
			name: "parse map without byte slice",
			args: args{
				obj: map[string]string{
					"field1": "field1",
					"field2": "field2",
				},
			},
			want: `field1: field1
field2: field2`,
		},
		{
			name: "parse map with byte slice",
			args: args{
				obj: map[string][]byte{
					"field1": []byte("field1"),
					"field2": []byte("field2"),
				},
				paths: []string{
					"$.field1",
					"$.field2",
				},
			},
			want: `field1: field1
field2: field2`,
		},
		{
			name: "parse struct without byte slice",
			args: args{
				obj: sParser,
			},
			want: `field1: field1
field2: field2`,
		},
		{
			name: "parse struct with byte slice",
			args: args{
				obj: bParser,
				paths: []string{
					"$.field2.field1",
					"$.field2.field2.field1",
				},
			},
			want: `field1: field1
field2:
  field1: field1
  field2:
    field1: field1
    field2: field2
  field3: field3
field3: field3`,
		},
		{
			name: "parse struct with byte slice and wrong path",
			args: args{
				obj: bParser,
				paths: []string{
					"$.field2.field30",
					"$.field2.field31",
				},
			},
			want: `field1: field1
field2:
  field1:
  - 102
  - 105
  - 101
  - 108
  - 100
  - 49
  field2:
    field1:
    - 102
    - 105
    - 101
    - 108
    - 100
    - 49
    field2: field2
  field3: field3
field3: field3`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(context.Background(), tt.args.obj, tt.args.paths...)
			if tt.err != nil {
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
