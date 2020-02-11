package reader

import "github.com/Juli3nnicolas/http_log_monitor/pkg/log"

// Stub is a stub for the reader.Reader interface. Open calls OpenStub, Read calls ReadStub...
type Stub struct {
	OpenStub  func(...interface{}) error
	ReadStub  func() ([]log.Info, error)
	CloseStub func()
}

// Open calls OpenStub
func (r *Stub) Open(args ...interface{}) error {
	return r.OpenStub(args)
}

// Read calls ReadStub
func (r *Stub) Read() ([]log.Info, error) {
	return r.ReadStub()
}

// Close calls CloseStub
func (r *Stub) Close() {
	r.CloseStub()
}
