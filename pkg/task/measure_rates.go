package task

import (
	"fmt"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
)

// MeasureRates is a task measuring different rates and measures usefull for the whole app
type MeasureRates struct {
	done  bool
	rates Rates
}

// Rates contains all type of rates and measures taken by the task
type Rates struct {
	Global GlobalRates
	Frame  FrameRates
}

// GlobalRates global mesures taking into account the whole log file
type GlobalRates struct {
	AvgReqPerS uint64
	nbMeasures uint64
	MaxReqPerS uint64
}

// FrameRates measures related to the current time-frame
type FrameRates struct {
	Duration   uint64 // frame's duration expressed in seconds
	ReqPerS    uint64
	NbRequests uint64
	NbSuccess  uint64
	NbFailures uint64
}

// Init does nothing, implements the task interface
func (o *MeasureRates) Init(args ...interface{}) error {
	return nil
}

// BeforeRun flags the task as not done.
func (o *MeasureRates) BeforeRun(...interface{}) error {
	o.done = false

	return nil
}

// Run computes the measures present in Rates. It takes two parameters:
// logs []log.Info : slice of logs to base the computing on
// frame uint64 : time-frame duration in seconds
func (o *MeasureRates) Run(args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("wrong parameters - in order parameters must be (logs []log.Info, frameDurationInSeconds uint64)")
	}

	logs, ok := args[0].([]log.Info)
	if !ok {
		return fmt.Errorf("type error - got %T instead of []log.Info", args[0])
	}

	frame, ok := args[1].(uint64)
	if !ok {
		return fmt.Errorf("type error - got %T instead of uint64", args[1])
	}

	o.computeFrameRates(logs, frame)
	o.computeGlobalRates(logs)

	o.done = true
	return nil
}

func (o *MeasureRates) computeFrameRates(logs []log.Info, frame uint64) {
	losgLen := len(logs)
	f := &o.rates.Frame

	// All log lines always correspond to HTTP requests
	f.NbRequests = uint64(losgLen)
	f.ReqPerS = f.NbRequests / frame
	f.Duration = frame

	var nbSuccess uint64
	for i := 0; i < losgLen; i++ {
		if logs[i].Request.Code < 400 {
			nbSuccess++
		}
	}
	f.NbSuccess = nbSuccess
	f.NbFailures = f.NbRequests - nbSuccess
}

func (o *MeasureRates) computeGlobalRates(logs []log.Info) {
	f := &o.rates.Frame
	g := &o.rates.Global

	if f.ReqPerS > g.MaxReqPerS {
		g.MaxReqPerS = f.ReqPerS
	}

	if g.nbMeasures >= 1 {
		g.AvgReqPerS = (g.AvgReqPerS*g.nbMeasures + f.ReqPerS) / (g.nbMeasures + 1)
	} else {
		g.AvgReqPerS = f.ReqPerS
	}
	g.nbMeasures++
}

// AfterRun does nothing, implements the Task interface
func (o *MeasureRates) AfterRun() error {
	return nil
}

// Result returns a copy of the measure-data
func (o *MeasureRates) Result() Rates {
	return o.rates
}

// IsDone returns true if the task has complted its work. False otherwise.
func (o *MeasureRates) IsDone() bool {
	return o.done
}

// Close wipes the object's content. Call Open to use it again.
func (o *MeasureRates) Close() error {
	o.rates.Global = GlobalRates{}
	o.rates.Frame = FrameRates{}
	return nil
}
