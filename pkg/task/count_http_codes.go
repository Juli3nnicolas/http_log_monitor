package task

import (
	"fmt"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
)

// CountHTTPCodes is a task that counts all different codes present in a []log.Info
type CountHTTPCodes struct {
	done  bool
	codes map[uint32]uint64
}

// Init does nothing, implements the task interface
func (o *CountHTTPCodes) Init(args ...interface{}) error {
	return nil
}

// BeforeRun flags the task as not done.
func (o *CountHTTPCodes) BeforeRun(...interface{}) error {
	o.done = false
	o.codes = make(map[uint32]uint64)

	return nil
}

// Run parses the log slice to count all occuring codes
// logs []log.Info : slice of logs to base the computing on
func (o *CountHTTPCodes) Run(args ...interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong parameters - the only required parameter is (logs []log.Info)")
	}

	logs, ok := args[0].([]log.Info)
	if !ok {
		return fmt.Errorf("type error - got %T instead of []log.Info", args[0])
	}

	logsLen := len(logs)
	for i := 0; i < logsLen; i++ {
		o.codes[logs[i].Request.Code]++
	}

	o.done = true
	return nil
}

// AfterRun does nothing, implements the Task interface
func (o *CountHTTPCodes) AfterRun() error {
	return nil
}

// Result returns a map of error codes (the key) and their number of occurence (the value)
// Exemple : nb404 := Result()[404] // gets the number of 404 present in the log
func (o *CountHTTPCodes) Result() map[uint32]uint64 {
	return o.codes
}

// IsDone returns true if the task has complted its work. False otherwise.
func (o *CountHTTPCodes) IsDone() bool {
	return o.done
}

// Close wipes the object's content. Call Open to use it again.
func (o *CountHTTPCodes) Close() error {
	o.codes = nil
	return nil
}
