package task

import (
	"testing"
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/timer"

	"github.com/stretchr/testify/assert"
)

func TestRunTurnTheAlertOnAndOffCorrectly(t *testing.T) {
	// Setup stage
	const frameDuration time.Duration = time.Second

	alert := Alert{}
	if err := alert.Init(time.Second, uint64(4)); err != nil {
		panic(err)
	}

	// Sets alert.start to time.Now()
	if err := alert.BeforeRun(); err != nil {
		panic(err)
	}

	alertOn := Rates{
		Frame: FrameRates{
			Duration: uint64(frameDuration.Seconds()),
			ReqPerS:  5,
		},
	}
	alertOff := Rates{
		Frame: FrameRates{
			ReqPerS: 2,
		},
	}

	now := time.Now()
	t1 := &timer.TimeStub{
		NowStub: func() func() time.Time {
			i := time.Duration(0)
			iptr := &i
			pnow := &now

			// Returns now at first call then always now + 1 sec
			return func() time.Time {
				*pnow = pnow.Add(*iptr * frameDuration)
				*iptr = 1
				return *pnow
			}
		}(),
	}

	// Exercise & validation stages

	// 5 requests at t = 0s, the alert is not activated as the
	// average req/s must be above 4 req/s for a second
	err := alert.Run(alertOn, t1)
	assert.Nil(t, err)
	res := alert.Result()
	assert.False(t, res.IsOn)
	assert.Equal(t, uint64(0), res.Avg)
	assert.Equal(t, uint64(0), res.NbReqs)
	assert.Equal(t, time.Time{}, res.Date)

	// A second has elapsed, with 5 other requests
	// The average on a second is therefore of 5 req/s > 4 req/s
	// So the alert is switched on
	err = alert.Run(alertOn, t1)
	assert.Nil(t, err)
	res = alert.Result()
	tAlertOn := now
	assert.True(t, res.IsOn)
	assert.Equal(t, alertOn.Frame.ReqPerS, res.Avg)
	assert.Equal(t, alertOn.Frame.ReqPerS, res.NbReqs)
	assert.Equal(t, tAlertOn, res.Date)

	// Another second elapse, now the request rate dived to 2 req/s
	// 2 req/s < 4 req/s so the alert is switched off
	err = alert.Run(alertOff, t1)
	assert.Nil(t, err)
	res = alert.Result()
	tAlertOff := now
	assert.False(t, res.IsOn)
	assert.Equal(t, uint64(0), res.Avg)
	assert.Equal(t, uint64(0), res.NbReqs)
	assert.Equal(t, tAlertOff, res.Date)
}

func TestRunTurnTheAlertOnAndOffCorrectlyWithMeasuresSpanningOnTwoTimeFrames(t *testing.T) {
	// Setup stage
	const duration time.Duration = 2 * time.Minute
	const frameDuration time.Duration = time.Minute

	alert := Alert{}
	if err := alert.Init(duration, uint64(10)); err != nil {
		panic(err)
	}

	// Sets alert.start to time.Now()
	if err := alert.BeforeRun(); err != nil {
		panic(err)
	}

	alertOn := Rates{
		Frame: FrameRates{
			Duration: uint64(frameDuration.Seconds()),
			ReqPerS:  15,
		},
	}
	alertOff := Rates{
		Frame: FrameRates{
			ReqPerS: 2,
		},
	}

	now := time.Now()
	ti := &timer.TimeStub{
		NowStub: func() func() time.Time {
			i := time.Duration(0)
			iptr := &i
			pnow := &now

			// Returns now at first call then always now + 1 sec
			return func() time.Time {
				*pnow = pnow.Add(*iptr * frameDuration)
				*iptr = 1
				return *pnow
			}
		}(),
	}

	// Exercise & validation stages

	// 15 requests/s at t = 0s, the alert is not activated as the
	// average req/s must be above 10 req/s for 2 minutes
	err := alert.Run(alertOn, ti)
	assert.Nil(t, err)
	res := alert.Result()
	assert.False(t, res.IsOn)
	assert.Equal(t, uint64(0), res.Avg)
	assert.Equal(t, uint64(0), res.NbReqs)
	assert.Equal(t, time.Time{}, res.Date)

	// A minute has elapsed, with 15 other requests/s
	// The two minutes haven't elapsed so the alert is not switched on
	err = alert.Run(alertOn, ti)
	assert.Nil(t, err)
	res = alert.Result()
	assert.False(t, res.IsOn)
	assert.Equal(t, uint64(0), res.Avg)
	assert.Equal(t, uint64(0), res.NbReqs)
	assert.Equal(t, time.Time{}, res.Date)

	// Another minute has elapsed now totalling 2 minutes with 15 other requests/s
	// The average on 2 minutes is therefore of 15 req/s > 10 req/s
	// So the alert is switched on
	err = alert.Run(alertOn, ti)
	assert.Nil(t, err)
	res = alert.Result()
	tAlertOn := now
	assert.True(t, res.IsOn)
	assert.Equal(t, alertOn.Frame.ReqPerS, res.Avg)
	assert.Equal(t, int(alertOn.Frame.ReqPerS*uint64(duration.Seconds())), int(res.NbReqs))
	assert.Equal(t, tAlertOn, res.Date)

	// A minute has elapsed in the new time-frame with request rate diving to 2 req/s
	// However, one minute still remains so the alert is still on
	err = alert.Run(alertOff, ti)
	assert.Nil(t, err)
	res = alert.Result()
	assert.True(t, res.IsOn)
	assert.Equal(t, alertOn.Frame.ReqPerS, res.Avg)
	assert.Equal(t, int(alertOn.Frame.ReqPerS*uint64(duration.Seconds())), int(res.NbReqs))
	//assert.Equal(t, int(alertOn.Frame.ReqPerS*uint64(duration.Seconds())), int(res.NbReqs))
	assert.Equal(t, tAlertOn, res.Date)

	// The missing minute has passed with another rate of 2 req/s
	// the average being 2 req/s < 4 req/s so the alert is switched off
	err = alert.Run(alertOff, ti)
	assert.Nil(t, err)
	res = alert.Result()
	tAlertOff := now
	assert.False(t, res.IsOn)
	assert.Equal(t, 0, int(res.Avg))
	assert.Equal(t, 0, int(res.NbReqs))
	assert.Equal(t, tAlertOff, res.Date)
}
