package prefixgen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGen(t *testing.T) {
	cases := []struct {
		expected string
		given    string
	}{
		{"[TENDERMINT] ", New("Tendermint", Common()...).Gen()},
		{"Tendermint", New("Tendermint").Gen()},
		{"appd", New("%sd").Gen("app")},
	}
	for _, tt := range cases {
		t.Run(tt.expected, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.given)
		})
	}
}
