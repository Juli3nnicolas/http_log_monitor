package log

import (
	"time"
)

// RFC931User RFC931-compliant log user, this field doesn't seem to be used anymore
type RFC931User struct{}

// Info details a log line's content
type Info struct {
	Host      string
	LogUser   RFC931User
	AuthUser  string
	LocalTime time.Time
	Request   HTTP
}

// HTTP describes an HTTP-request-log
type HTTP struct {
	Method  string
	Route   string
	Code    uint32
	Size    uint64
	Version string
}
