package task

import (
	"fmt"
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/timer"
)

// Alert is a task to alert users that a specific threshold has been exceeded
// At the moment, thresolds can only be defined on an average request-per-duration.
// Where duration is the time interval the alert is going to monitor the traffic.
// It would require a small amount of work to consider other metrics (global or
// per frame) as run already receives a task.Rates struct.
// Example : If an alert has a duration of 2 minutes and a threshold of 10 req/s,
// the alert is triggered if the average request per second during those 2 minutes
// is greater than 10 req/s.
// Conversely, the alert recovers if during 2 minutes, the traffic is below 10 req/s.
type Alert struct {
	// time when the monitoring session starts
	start time.Time
	// time when the monitoring session ends
	duration time.Duration
	// average request-per-duration
	avgReq uint64
	// number of conducted measures, used to compute avgReq
	nbMeasures uint64
	// threshold value, above the alert is triggered
	threshold uint64
	// nbReqs is the number of requests that occured during a duration period
	nbReqs float64
	// done is true if the task has finished measuring for the current frame
	done bool
	// alert's current state
	state AlertState
}

// AlertState is a struct describing an alert's state
type AlertState struct {
	// IsOn is true when the alert is active.
	// It is false if there wasn't any alert or the system recovered
	IsOn bool
	// Avg is the average req/s the alert was triggered at (always 0 if IsOn)
	Avg uint64
	// NbReqs is number of requests that triggered the alert. Always equals 0 when IsOn == false
	NbReqs uint64
	// Date is the time the alert has been switched on or off. It has a default value
	// in case the alert has never been activated.
	Date time.Time
}

// Init sets up the task, it needs in order :
// - duration time.Duration : monitoring interval, if you give it 2 minutes
// the traffic will be monitored in time-slices of 2 minutes.
// - threshold uint64 : average req/s value above which the alert is triggerd.
// The alert is in recover-state if it goes below. For an alert to be triggered,
// the threshold must be exceeded on average during "duration" time.
func (o *Alert) Init(args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("wrong parameters - the following parameters are needed (duration time.Duration, reqPSecThreshold uint64)")
	}

	var duration time.Duration
	duration, ok := args[0].(time.Duration)
	if !ok {
		return fmt.Errorf("type error - got %T instead of %T", args[0], duration)
	}

	var threshold uint64
	threshold, ok = args[1].(uint64)
	if !ok {
		return fmt.Errorf("type error - got %T instead of %T", args[1], threshold)
	}

	o.duration = duration
	o.threshold = threshold

	return nil
}

// BeforeRun inits the monitoring timer for the first time-slice
func (o *Alert) BeforeRun() error {
	o.done = false
	// If the alert monitoring hasn't started, then init the chrono
	if o.start.Unix() == (time.Time{}).Unix() {
		o.start = time.Now()
	}

	return nil
}

// Run executes the monitoring process, triggers or recovers the alert
// Parameters :
// - rates task.Rates : traffic information used to trigger the alert
// (only task.Rates.Frame.ReqPerS is used at the moment)
// - timer : a Timer interface, only uses Timer.Now()
func (o *Alert) Run(args ...interface{}) error {

	argsLen := len(args)
	if argsLen != 2 {
		return fmt.Errorf("wrong parameters -  this function expects a task.Rates parameter and a Timer struct")
	}

	var rates Rates
	rates, ok := args[0].(Rates)
	if !ok {
		return fmt.Errorf("type error - got %T instead of %T", args[0], rates)
	}

	t, ok := args[1].(timer.Timer)
	if !ok {
		return fmt.Errorf("type error - got %T instead of %T", args[1], t)
	}

	now := t.Now()

	// floats are used to be sure to work out the exact
	o.avgReq = (o.avgReq*o.nbMeasures + rates.Frame.ReqPerS) / (o.nbMeasures + 1)
	o.nbMeasures++

	// Count ongoing requests in current time-frame
	currentNbReqs := float64(rates.Frame.ReqPerS) * float64(rates.Frame.Duration)
	o.nbReqs += currentNbReqs

	if now.Sub(o.start) >= o.duration {
		if !o.state.IsOn && o.avgReq >= o.threshold {
			o.state.IsOn = true
			o.state.Date = now
			o.state.NbReqs = uint64(o.nbReqs - currentNbReqs)
			o.state.Avg = o.avgReq
		}

		if o.state.IsOn && o.avgReq < o.threshold {
			o.state.IsOn = false
			o.state.Date = now
			o.state.NbReqs = 0
			o.state.Avg = 0
		}

		// Restart a new monitoring process
		o.start = now
		o.avgReq = 0
		o.nbMeasures = 0
		o.nbReqs = 0
	}
	o.done = true

	return nil
}

// AfterRun does nothing, implements the Task interface
func (o *Alert) AfterRun() error {
	return nil
}

// Result returns a copy of the current alert state
func (o *Alert) Result() AlertState {
	return o.state
}

// IsDone always returns false as an alert is a never-ending monitoring-task
func (o *Alert) IsDone() bool {
	return false
}

// Close wipes the object's content. Call Init to use it again.
func (o *Alert) Close() error {
	*o = Alert{}
	return nil
}
