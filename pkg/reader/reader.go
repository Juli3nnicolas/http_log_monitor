package reader

import "github.com/Juli3nnicolas/http_log_monitor/pkg/log"

// Reader is an interface to read data
type Reader interface {
	Open(...interface{}) error
	Read() ([]log.Info, error)
	Close()
}
