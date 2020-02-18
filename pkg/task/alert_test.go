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
			pnow := &now

			// Returns now + 1 sec
			return func() time.Time {
				*pnow = pnow.Add(frameDuration)
				return *pnow
			}
		}(),
	}

	// Exercise & validation stages

	// A second has elapsed, with 5 requests
	// The average on a second is therefore of 5 req/s > 4 req/s
	// So the alert is switched on
	err := alert.Run(alertOn, t1)
	assert.Nil(t, err)
	res := alert.Result()
	tAlertOn := now
	assert.True(t, res.IsOn)
	assert.Equal(t, alertOn.Frame.ReqPerS, res.Avg)
	assert.Equal(t, alertOn.Frame.ReqPerS, res.NbReqs)
	assert.Equal(t, tAlertOn, res.Date)

	// Another second elapses, now the request rate dived to 2 req/s
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
			pnow := &now

			// Returns now + 1 sec
			return func() time.Time {
				*pnow = pnow.Add(frameDuration)
				return *pnow
			}
		}(),
	}

	// Exercise & validation stages

	// A minute has elapsed, with 15 other requests/s
	// The two minutes haven't elapsed so the alert is not switched on
	err := alert.Run(alertOn, ti)
	assert.Nil(t, err)
	res := alert.Result()
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

func TestRunReturnsFalseForIsOnIfThresholdIsNotReachedButTimeLimitIs(t *testing.T) {
	// Setup stage
	const frameDuration time.Duration = time.Minute

	alert := Alert{}
	if err := alert.Init(frameDuration, uint64(4)); err != nil {
		panic(err)
	}

	// Sets alert.start to time.Now()
	if err := alert.BeforeRun(); err != nil {
		panic(err)
	}

	alertOff := Rates{
		Frame: FrameRates{
			ReqPerS: 2,
		},
	}

	now := time.Now()
	t1 := &timer.TimeStub{
		NowStub: func() func() time.Time {
			pnow := &now

			// Returns now + 1 sec
			return func() time.Time {
				*pnow = pnow.Add(frameDuration)
				return *pnow
			}
		}(),
	}

	// Exercise & validation stages

	// Enough time spent checking but the rate is not high enough to trigger the alert
	err := alert.Run(alertOff, t1)
	assert.Nil(t, err)
	res := alert.Result()
	assert.False(t, res.IsOn)
	assert.Equal(t, uint64(0), res.Avg)
	assert.Equal(t, uint64(0), res.NbReqs)
	assert.Equal(t, time.Time{}, res.Date)
}

func TestRunReturnsFalseForIsOnIfThresholdIsReachedButTimeLimitIsNot(t *testing.T) {
	// Setup stage
	const frameDuration time.Duration = time.Minute

	alert := Alert{}
	if err := alert.Init(frameDuration, uint64(4)); err != nil {
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

	now := time.Now()
	t1 := &timer.TimeStub{
		NowStub: func() func() time.Time {
			pnow := &now

			// Returns now + 1 sec
			return func() time.Time {
				*pnow = pnow.Add(frameDuration / 2)
				return *pnow
			}
		}(),
	}

	// Exercise & validation stages

	// The rate is high enough but too little spent checking so the alert is not triggered
	err := alert.Run(alertOn, t1)
	assert.Nil(t, err)
	res := alert.Result()
	assert.False(t, res.IsOn)
	assert.Equal(t, uint64(0), res.Avg)
	assert.Equal(t, uint64(0), res.NbReqs)
	assert.Equal(t, time.Time{}, res.Date)
}
