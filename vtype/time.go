package vtype

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Time is a wrapper around time.Time that provides UTC-based time handling.
// All time values are stored and handled in UTC to ensure consistency.
type Time struct {
	time.Time
}

// NewTime creates a new Time value from a time.Time pointer.
// If t is nil, returns a zero Time value.
// The input time is converted to UTC.
func NewTime(t *time.Time) *Time {
	if t == nil {
		return &Time{time.Time{}}
	}
	return &Time{t.UTC()}
}

// NewTimeFromUnix creates a new Time value from a Unix timestamp (seconds since epoch).
// If u is 0, returns a zero Time value.
func NewTimeFromUnix(u int64) *Time {
	if u == 0 {
		return &Time{time.Time{}}
	}
	tp := time.Unix(u, 0).UTC()
	return NewTime(&tp)
}

// NewTimeFromUnixMilli creates a new Time value from a Unix timestamp in milliseconds.
// If m is 0, returns a zero Time value.
func NewTimeFromUnixMilli(m int64) *Time {
	if m == 0 {
		return &Time{time.Time{}}
	}
	tp := time.UnixMilli(m).UTC()
	return NewTime(&tp)
}

// NewTimeFromUnixNano creates a new Time value from a Unix timestamp in nanoseconds.
// If n is 0, returns a zero Time value.
func NewTimeFromUnixNano(n int64) *Time {
	if n == 0 {
		return &Time{time.Time{}}
	}
	tp := time.Unix(0, n).UTC()
	return NewTime(&tp)
}

// NewTimeFromString parses a date string using the given format and returns a Time value.
// The format is the same as time.Parse format.
// Returns an error if the parsing fails.
func NewTimeFromString(dateStr string, format string) (*Time, error) {
	t, err := time.Parse(format, dateStr)
	if err != nil {
		return nil, err
	}
	return NewTime(&t), nil
}

// NewTimeFromStringWithTimeZone parses a date string using the given format and location.
// The location specifies the time zone for parsing.
// Returns an error if the parsing fails.
func NewTimeFromStringWithTimeZone(dateStr string, format string, location *time.Location) (*Time, error) {
	t, err := time.ParseInLocation(format, dateStr, location)
	if err != nil {
		return nil, err
	}
	return NewTime(&t), nil
}

// NewTimeNow returns the current time in UTC.
func NewTimeNow() *Time {
	now := time.Now().UTC()
	return &Time{now}
}

// Unix returns the Unix timestamp in seconds.
func (t *Time) Unix() int64 {
	return t.Time.Unix()
}

// UnixMilli returns the Unix timestamp in milliseconds.
// Returns 0 if the time is zero.
func (t *Time) UnixMilli() int64 {
	if t.Time.Equal(time.Time{}) {
		return 0
	}
	return t.Time.UnixMilli()
}

// UnixNano returns the Unix timestamp in nanoseconds.
// Returns 0 if the time is zero.
func (t *Time) UnixNano() int64 {
	if t.Time.Equal(time.Time{}) {
		return 0
	}
	return t.Time.UnixNano()
}

// PbTimestamp converts the Time to a protobuf Timestamp.
// Returns nil if the time is zero.
func (t *Time) PbTimestamp() *timestamppb.Timestamp {
	if t.Time.Equal(time.Time{}) {
		return nil
	}
	tp := t.Time.UTC()
	return timestamppb.New(tp)
}
