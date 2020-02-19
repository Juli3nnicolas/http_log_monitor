package task

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
)

// FindMostHitSections reads the logs and agregate them by sections. A log occurence is called a hit.
// This task returns information on hits (how many times they were found, how were they called...).
type FindMostHitSections struct {
	sectionHits []Hit
	done        bool
}

// Hit is a structure representing a request occurence in the log file.
type Hit struct {
	Section string
	Total   uint64
	Methods map[string]uint64
}

// set is used to properly fill a Hit. It requires a url section and an HTTP method
func (o *Hit) set(section, method string) {
	if o.Methods == nil {
		o.Methods = make(map[string]uint64)
	}
	o.Section = section
	o.Total++
	o.Methods[method]++
}

// Init does nothing, implemented to comply with the task interface.
func (o *FindMostHitSections) Init(args ...interface{}) error {
	return nil
}

// BeforeRun sets the task as not done. IsDone is going to return to false.
func (o *FindMostHitSections) BeforeRun(...interface{}) error {
	o.done = false

	return nil
}

// Run parses a []log.Info and aggregates its content by sections.
// Then every occurence's data is accounted for so that the summary
// can returned by Result.
// IMPORTAN : this method expects a []log.Info input to function
func (o *FindMostHitSections) Run(args ...interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong parameters - only one parameter is supported, it must be a []log.Info")
	}

	logs, ok := args[0].([]log.Info)
	if !ok {
		return fmt.Errorf("type error - got %T instead of []log.Info", args[0])
	}

	hitMap := make(map[string]int)
	o.sectionHits = make([]Hit, 0, 1)

	for i := range logs {
		section := extractSection(logs[i].Request.Route)
		if section != "" {
			j, found := hitMap[section]

			if found {
				o.sectionHits[j].set(section, logs[i].Request.Method)
			} else {
				hit := Hit{}
				hit.set(section, logs[i].Request.Method)
				o.sectionHits = append(o.sectionHits, hit)
				hitMap[section] = len(o.sectionHits) - 1
			}
		}
	}

	sort.Sort(sectionHits(o.sectionHits))
	o.done = true

	return nil
}

type sectionHits []Hit

func (s sectionHits) Len() int {
	return len(s)
}

func (s sectionHits) Less(i, j int) bool {
	return s[i].Total > s[j].Total
}

func (s sectionHits) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

// Result returns the result of the work carried out by the task if the task is done. Returns nil otherwise.
// Returns a slice of section-hits in decreasing-order on Hit.Total
func (o *FindMostHitSections) Result() []Hit {
	if o.IsDone() {
		return o.sectionHits
	}

	return nil
}

// AfterRun does nothing, must be implemented to implement to task interface
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
