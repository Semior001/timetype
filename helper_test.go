package timetype

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWeekday(t *testing.T) {
	mustParse := func(s string) time.Weekday {
		wd, err := ParseWeekday(s)
		require.NoError(t, err)
		return wd
	}
	assert.Equal(t, time.Sunday, mustParse("Sunday"))
	assert.Equal(t, time.Monday, mustParse("Monday"))
	assert.Equal(t, time.Tuesday, mustParse("Tuesday"))
	assert.Equal(t, time.Wednesday, mustParse("Wednesday"))
	assert.Equal(t, time.Thursday, mustParse("Thursday"))
	assert.Equal(t, time.Friday, mustParse("Friday"))
	assert.Equal(t, time.Saturday, mustParse("Saturday"))

	_, err := ParseWeekday("Workday")
	assert.EqualError(t, err, "timetype: invalid weekday")
	assert.Equal(t, ErrInvalidWeekday, err)
}
