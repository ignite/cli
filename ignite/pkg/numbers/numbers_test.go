package numbers

import (
	"errors"
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
		{"8-11, 1-3, ", []uint64{8, 9, 10, 11, 1, 2, 3}},
		{"1-3,33, 8-11, ", []uint64{1, 2, 3, 33, 8, 9, 10, 11}},
		{"1-3,8-11,33-36 ", []uint64{1, 2, 3, 8, 9, 10, 11, 33, 34, 35, 36}},
		{"2-7,2-5,9-11,1-8", []uint64{2, 3, 4, 5, 6, 7, 9, 10, 11, 1, 8}},
		{",", []uint64{}},
		{",-", []uint64{}},
		{",10-", []uint64{10}},
		{"10-", []uint64{10}},
		{"-10", []uint64{10}},
		{"10-10", []uint64{10}},
	}
	for _, tt := range cases {
		t.Run("list "+tt.list, func(t *testing.T) {
			parsed, err := ParseList(tt.list)
			require.NoError(t, err)
			require.Equal(t, tt.parsed, parsed)
		})
	}
}

func TestParseListErrors(t *testing.T) {
	cases := []struct {
		list string
		err  error
	}{
		{"12-8", errors.New("cannot parse a reverse ordering range: 12-8")},
		{"1-2-3", errors.New("cannot parse the number range: 1-2-3")},
	}
	for _, tt := range cases {
		t.Run("list "+tt.list, func(t *testing.T) {
			_, err := ParseList(tt.list)
			require.Error(t, err)
			require.Equal(t, tt.err, err)
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
