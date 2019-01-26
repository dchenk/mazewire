package time

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const (
	// minSeconds is the lowest value valid for the Time.Seconds field.
	// This is time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	minSeconds = -62135596800

	// maxSeconds is the maximum valid value for the Time.Seconds field.
	// This is time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC).Unix().
	maxSeconds = 253402300799
)

// validate determines whether a Time is within the range [0001-01-01, 10000-01-01) and
// has a Nanos field within the range [0, 1e9).
// If the Time is valid, validate returns nil.
//
// Every valid Time can be represented by a time.Time, but the converse is not true.
func validate(t *Time) error {
	if t == nil {
		return errors.New("time: nil Time")
	}
	if t.Seconds < minSeconds {
		return fmt.Errorf("time: %v before 0001-01-01", t)
	}
	if t.Seconds > maxSeconds {
		return fmt.Errorf("time: %v after 9999-12-31", t)
	}
	if t.Nanos < 0 || t.Nanos >= 1e9 {
		return fmt.Errorf("time: %v: nanoseconds not in range [0, 1e9)", t)
	}
	return nil
}

// ToTime converts a Time as defined in this package to a time.Time.
// If the error returned is not nil, then the time.Time returned is invalid.
// A nil Time returns an error.
func (t *Time) ToTime() (time.Time, error) {
	if t == nil {
		return time.Time{}, validate(t)
	}
	return time.Unix(t.Seconds, int64(t.Nanos)).UTC(), validate(t)
}

// Value implements database/sql/driver.Valuer for the *Time type.
func (t *Time) Value() (driver.Value, error) {
	return t.ToTime()
}

// Scan implements database/sql.Scanner for the *Time type.
func (t *Time) Scan(src interface{}) error {
	tt, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("time: could not scan unexpected value of type %T", src)
	}
	tp, err := TimeProto(tt)
	t.Seconds = tp.Seconds
	t.Nanos = tp.Nanos
	return err
}

// TimeProto converts the time.Time to a Time as defined in this package.
// An error is returned if the resulting Time is invalid.
// If the error returned is not nil, then the Time is invalid.
func TimeProto(t time.Time) (*Time, error) {
	seconds := t.Unix()
	ts := &Time{
		Seconds: seconds,
		Nanos:   int32(t.Sub(time.Unix(seconds, 0))),
	}
	return ts, validate(ts)
}
