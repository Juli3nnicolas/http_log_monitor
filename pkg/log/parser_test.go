package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseReturnsValidStructFromTotallyStandardEntry(t *testing.T) {
	// Setup stage
	const commonLogLine string = `172.17.0.1 - - [09/Feb/2020:16:27:00 +0000] "GET / HTTP/1.1" 200 612`
	const timeFormat string = `02/Jan/2006:15:04:05 -0700`
	refTime, err := time.Parse(timeFormat, "09/Feb/2020:16:27:00 +0000")
	if err != nil {
		panic(err)
	}

	// Exercise
	info, err := Parse(commonLogLine)

	// Validation
	assert.Nil(t, err)
	assert.Equal(t, "172.17.0.1", info.Host)
	assert.Equal(t, RFC931User{}, info.LogUser)
	assert.Equal(t, "", info.AuthUser)
	assert.Equal(t, refTime, info.LocalTime)
	assert.Equal(t, refTime, info.LocalTime)
	assert.Equal(t, "/", info.Request.Route)
	assert.Equal(t, uint32(200), info.Request.Code)
	assert.Equal(t, uint64(612), info.Request.Size)
	assert.Equal(t, "HTTP/1.1", info.Request.Version)
}

func TestParseReturnsValidStructFromStandardEntryWithAdditionalFields(t *testing.T) {
	// Setup stage
	const logEntry string = `172.17.0.1 - - [09/Feb/2020:16:27:00 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.54.0" "-"`
	const timeFormat string = `02/Jan/2006:15:04:05 -0700`
	refTime, err := time.Parse(timeFormat, "09/Feb/2020:16:27:00 +0000")
	if err != nil {
		panic(err)
	}

	// Exercise
	info, err := Parse(logEntry)

	// Validation
	assert.Nil(t, err)
	assert.Equal(t, "172.17.0.1", info.Host)
	assert.Equal(t, RFC931User{}, info.LogUser)
	assert.Equal(t, "", info.AuthUser)
	assert.Equal(t, refTime, info.LocalTime)
	assert.Equal(t, refTime, info.LocalTime)
	assert.Equal(t, "/", info.Request.Route)
	assert.Equal(t, uint32(200), info.Request.Code)
	assert.Equal(t, uint64(612), info.Request.Size)
	assert.Equal(t, "HTTP/1.1", info.Request.Version)
}

func TestParseReturnsAnErrorIfTheInputStringIsEmpty(t *testing.T) {
	_, err := Parse("")
	assert.NotNil(t, err)
}
