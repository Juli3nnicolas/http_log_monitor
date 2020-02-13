package task

import (
	"fmt"
	"strings"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
)

// FindMostHitSections asynchronously reads a log-file and gets its content.
type FindMostHitSections struct {
	sectionsCount map[string]uint64
	done          bool
}

// Init sets up the async loader for reading the log file
func (o *FindMostHitSections) Init(args ...interface{}) error {
	return nil
}

// BeforeRun Starts the reading process. Just internally fetch data.
func (o *FindMostHitSections) BeforeRun() error {
	o.done = false

	return nil
}

// Run copies the files content to its internal buffer for sharing with other
func (o *FindMostHitSections) Run(args ...interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong parameters - only one parameter is supported, it must be a []log.Info")
	}

	logs, ok := args[0].([]log.Info)
	if !ok {
		return fmt.Errorf("type error - got %T instead of []log.Info", args[0])
	}

	logsLen := len(logs)
	o.sectionsCount = make(map[string]uint64, logsLen)

	for i := 0; i < logsLen; i++ {
		s := extractSection(logs[i].Request.Route)
		if s != "" {
			o.sectionsCount[s]++
		}
	}
	o.done = true

	return nil
}

// Result returns the result of the work carried out by the task
func (o *FindMostHitSections) Result() map[string]uint64 {
	if o.IsDone() {
		return o.sectionsCount
	}

	return nil
}

// AfterRun ceases reading and prepares a new buffer for reading data during
func (o *FindMostHitSections) AfterRun() error {
	return nil
}

// IsDone is true when the task has completed
func (o *FindMostHitSections) IsDone() bool {
	return o.done
}

// Close closes the task. Call Init to use it again.
func (o *FindMostHitSections) Close() error {
	o.done = false
	return nil
}

// section returns the url-part before the second "/". It returns / when finds
// an empty section (the root). It returns nil in case no "/" is found.
// Note -  I think it is useful to log "/" as redirections to proxies or
// other services might occur (especially with rather complex web server configs).
func extractSection(url string) string {
	parts := strings.Split(url, "/")
	partsLen := len(parts)

	if partsLen <= 2 {
		if partsLen == 2 {
			return "/"
		}

		return ""
	}

	return "/" + parts[1]
}
