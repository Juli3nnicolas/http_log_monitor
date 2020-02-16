package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunTurnTheAlertOnAndOffCorrectly(t *testing.T) {
	// Setup stage
	alert := Alert{}
	if err := alert.Init(time.Second, uint64(4)); err != nil {
		panic(err)
	}

	if err := alert.BeforeRun(); err != nil {
		panic(err)
	}

	alertOn := Rates{
		Frame: FrameRates{
			ReqPerS: 5,
		},
	}
	alertOff := Rates{
		Frame: FrameRates{
			ReqPerS: 2,
		},
	}

	// Exercise & validation stages
	err := alert.Run(alertOn)
	assert.Nil(t, err)
	res := alert.Result()
	assert.False(t, res.IsOn)

	time.Sleep(time.Second)

	err = alert.Run(alertOn)
	assert.Nil(t, err)
	res = alert.Result()
	assert.True(t, res.IsOn)

	time.Sleep(time.Second)

	err = alert.Run(alertOff)
	assert.Nil(t, err)
	res = alert.Result()
	assert.False(t, res.IsOn)
}
