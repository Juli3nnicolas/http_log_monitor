package reader

import (
	"fmt"
	"io"
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
	"github.com/papertrail/go-tail/follower"
)

// Tail is a reader able to read an entire file content and seemlessly return file
// updates
type Tail struct {
	timeout time.Duration
	tailer  *follower.Follower
	Parse   Parser
}

// Open opens a file in read mode
func (r *Tail) Open(args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("wrong argument number")
	}

	p, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("invalid type - path must be a string")
	}

	timeout, ok := args[1].(time.Duration)
	if !ok {
		return fmt.Errorf("invalid type - timeout must be a time.Duration argument")
	}
	r.timeout = timeout

	t, err := follower.New(p, follower.Config{
		Whence: io.SeekEnd,
		Offset: 0,
		Reopen: true,
	})
	if err != nil {
		return err
	}
	r.tailer = t

	if r.Parse == nil {
		r.Parse = CommonLogFormatParser()
	}

	return nil
}

// Read reads a file content line by line
// Sleeps wawaiting for data when io.EOF is reached
func (r *Tail) Read() ([]log.Info, error) {

	select {
	case line := <-r.tailer.Lines():
		return r.parseLine(line)

	case <-time.After(r.timeout):
		return nil, nil
	}
}

func (r *Tail) parseLine(line follower.Line) ([]log.Info, error) {
	if r.tailer.Err() != nil {
		return nil, r.tailer.Err()
	}

	parsedLine, err := r.Parse(line.Bytes())
	if err != nil {
		return nil, err
	}

	return []log.Info{parsedLine}, nil
}

// Close closes the file opened with Open
func (r *Tail) Close() {
	r.tailer.Close()
}
