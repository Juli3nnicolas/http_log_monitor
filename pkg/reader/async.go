package reader

import "github.com/Juli3nnicolas/http_log_monitor/pkg/log"

// Async calls a reader asynchronously and stores its output into a dedicated buffer
// At the moment, each Async readers correctly handle one correct parralel read.
// Create several readers to read several streams in parallel.
// Only one parallel call is supported, because I don't need more fore the test.
type Async struct {
	Reader     Reader
	MinBufsize uint64
	err        error
	stop       chan bool
	buffer     []log.Info
}

const defaultMinBufSize uint64 = 500

func (r *Async) init() {
	if r.MinBufsize == 0 {
		r.MinBufsize = defaultMinBufSize
	}
	r.buffer = make([]log.Info, 0, r.MinBufsize)
	r.stop = make(chan bool)
	r.err = nil
}

// Open inits the sync reader and opens Reader.
func (r *Async) Open(args ...interface{}) error {
	r.init()
	return r.Reader.Open(args...)
}

// Close stops reading data and closes Reader.
func (r *Async) Close() {
	r.Stop()
	r.Reader.Close()
}

// Read returns the content that has already been read. Call r.Start() to initiate the process.
// Make sure to call Stop before if you want accurate data.
// Otherwise, the content can still be updated while you read it.
func (r *Async) Read() ([]log.Info, error) {
	return r.buffer, r.err
}

// Start starts a new parallel reading process. Calls Reader to read.
// At the moment only one parallel read is supported (One call to Start)
// This call is possible though: r.Start(); r.Stop(); r.Start();
func (r *Async) Start() {
	go r.read()
}

// Stop stops the reading process
func (r *Async) Stop() {
	r.stop <- true
}

// Flush empties all internal buffers. Make sure to call Stop before otherwise
// the read function may crash.
// If Reader implements a similar method, it must be called manually.
func (r *Async) Flush() {
	r.init()
}

// read reads the file content asynchronously. It fills r.buffer.
func (r *Async) read() error {
	for {
		select {
		case s, ok := <-r.stop:
			if s || !ok {
				r.err = nil
				return nil
			}
		default:
			output, err := r.Reader.Read()
			if err != nil {
				r.err = err
				return err
			}
			if output != nil {
				r.buffer = append(r.buffer, output...)
			}
		}
	}
}
