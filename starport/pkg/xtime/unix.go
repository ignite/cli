package xtime

import (
	"time"
)

func NowAfter(unix uint64) string {
	date := time.Now().Add(time.Duration(unix) * time.Second)
	return Format(date)
}

func Format(date time.Time) string {
	return date.Format(time.UnixDate)
}
