package testutil

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// AssertTimeEqual asserts that two times are equal within a tolerance
func AssertTimeEqual(t *testing.T, expected, actual time.Time, tolerance time.Duration, msgAndArgs ...interface{}) bool {
	t.Helper()

	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}

	return assert.True(t, diff <= tolerance, msgAndArgs...)
}

// AssertTimeNow asserts that a time is close to now
func AssertTimeNow(t *testing.T, actual time.Time, tolerance time.Duration, msgAndArgs ...interface{}) bool {
	t.Helper()
	return AssertTimeEqual(t, time.Now().UTC(), actual, tolerance, msgAndArgs...)
}

// AssertUUIDValid asserts that a UUID is valid
func AssertUUIDValid(t *testing.T, id uuid.UUID, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.NotEqual(t, uuid.Nil, id, msgAndArgs...)
}

// AssertUUIDNotNil asserts that a UUID is not nil
func AssertUUIDNotNil(t *testing.T, id uuid.UUID, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.NotEqual(t, uuid.Nil, id, msgAndArgs...)
}

// AssertNoError is a helper that calls require.NoError with better error messages
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.NoError(t, err, msgAndArgs...)
}

// AssertError is a helper that calls require.Error
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.Error(t, err, msgAndArgs...)
}

// AssertErrorContains asserts that an error contains a specific message
func AssertErrorContains(t *testing.T, err error, contains string, msgAndArgs ...interface{}) bool {
	t.Helper()

	if !assert.Error(t, err, msgAndArgs...) {
		return false
	}

	return assert.Contains(t, err.Error(), contains, msgAndArgs...)
}

// AssertJSONField asserts that a JSON map contains a specific field
func AssertJSONField(t *testing.T, jsonMap map[string]interface{}, field string, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.Contains(t, jsonMap, field, msgAndArgs...)
}

// AssertJSONFieldValue asserts that a JSON map field has a specific value
func AssertJSONFieldValue(t *testing.T, jsonMap map[string]interface{}, field string, expected interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if !AssertJSONField(t, jsonMap, field, msgAndArgs...) {
		return false
	}

	return assert.Equal(t, expected, jsonMap[field], msgAndArgs...)
}

// AssertSliceNotEmpty asserts that a slice is not empty
func AssertSliceNotEmpty(t *testing.T, slice interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.NotEmpty(t, slice, msgAndArgs...)
}

// AssertSliceLength asserts that a slice has a specific length
func AssertSliceLength(t *testing.T, slice interface{}, length int, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.Len(t, slice, length, msgAndArgs...)
}

// AssertMapContainsKey asserts that a map contains a specific key
func AssertMapContainsKey(t *testing.T, m map[string]interface{}, key string, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.Contains(t, m, key, msgAndArgs...)
}

// AssertStringNotEmpty asserts that a string is not empty
func AssertStringNotEmpty(t *testing.T, s string, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.NotEmpty(t, s, msgAndArgs...)
}

// AssertPointerNotNil asserts that a pointer is not nil
func AssertPointerNotNil(t *testing.T, ptr interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	return assert.NotNil(t, ptr, msgAndArgs...)
}

