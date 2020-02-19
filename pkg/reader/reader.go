package reader

import "github.com/Juli3nnicolas/http_log_monitor/pkg/log"

// Reader is an interface to read data
type Reader interface {
	// Open prepares an object for reading
	Open(...interface{}) error
	// Read reads the object content and returns formatted logs if any
	Read() ([]log.Info, error)
	// Close closes the object and all resources used for reading the file
	Close()
}
