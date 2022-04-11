package network

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want uint64
		err  error
	}{
		{
			name: "valid number",
			id:   "10",
			want: 10,
		},
		{
			name: "invalid uint",
			id:   "-10",
			err:  errors.New("error parsing launchID: strconv.ParseUint: parsing \"-10\": invalid syntax"),
		},
		{
			name: "invalid string",
			id:   "test",
			err:  errors.New("error parsing launchID: strconv.ParseUint: parsing \"test\": invalid syntax"),
		},
		{
			name: "invalid launch id",
			id:   "0",
			err:  errors.New("launch ID must be greater than 0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseID(tt.id)
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
