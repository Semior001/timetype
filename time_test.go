package timetype

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.Equal(t, Clock(time.Date(0, 0, 0, 13, 24, 32, 0, time.Local)),
		NewClock(13, 24, 32, time.Local))

	assert.Equal(t, Clock(time.Date(0, 0, 0, 23, 59, 59, 0, time.UTC)),
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

	err = d.UnmarshalJSON([]byte("\"\""))
	require.Error(t, err)
	assert.EqualError(t, err, "time: invalid duration ", "passed empty string to time parser")
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
