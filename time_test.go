package timetype

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClock_GoString(t *testing.T) {
	s := Clock(time.Date(0, time.January, 1, 13, 24, 0, 0, time.UTC)).GoString()
	assert.Equal(t, "timetype.NewClock(13, 24, 0, UTC)", s)
}

func TestClock_String(t *testing.T) {
	s := Clock(time.Date(0, time.January, 1, 17, 54, 0, 0, time.UTC)).String()
	assert.Equal(t, "17:54:00 UTC", s)
}

func TestClock_UnmarshalJSON(t *testing.T) {
	var c Clock
	err := c.UnmarshalJSON([]byte("\"19:24:00\""))
	require.NoError(t, err)
	assert.Equal(t, Clock(time.Date(0, time.January, 1, 19, 24, 0, 0, time.UTC)), c)

	// errors
	err = c.UnmarshalJSON([]byte("19:24:00")) // time should be presented as string
	require.Error(t, err)
	assert.IsType(t, &json.SyntaxError{}, err, "time should be escaped in quotes")

	err = c.UnmarshalJSON([]byte("32145")) // invalid clock format
	assert.EqualError(t, err, "timetype: invalid clock", "clock should be in format \"15:04:05\"")
	assert.Equal(t, ErrInvalidClock, err)

	err = c.UnmarshalJSON([]byte("\"19:24:c00\""))
	require.Error(t, err)
	assert.IsType(t, &time.ParseError{}, err, "invalid character \"c\" in seconds")
}

func TestNewClock(t *testing.T) {
	assert.Equal(t, Clock(time.Date(0, time.January, 1, 13, 24, 32, 0, time.Local)),
		NewClock(13, 24, 32, time.Local))

	assert.Equal(t, Clock(time.Date(0, time.January, 1, 23, 59, 59, 0, time.UTC)),
		NewUTCClock(23, 59, 59))
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	var d Duration
	err := d.UnmarshalJSON([]byte("\"1h5m3s\""))
	require.NoError(t, err)
	assert.Equal(t, Duration(time.Hour+5*time.Minute+3*time.Second), d)

	err = d.UnmarshalJSON([]byte("3903000000000"))
	require.NoError(t, err)
	assert.Equal(t, Duration(time.Hour+5*time.Minute+3*time.Second), d)

	// errors
	err = d.UnmarshalJSON([]byte("true"))
	require.EqualError(t, err, "timetype: invalid duration", "passed bool to type duration")
	assert.Equal(t, ErrInvalidDuration, err)

	err = d.UnmarshalJSON([]byte("1h5m3s"))
	require.Error(t, err)
	assert.IsType(t, &json.SyntaxError{}, err, "duration should be escaped in quotes or passed as integer")

	err = d.UnmarshalJSON([]byte("\"123\""))
	require.Error(t, err)
	assert.EqualError(t, err, "time: missing unit in duration \"123\"", "passed empty string to time parser")
}

func TestDuration_Scan(t *testing.T) {
	tbl := []struct {
		arg      interface{}
		expected Duration
		err      string
	}{
		{
			arg:      nil,
			expected: Duration(0),
		},
		{
			arg:      5 * time.Minute,
			expected: Duration(5 * time.Minute),
		},
		{
			arg:      float64(10*time.Second + 1*time.Microsecond),
			expected: Duration(10*time.Second + 1*time.Microsecond),
		},
		{
			arg:      time.Duration(32 * time.Hour),
			expected: Duration(32 * time.Hour),
		},
		{
			arg:      `"5h3m2s"`,
			expected: Duration(5*time.Hour + 3*time.Minute + 2*time.Second),
		},
		{
			arg:      []byte(`"2h3m"`),
			expected: Duration(2*time.Hour + 3*time.Minute),
		},
		{
			arg: 'c',
			err: "timetype: invalid duration",
		},
	}
	for i, tt := range tbl {
		var d Duration
		err := d.Scan(tt.arg)
		if tt.err != "" {
			assert.EqualError(t, err, tt.err, "case #%d", i)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tt.expected, d, "case #%d", i)
	}
}

func TestDuration_Value(t *testing.T) {
	tbl := []struct {
		arg      Duration
		expected driver.Value
	}{
		{
			arg:      Duration(2*time.Hour + 3*time.Minute),
			expected: driver.Value([]byte(`"2h3m0s"`)),
		},
		{
			arg:      Duration(5*time.Hour + 3*time.Minute + 2*time.Second),
			expected: driver.Value([]byte(`"5h3m2s"`)),
		},
		{
			arg:      Duration(1 * time.Second),
			expected: driver.Value([]byte(`"1s"`)),
		},
		{
			arg:      Duration(1 * time.Millisecond),
			expected: driver.Value([]byte(`"1ms"`)),
		},
		{
			arg:      Duration(1 * time.Nanosecond),
			expected: driver.Value([]byte(`"1ns"`)),
		},
	}

	for i, tt := range tbl {
		actual, err := tt.arg.Value()
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, actual, "case #%d", i)
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	bytes, err := Duration(time.Hour + 5*time.Minute + 3*time.Second).MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`"1h5m3s"`), bytes)
}

func TestClock_MarshalJSON(t *testing.T) {
	bytes, err := Clock(time.Date(0, time.January, 1, 19, 24, 0, 0, time.UTC)).MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`"19:24:00"`), bytes)
}

func TestClock_Scan(t *testing.T) {
	tbl := []struct {
		arg      interface{}
		expected Clock
		err      string
	}{
		{
			arg:      nil,
			expected: Clock(time.Time{}),
		},
		{
			arg:      time.Date(0, time.January, 1, 2, 19, 30, 0, time.UTC),
			expected: Clock(time.Date(0, time.January, 1, 2, 19, 30, 0, time.UTC)),
		},
		{
			arg:      `"19:24:00"`,
			expected: Clock(time.Date(0, time.January, 1, 19, 24, 0, 0, time.UTC)),
		},
		{
			arg:      []byte(`"2:21:55"`),
			expected: Clock(time.Date(0, time.January, 1, 2, 21, 55, 0, time.UTC)),
		},
		{
			arg:      2567,
			expected: Clock{},
			err:      "timetype: invalid clock",
		},
	}

	for i, tt := range tbl {
		c := Clock{}
		err := c.Scan(tt.arg)
		if tt.err != "" {
			assert.EqualError(t, err, tt.err, "case #%d", i)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tt.expected, c, "case #%d", i)
	}
}

func TestClock_Value(t *testing.T) {
	tbl := []struct {
		arg      Clock
		expected driver.Value
	}{
		{
			arg:      Clock(time.Date(0, time.January, 1, 19, 24, 0, 0, time.UTC)),
			expected: driver.Value([]byte(`"19:24:00"`)),
		},
		{
			arg:      Clock(time.Date(0, time.January, 1, 2, 21, 55, 0, time.UTC)),
			expected: driver.Value([]byte(`"02:21:55"`)),
		},
		{
			arg:      Clock(time.Date(0, time.January, 1, 2, 19, 30, 0, time.UTC)),
			expected: driver.Value([]byte(`"02:19:30"`)),
		},
	}

	for i, tt := range tbl {
		actual, err := tt.arg.Value()
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, actual, "case #%d", i)
	}
}
