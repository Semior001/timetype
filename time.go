package timetype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Parsing errors
var (
	ErrInvalidClock    = errors.New("timetype: invalid clock")
	ErrInvalidDuration = errors.New("timetype: invalid duration")
)

// ISO8601Clock describes time layout in ISO 8601 standard
const ISO8601Clock = "15:04:05"

// Clock is a wrapper for time.time to allow parsing datetime stamp with time only in
// ISO 8601 format, like "15:04:05"
type Clock time.Time

// NewClock returns the Clock in the given location with given hours, minutes and secs
func NewClock(h, m, s int, loc *time.Location) Clock {
	return Clock(time.Date(0, time.January, 1, h, m, s, 0, loc))
}

// NewUTCClock returns new clock with given hours, minutes and seconds in the UTC location
func NewUTCClock(h, m, s int) Clock {
	return NewClock(h, m, s, time.UTC)
}

// MarshalJSON marshals time into time
func (h Clock) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(h).Format(ISO8601Clock))
}

// String implements fmt.Stringer to print and log Clock properly
func (h Clock) String() string {
	t := time.Time(h)
	return fmt.Sprintf("%02d:%02d:%02d %s", t.Hour(), t.Minute(), t.Second(), t.Location())
}

// GoString implements fmt.GoStringer to use Clock in %#v formats
func (h Clock) GoString() string {
	t := time.Time(h)
	return fmt.Sprintf("timetype.NewClock(%d, %d, %d, %s)", t.Hour(), t.Minute(), t.Second(), t.Location())
}

// UnmarshalJSON converts time to ISO 8601 representation
func (h *Clock) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	val, ok := v.(string)
	if !ok {
		return ErrInvalidClock
	}
	t, err := time.Parse(ISO8601Clock, val)
	if err != nil {
		return err
	}
	*h = Clock(t)
	return nil
}

// Scan the given SQL value as Clock
func (h *Clock) Scan(src interface{}) (err error) {
	switch v := src.(type) {
	case nil:
		*h = Clock{}
	case time.Time:
		*h = Clock(v)
	case string:
		err = h.UnmarshalJSON([]byte(v))
	case []byte:
		err = h.UnmarshalJSON(v)
	default:
		return ErrInvalidClock
	}

	return err
}

// Value returns the SQL value of the given Clock
func (h Clock) Value() (driver.Value, error) {
	return h.MarshalJSON()
}

// Duration is a wrapper of time.Duration, that allows to marshal and unmarshal time in RFC3339 format
type Duration time.Duration

// MarshalJSON simply marshals duration into nanoseconds
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON converts time duration from RFC3339 format into time.Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return ErrInvalidDuration
	}
}

// Scan the given SQL value as Duration
func (d *Duration) Scan(src interface{}) (err error) {
	switch v := src.(type) {
	case nil:
		*d = 0
	case time.Duration:
		*d = Duration(v)
	case float64:
		*d = Duration(time.Duration(v))
	case string:
		err = d.UnmarshalJSON([]byte(v))
	case []byte:
		err = d.UnmarshalJSON(v)
	default:
		return ErrInvalidDuration
	}

	return err
}

// Value returns the SQL value of the given Duration
func (d Duration) Value() (driver.Value, error) {
	return d.MarshalJSON()
}
