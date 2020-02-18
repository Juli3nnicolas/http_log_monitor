package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// New creates a new logger. Do so if you may log from another go routine
func init() {
	// DefaultAppErrorLogFile is the path to the app's error log file
	const DefaultAppErrorLogFile string = "/usr/local/var/log/http_log_monitor.log"

	log = logrus.New()
	log.ReportCaller = true
	file, err := os.OpenFile(DefaultAppErrorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
	}

	log.Out = file
}
