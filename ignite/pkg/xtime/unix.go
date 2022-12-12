package xtime

import (
	"time"
)

// Seconds creates a time.Duration based on the seconds parameter.
func Seconds(seconds int64) time.Duration {
	return time.Duration(seconds) * time.Second
}

// NowAfter returns a unix date string from now plus the duration.
func NowAfter(unix time.Duration) string {
	date := time.Now().Add(unix)
	return FormatUnix(date)
}

// FormatUnix formats the time.Time to unix date string.
func FormatUnix(date time.Time) string {
	return date.Format(time.UnixDate)
}

// FormatUnixInt formats the int timestamp to unix date string.
func FormatUnixInt(unix int64) string {
	return FormatUnix(time.Unix(unix, 0))
}
