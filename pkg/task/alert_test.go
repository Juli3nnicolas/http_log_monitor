package task

import (
	"testing"
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/timer"

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

	ti := &timer.TimeStub{
		NowStub: func() func() time.Time {
			i := -1
			pi := &i
			t2 := time.Now()

			return func() time.Time {
				*pi++
				return t2.Add(time.Duration(*pi) * time.Second)
			}
		}(),
	}

	// Exercise & validation stages
	err := alert.Run(alertOn, ti)
	assert.Nil(t, err)
	res := alert.Result()
	assert.False(t, res.IsOn)

	err = alert.Run(alertOn, ti)
	assert.Nil(t, err)
	res = alert.Result()
	assert.True(t, res.IsOn)

	err = alert.Run(alertOff, ti)
	assert.Nil(t, err)
	res = alert.Result()
	assert.False(t, res.IsOn)
}
