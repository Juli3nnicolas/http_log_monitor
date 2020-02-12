package task

import (
	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/reader"
)

// Poller is a task asynchronously reading a file.
// It implements the double buffering technique to ensure
// lock-free thread-safe access.
type Poller struct {
	reader  reader.Async
	done    bool
	wBuf    uint8 // write buffer index
	buffers [2][]log.Info
}

// Init sets up the poll to asynchronously read the file pointed to by path
// The parser is used by Run to fill the readable log.Info buffer.
func (o *Poller) Init(path string, parser reader.Parser) error {
	fileReader := reader.File{Parse: parser}
	o.reader.Reader = &fileReader

	err := o.reader.Open(path)
	if err != nil {
		return err
	}

	return nil
}

// Run starts the Poller task. Data will have been updated once
// IsDone is true.
// IMPORTANT ; It IS SAFE to call Poll whenever you want.
func (o *Poller) Run() error {
	o.done = false
	o.reader.Start()
	return nil
}

// Poll returns the file content read before Swap was called.
func (o *Poller) Poll() []log.Info {
	return o.buffers[nextBuf(o.wBuf)]
}

// IsDone returns true when data have been updated, false while they are.
// It is safe to call Poll whenever you want since this class implements double
// buffering.
func (o *Poller) IsDone() bool {
	return o.done
}

// Swap copies written data to the back buffer so that they can be Polled.
// The back buffer is swapped to the front position so that it can be read.
// Consequently, the front buffer is swapped to the back position so that it
// can be written to without causing any concurrency issue.
// The Task is flagged as done after calling this.
func (o *Poller) Swap() error {
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

	// Tell everyone data have been updated
	o.done = true

	return nil
}

func nextBuf(index uint8) uint8 {
	return (index + 1) % 2
}
