package date

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNow(t *testing.T) {
	tests := []uint64{
		9999999999,
		10000,
		100,
		0,
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("test %d value", tt), func(t *testing.T) {
			got := Now(tt)
			date := time.Now().Add(time.Duration(tt) * time.Second)
			require.Equal(t, date.Format(time.UnixDate), got)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		date time.Time
		want string
	}{
		{
			date: time.Time{},
			want: "Mon Jan  1 00:00:00 UTC 0001",
		},
		{
			date: time.Unix(10000000000, 100),
			want: "Sat Nov 20 14:46:40 -03 2286",
		},
		{
			date: time.Date(2020, 10, 11, 12, 30, 50, 0, time.FixedZone("Europe/Berlin", 3*60*60)),
			want: "Sun Oct 11 12:30:50 Europe/Berlin 2020",
		},
	}
	for _, tt := range tests {
		t.Run("test date "+tt.date.String(), func(t *testing.T) {
			got := ToString(tt.date)
			require.Equal(t, tt.want, got)
		})
	}
}
