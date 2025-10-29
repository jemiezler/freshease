package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPtrIfNotNil(t *testing.T) {
	t.Run("returns nil when input is nil", func(t *testing.T) {
		var s *string
		result := PtrIfNotNil(s)
		assert.Nil(t, result)
	})

	t.Run("returns pointer to string when input is not nil", func(t *testing.T) {
		input := "test string"
		result := PtrIfNotNil(&input)
		assert.NotNil(t, result)
		assert.Equal(t, "test string", *result)
	})

	t.Run("returns pointer to empty string", func(t *testing.T) {
		input := ""
		result := PtrIfNotNil(&input)
		assert.NotNil(t, result)
		assert.Equal(t, "", *result)
	})

	t.Run("handles special characters", func(t *testing.T) {
		input := "test@#$%^&*()_+-=[]{}|;':\",./<>?"
		result := PtrIfNotNil(&input)
		assert.NotNil(t, result)
		assert.Equal(t, "test@#$%^&*()_+-=[]{}|;':\",./<>?", *result)
	})

	t.Run("handles unicode characters", func(t *testing.T) {
		input := "测试字符串 🚀 émojis"
		result := PtrIfNotNil(&input)
		assert.NotNil(t, result)
		assert.Equal(t, "测试字符串 🚀 émojis", *result)
	})
}

func TestTimeToISOString(t *testing.T) {
	t.Run("returns nil when input is nil", func(t *testing.T) {
		var time *time.Time
		result := TimeToISOString(time)
		assert.Nil(t, result)
	})

	t.Run("converts time to ISO string", func(t *testing.T) {
		// Create a specific time for testing
		testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		result := TimeToISOString(&testTime)

		assert.NotNil(t, result)
		assert.Equal(t, "2023-12-25T15:30:45Z", *result)
	})

	t.Run("handles different time zones", func(t *testing.T) {
		// Test with different timezone
		loc, err := time.LoadLocation("America/New_York")
		require.NoError(t, err)

		testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, loc)
		result := TimeToISOString(&testTime)

		assert.NotNil(t, result)
		// The result should be in RFC3339 format
		assert.Contains(t, *result, "2023-12-25T")
		assert.Contains(t, *result, "15:30:45")
	})

	t.Run("handles zero time", func(t *testing.T) {
		zeroTime := time.Time{}
		result := TimeToISOString(&zeroTime)

		assert.NotNil(t, result)
		assert.Equal(t, "0001-01-01T00:00:00Z", *result)
	})

	t.Run("handles current time", func(t *testing.T) {
		now := time.Now()
		result := TimeToISOString(&now)

		assert.NotNil(t, result)
		// Verify it's a valid RFC3339 format
		parsedTime, err := time.Parse(time.RFC3339, *result)
		assert.NoError(t, err)
		// Compare with some tolerance for time differences
		assert.True(t, parsedTime.UTC().Sub(now.UTC()).Abs() < time.Second)
	})

	t.Run("handles time with nanoseconds", func(t *testing.T) {
		testTime := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
		result := TimeToISOString(&testTime)

		assert.NotNil(t, result)
		// RFC3339 format may not include nanoseconds, so just check it contains the time
		assert.Contains(t, *result, "2023-12-25T15:30:45")
		assert.Contains(t, *result, "Z")
	})
}

func TestHelpersIntegration(t *testing.T) {
	t.Run("combined usage with real data", func(t *testing.T) {
		// Simulate a scenario where we have a time and want to convert it
		now := time.Now()
		timeStr := TimeToISOString(&now)

		// Then use PtrIfNotNil on the result
		result := PtrIfNotNil(timeStr)

		assert.NotNil(t, result)
		assert.NotEmpty(t, *result)

		// Verify it's a valid time string
		parsedTime, err := time.Parse(time.RFC3339, *result)
		assert.NoError(t, err)
		// Compare with some tolerance for time differences
		assert.True(t, parsedTime.UTC().Sub(now.UTC()).Abs() < time.Second)
	})

	t.Run("edge case with nil time and string", func(t *testing.T) {
		var nilTime *time.Time
		var nilString *string

		timeResult := TimeToISOString(nilTime)
		stringResult := PtrIfNotNil(nilString)

		assert.Nil(t, timeResult)
		assert.Nil(t, stringResult)
	})
}
