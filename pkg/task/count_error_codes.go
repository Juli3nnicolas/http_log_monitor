package task

import (
	"fmt"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
)

// CountErrorCodes is a task that counts all different error codes present in a []log.Info
// TODO: rename this task, the name is not accurate
type CountErrorCodes struct {
	done  bool
	codes map[uint32]uint64
}

// Init does nothing, implements the task interface
func (o *CountErrorCodes) Init(args ...interface{}) error {
	return nil
}

// BeforeRun flags the task as not done.
func (o *CountErrorCodes) BeforeRun(...interface{}) error {
	o.done = false
	o.codes = make(map[uint32]uint64)

	return nil
}

// Run parses the log slice to count all occuring error codes
// logs []log.Info : slice of logs to base the computing on
func (o *CountErrorCodes) Run(args ...interface{}) error {
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
func (o *CountErrorCodes) AfterRun() error {
	return nil
}

// Result returns a map of error codes (the key) and their number of occurence (the value)
// Exemple : nb404 := Result()[404] // gets the number of 404 present in the log
func (o *CountErrorCodes) Result() map[uint32]uint64 {
	return o.codes
}

// IsDone returns true if the task has complted its work. False otherwise.
func (o *CountErrorCodes) IsDone() bool {
	return o.done
}

// Close wipes the object's content. Call Open to use it again.
func (o *CountErrorCodes) Close() error {
	o.codes = nil
	return nil
}
