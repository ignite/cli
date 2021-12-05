package date

import (
	"time"
)

func Now(unix uint64) string {
	date := time.Now().Add(time.Duration(unix) * time.Second)
	return ToString(date)
}

func ToString(date time.Time) string {
	return date.Format(time.UnixDate)
}
