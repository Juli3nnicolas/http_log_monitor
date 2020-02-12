package task

import (
	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/reader"
)

// FetchLogs asynchronously reads a log-file and gets its content.
type FetchLogs struct {
	logs []log.Info
	dbuf reader.ASyncDBuf
	done bool
}

// Init sets up the async loader for reading the log file
// Give it a :
// path string : path to log file
// parser Parser : a log-parsing function such as one returned
// by reader.CommonLogFormatParser
func (o *FetchLogs) Init(args ...interface{}) error {
	err := o.dbuf.Open(args...)
	if err != nil {
		return err
	}

	return nil
}

// BeforeRun Starts the reading process. Just internally fetch data.
func (o *FetchLogs) BeforeRun() error {
	o.done = false
	err := o.dbuf.Run()
	return err
}

// Run copies the files content to its internal buffer for sharing with other
// tasks. The read content depends on time-frames' duration and the log-file's
// writing-rate
func (o *FetchLogs) Run() error {
	var err error
	o.logs, err = o.dbuf.Read()
	if err != nil {
		return err
	}

	o.done = true

	return err
}

// Fetch returns the logs from the input log-file.
// The returned entry should only be read from. Otherwise
// the effects can be unpredictable.
func (o *FetchLogs) Fetch() []log.Info {
	return o.logs
}

// AfterRun ceases reading and prepares a new buffer for reading data during
// the next time-frame
func (o *FetchLogs) AfterRun() error {
	return o.dbuf.Swap()
}

// IsDone is true when the task has completed
func (o *FetchLogs) IsDone() bool {
	return o.done
}

// Close closes the task. Call Init to use it again.
func (o *FetchLogs) Close() error {
	o.done = false
	o.dbuf.Close()

	return nil
}
