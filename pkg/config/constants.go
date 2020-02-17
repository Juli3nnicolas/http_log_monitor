package config

import "time"

const (
	// DefaultAppErrorLogFile is the path to the app's error log file
	DefaultAppErrorLogFile string = "/usr/local/var/log/http_log_monitor.log"
	// DefaultLogFilePath refers to the first file the app will try to read the logs from
	DefaultLogFilePath string = "/tmp/access.log"
	// DefaultUpdateFrameDuration refers to the default time the app will carry out all its measures
	// Said diferently, this value defines the app's backend refresh rate
	DefaultUpdateFrameDuration time.Duration = 10 * time.Second
	// DefaultAlertFrameDuration corresponds to a time-frame during which work related to a single alert will be carried to compute its state
	// It's an alert refreshing rate
	DefaultAlertFrameDuration time.Duration = 2 * time.Minute
	// DefaultAlertThreshold is the default threshold (in req/s) triggering an alert
	DefaultAlertThreshold uint64 = 10
)
