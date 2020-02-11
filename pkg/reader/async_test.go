package reader

import (
	"testing"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestOpenInitsBuffer(t *testing.T) {
	// Setup stage
	reader := Stub{OpenStub: func(...interface{}) error { return nil }}
	ar := Async{Reader: &reader}

	// Exercise stage
	err := ar.Open(nil)

	// Validation stage
	assert.Nil(t, err)
	assert.Equal(t, defaultMinBufSize, ar.MinBufsize)
	assert.Equal(t, cap(ar.buffer), int(defaultMinBufSize))
}

func TestOpenInitsBufferWithCustomMinBufSize(t *testing.T) {
	// Setup stage
	reader := Stub{OpenStub: func(...interface{}) error { return nil }}
	ar := Async{Reader: &reader}

	const customMinSize uint64 = 800
	ar.MinBufsize = customMinSize

	// Exercise stage
	err := ar.Open(nil)

	// Validation stage
	assert.Nil(t, err)
	assert.Equal(t, customMinSize, ar.MinBufsize)
	assert.Equal(t, cap(ar.buffer), int(customMinSize))
}

func TestStartWritesToBuffer(t *testing.T) {
	// Setup stage
	const host1 string = "host1"
	const host2 string = "host2"
	sharedCount := 0
	callCount := &sharedCount
	done := make(chan bool)

	reader := Stub{
		OpenStub: func(...interface{}) error { return nil },
		ReadStub: func() ([]log.Info, error) {
			if *callCount == 0 {
				*callCount++
				return []log.Info{log.Info{Host: host1}, log.Info{Host: host2}}, nil
			}
			if *callCount == 1 {
				*callCount++
				done <- true
				close(done)
			}
			return nil, nil
		},
	}
	ar := Async{Reader: &reader}

	err := ar.Open(nil)
	assert.Nil(t, err)

	// Exercise stage
	ar.Start()
	<-done
	ar.Stop() // otherwise it wouldn't stop

	// Validation stage
	assert.Len(t, ar.buffer, 2)
	assert.Equal(t, host1, ar.buffer[0].Host)
	assert.Equal(t, host2, ar.buffer[1].Host)
}
