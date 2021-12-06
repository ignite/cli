package xtime

import (
	"time"
)

// Seconds creates a time.Duration based on the seconds parameter
func Seconds(seconds uint64) time.Duration {
	return time.Duration(seconds) * time.Second
}

// NowAfter returns a unix date string from now plus the duration
func NowAfter(unix time.Duration) string {
	date := time.Now().Add(unix)
	return Format(date)
}

// Format formats the time.Time to unix date string
func Format(date time.Time) string {
	return date.Format(time.UnixDate)
}
