package logger

import (
	"os"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/config"
	"github.com/sirupsen/logrus"
)

// New creates a new logger. Do so if you may log from another go routine
func New() *logrus.Logger {
	log := logrus.New()
	log.ReportCaller = true
	file, err := os.OpenFile(config.DefaultAppErrorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
		return log
	}

	log.Out = file

	return log
}
