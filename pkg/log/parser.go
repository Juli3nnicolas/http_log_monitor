package log

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// This parser is able to read the Common Log Format used by Apache's httpd
// The common log format reads as follows
// 172.17.0.1 - - [09/Feb/2020:16:27:00 +0000] "GET / HTTP/1.1" 200 612
// Where the orderly enumerated values correspond to :
// - 172.17.0.1 is the host value
// - "-" is an RFC931User. This value doesn't seem to be used anymore
// - "-" is an authenticated user (if any)
// - [09/Feb/2020:16:27:00 +0000] is a local timestamp.
// This time format is specific to the common log format.
// - "GET / HTTP/1.1" 200 612 describes several http-related data :
// - GET is a http method
// - / is the route that have been served
// - HTTP/1.1 is the protocol version
// - 200 is the http-return-code
// - 612 is the served-content's size
// Please bear in mind that additional, non-standard information can be provided
// by servers.

// Parse reads a log string and returns a properly hydrated Info struct
func Parse(line string) (Info, error) {
	if line == "" {
		return Info{}, fmt.Errorf("log.Parse error - empty line")
	}

	fields := strings.Split(line, " ")

	localTime, err := parseLocalTime(strings.Join(fields[3:5], " "))
	if err != nil {
		return Info{}, err
	}

	httpReq, err := parseHTTP(fields[5:10])
	if err != nil {
		return Info{}, err
	}

	info := Info{
		Host:      parseHost(fields[0]),
		LogUser:   parseRFC931User(fields[1]),
		AuthUser:  parseAuthUser(fields[2]),
		LocalTime: localTime,
		Request:   httpReq,
	}

	return info, nil
}

func parseHost(field string) string {
	return field
}

func parseRFC931User(field string) (user RFC931User) {
	return
}

func parseAuthUser(field string) string {
	if field == "-" {
		return ""
	}
	return field
}

func parseLocalTime(field string) (time.Time, error) {
	const commonLogFormat string = `02/Jan/2006:15:04:05 -0700`

	field = field[1 : len(field)-1]
	t, err := time.Parse(commonLogFormat, field)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func parseHTTP(fields []string) (HTTP, error) {
	code, err := strconv.Atoi(fields[3])
	if err != nil {
		return HTTP{}, err
	}

	size, err := strconv.Atoi(fields[4])
	if err != nil {
		return HTTP{}, err
	}

	return HTTP{
		Method:  strings.Trim(fields[0], `"`),
		Route:   fields[1],
		Version: strings.Trim(fields[2], `"`),
		Code:    uint32(code),
		Size:    uint64(size),
	}, nil
}
