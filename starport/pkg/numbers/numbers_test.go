package numbers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseList(t *testing.T) {
	cases := []struct {
		list   string
		parsed []int
	}{
		{"1,2,3", []int{1, 2, 3}},
		{"1, 2,3 ", []int{1, 2, 3}},
		{",1, 2,", []int{1, 2}},
		{",", []int{}},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("no: %d", i), func(t *testing.T) {
			parsed, err := ParseList(tt.list)
			require.NoError(t, err)
			require.Equal(t, tt.parsed, parsed)
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		parsed []int
		list   string
	}{
		{[]int{1, 2, 3}, "#1, #2, #3"},
		{[]int{1}, "#1"},
		{[]int{}, ""},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("no: %d", i), func(t *testing.T) {
			require.Equal(t, tt.list, List(tt.parsed, "#"))
		})
	}
}
