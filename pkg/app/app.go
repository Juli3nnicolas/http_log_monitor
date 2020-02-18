package app

import (
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/logger"
)

// Config is a struct to initialise the application
type Config struct {
	// logFilePath refers to the first file the app will try to read the logs from
	LogFilePath string
	// updateFrameDuration refers to the default time the app will carry out all its measures
	// Said diferently, this value defines the app's backend refresh rate
	UpdateFrameDuration time.Duration
	// alertFrameDuration corresponds to a time-frame during which work related to a single alert will be carried to compute its state
	// It's an alert refreshing rate
	AlertFrameDuration time.Duration
	// alertThreshold is the default threshold (in req/s) triggering an alert
	AlertThreshold uint64
}

// Run executes the entire application (both frontend and backend)
func Run(conf *Config) error {
	// Init backend
	b := Backend{}

	// Init view
	r := renderer{}
	ctx, err := r.init()
	if err != nil {
		logger.Get().Fatalf(err.Error())
		return err
	}
	defer r.shutdown()

	updateChan := make(chan ViewFrame)
	go r.update(updateChan, LogUpdateError())
	go b.run(conf, updateChan)

	err = r.render(ctx)
	if err != nil {
		logger.Get().Fatalf(err.Error())
		return err
	}

	return nil
}
