package app

import (
	"github.com/Juli3nnicolas/http_log_monitor/pkg/logger"
)

// Run executes the entire application (both frontend and backend)
func Run(conf *Config) error {
	l := logger.Get()

	// Init backend
	b := Backend{}
	err := b.init(conf)
	if err != nil {
		l.Fatalf(err.Error())
		return err
	}
	defer func() {
		if err := b.shutdown(); err != nil {
			logger.Get().Fatalf(err.Error())
		}
	}()

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
