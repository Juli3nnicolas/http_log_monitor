package reader

import (
	"fmt"
	"io"
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/config"
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

	// TODO : Remove this value, it should be injected from open
	case <-time.After(config.DefaultUpdateFrameDuration):
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
