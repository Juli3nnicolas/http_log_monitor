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

// Parse reads a log string and returns a properly hydrated Info struct
func Parse(line string) (Info, error) {
	return Info{}, nil
}
