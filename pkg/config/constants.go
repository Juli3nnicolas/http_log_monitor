package config

const (
	// DefaultAppErrorLogFile is the path to the app's error log file
	DefaultAppErrorLogFile string = "/var/log/http_log_monitor.log"
	// DefaultLogFilePath refers to the first file the app will try to read the logs from
	DefaultLogFilePath string = "/tmp/access.log"
	// DefaultUpdateFrameDuration refers to the default time (in seconds) the app will carry out all its measures
	// Said diferently, this value defines the app's refresh rate
	DefaultUpdateFrameDuration uint64 = 10
	// DefaultAlertFrameDuration corresponds to a time-frame during which work related to a single alert will be carried to compute its state
	// It's an alert refreshing rate
	DefaultAlertFrameDuration uint64 = 120
)
