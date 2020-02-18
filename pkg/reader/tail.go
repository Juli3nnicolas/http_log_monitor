package reader

import (
	"fmt"
	"io"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
	"github.com/papertrail/go-tail/follower"
)

// Tail is a reader able to read an entire file content and seemlessly return file
// updates
type Tail struct {
	tailer *follower.Follower
	Parse  Parser
}

// Open opens a file in read mode
func (r *Tail) Open(path ...interface{}) error {
	if len(path) != 1 {
		return fmt.Errorf("wrong argument number")
	}

	p, ok := path[0].(string)
	if !ok {
		return fmt.Errorf("invalid type - path must be a string")
	}

	t, err := follower.New(p, follower.Config{
		Whence: io.SeekStart,
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

	line := <-r.tailer.Lines()
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
