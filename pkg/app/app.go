package app

import (
	"time"
)

// Config is a struct to initialise the application
type Config struct {
	// appErrorLogFile is the path to the app's error log file
	AppErrorLogFile string
	// logFilePath refers to the first file the app will try to read the logs from
	LogFilePath string
	BackendConfig
}

// BackendConfig is a struct to initailise the app's backend
type BackendConfig struct {
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
	if err := b.init(); err != nil {
		return err
	}
	defer b.shutdown()

	// Init view
	r := renderer{}
	ctx, err := r.init()
	if err != nil {
		return err
	}
	defer r.shutdown()

	updateChan := make(chan ViewFrame)
	go r.update(updateChan, LogUpdateError(""))
	go b.run(conf.UpdateFrameDuration, updateChan)
	//go b.run(&conf.BackendConfig, updateChan)

	err = r.render(ctx)
	if err != nil {
		return err
	}

	return nil
}
