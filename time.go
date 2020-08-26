package timetype

import (
	"encoding/json"
	"errors"
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

// MarshalJSON marshals time into time
func (h Clock) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(h).Format(ISO8601Clock))
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
