package numbers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseList(t *testing.T) {
	cases := []struct {
		list   string
		parsed []uint64
	}{
		{"1,2,3", []uint64{1, 2, 3}},
		{"1, 2,3 ", []uint64{1, 2, 3}},
		{",1, 2,", []uint64{1, 2}},
		{"1-3 ", []uint64{1, 2, 3}},
		{"1-3,8 ", []uint64{1, 2, 3, 8}},
		{"1-3,8-11 ", []uint64{1, 2, 3, 8, 9, 10, 11}},
		{"1-3,8-11,33 ", []uint64{1, 2, 3, 8, 9, 10, 11, 33}},
		{"1-3,8-11,33-36 ", []uint64{1, 2, 3, 8, 9, 10, 11, 33, 34, 35, 36}},
		{"1-5,2-7,9-11,1-8 ", []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
		{",", []uint64{}},
		{",-", []uint64{}},
		{",10-", []uint64{10}},
		{"10-", []uint64{10}},
		{"-10", []uint64{10}},
		{"10-10", []uint64{10}},
		{"12-8", []uint64{8, 9, 10, 11, 12}},
		{"12-8,4-1", []uint64{1, 2, 3, 4, 8, 9, 10, 11, 12}},
	}
	for _, tt := range cases {
		t.Run("list "+tt.list, func(t *testing.T) {
			parsed, err := ParseList(tt.list)
			require.NoError(t, err)
			require.Equal(t, tt.parsed, parsed)
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		parsed []uint64
		list   string
	}{
		{[]uint64{1, 2, 3}, "#1, #2, #3"},
		{[]uint64{1}, "#1"},
		{[]uint64{}, ""},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("no: %d", i), func(t *testing.T) {
			require.Equal(t, tt.list, List(tt.parsed, "#"))
		})
	}
}
