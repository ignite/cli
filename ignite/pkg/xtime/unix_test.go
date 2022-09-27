package xtime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ignite/cli/ignite/pkg/xtime"

	"github.com/stretchr/testify/require"
)

func TestSeconds(t *testing.T) {
	tests := []int64{
		9999999999,
		10000,
		100,
		0,
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("test %d value", tt), func(t *testing.T) {
			got := xtime.Seconds(tt)
			require.Equal(t, time.Duration(tt)*time.Second, got)
		})
	}
}

func TestNowAfter(t *testing.T) {
	tests := []int64{
		9999999999,
		10000,
		100,
		0,
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("test %d value", tt), func(t *testing.T) {
			got := xtime.NowAfter(xtime.Seconds(tt))
			date := time.Now().Add(time.Duration(tt) * time.Second)
			require.Equal(t, date.Format(time.UnixDate), got)
		})
	}
}

func TestFormatUnix(t *testing.T) {
	tests := []struct {
		date time.Time
		want string
	}{
		{
			date: time.Time{},
			want: "Mon Jan  1 00:00:00 UTC 0001",
		},
		{
			date: time.Unix(10000000000, 100).In(time.UTC),
			want: "Sat Nov 20 17:46:40 UTC 2286",
		},
		{
			date: time.Date(2020, 10, 11, 12, 30, 50, 0, time.FixedZone("Europe/Berlin", 3*60*60)),
			want: "Sun Oct 11 12:30:50 Europe/Berlin 2020",
		},
	}
	for _, tt := range tests {
		t.Run("test date "+tt.date.String(), func(t *testing.T) {
			got := xtime.FormatUnix(tt.date)
			require.Equal(t, tt.want, got)
		})
	}
}
