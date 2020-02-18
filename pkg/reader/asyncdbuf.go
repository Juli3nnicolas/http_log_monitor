package reader

import (
	"fmt"
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
)

// ASyncDBuf asynchronously reads a file.
// It implements the double buffering technique to ensure
// lock-free thread-safe access.
// ASyncDBuf stands for asynchronous double-buffering-reader
type ASyncDBuf struct {
	reader  Async
	wBuf    uint8 // write buffer index
	buffers [2][]log.Info
}

// Open inits the reader to asynchronously read the file pointed to by path
// The parser is used by Run to fill the readable log.Info buffer.
// 1st param : path to file to read
// 2nd param : a reader.Parser function, used to fill the buffer with log.Info data
func (o *ASyncDBuf) Open(args ...interface{}) error {
	if len(args) != 3 {
		return fmt.Errorf("wrong parameter number")
	}

	path, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("first parameter should be a string")
	}

	parser, ok := args[1].(Parser)
	if !ok {
		return fmt.Errorf("second parameter should be a reader.Parser")
	}

	timeout, ok := args[2].(time.Duration)
	if !ok {
		return fmt.Errorf("invalid type - timeout must be a time.Duration argument")
	}

	fileReader := Tail{Parse: parser}
	o.reader.Reader = &fileReader

	err := o.reader.Open(path, timeout)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the reader. To reuse it, call Open.
func (o *ASyncDBuf) Close() {
	o.wBuf = 0
	o.reader.Close()
}

// Run starts the ASyncDBuf task. Data will have been updated once
// IsDone is true.
// IMPORTANT ; It IS SAFE to call Read whenever you want.
func (o *ASyncDBuf) Run() error {
	o.reader.Start()
	return nil
}

// Read returns the file content read before Swap was called.
func (o *ASyncDBuf) Read() ([]log.Info, error) {
	return o.buffers[nextBuf(o.wBuf)], nil
}

// Swap copies written data to the back buffer so that they can be Read.
// The back buffer is swapped to the front position so that it can be read.
// Consequently, the front buffer is swapped to the back position so that it
// can be written to without causing any concurrency issue.
// The Task is flagged as done after calling this.
func (o *ASyncDBuf) Swap() error {
	o.reader.Stop()
	buf, err := o.reader.Read()
	if err != nil {
		return err
	}

	// Copy the current work to the front buffer
	o.buffers[o.wBuf] = make([]log.Info, len(buf))
	copy(o.buffers[o.wBuf], buf)

	// Swap front and back buffers
	o.wBuf = nextBuf(o.wBuf)

	// Flush reader to avoid rewritting previous results to the new buffer
	o.reader.Flush()

	return nil
}

func nextBuf(index uint8) uint8 {
	return (index + 1) % 2
}
