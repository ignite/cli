package xtime

import "time"

// Clock represents a clock that can retrieve current time.
type Clock interface {
	Now() time.Time
	Add(duration time.Duration)
}

// ClockSystem is a clock that retrieves system time.
type ClockSystem struct{}

// NewClockSystem returns a new ClockSystem.
func NewClockSystem() ClockSystem {
	return ClockSystem{}
}

// Now implements Clock.
func (ClockSystem) Now() time.Time {
	return time.Now()
}

// Add implements Clock.
func (ClockSystem) Add(_ time.Duration) {
	panic("Add can't be called for ClockSystem")
}

// ClockMock is a clock mocking time with an internal counter.
type ClockMock struct {
	t time.Time
}

// NewClockMock returns a new ClockMock.
func NewClockMock(originalTime time.Time) *ClockMock {
	return &ClockMock{
		t: originalTime,
	}
}

// Now implements Clock.
func (c ClockMock) Now() time.Time {
	return c.t
}

// Add implements Clock.
func (c *ClockMock) Add(duration time.Duration) {
	c.t = c.t.Add(duration)
}
