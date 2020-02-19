package app

import (
	"github.com/Juli3nnicolas/http_log_monitor/pkg/logger"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/reader"
)

// Run executes the entire application (both frontend and backend)
func Run(conf *Config) error {
	l := logger.Get()

	// Init backend
	b := Backend{}

	// Add your tasks to the backend so that it can execute them
	// HACK : the tasks passed to add are already part of the backend
	// struct. It is done so because passing certain parameters to Run
	// and return the task values require  a few abstractions layers
	// that I don't have time to code (define dependencies between tasks
	// a suited view struct and its abstraction layer to convert the tasks
	// into strongly typed values...)
	//
	// However, I hope the code below and the other backend functions will
	// give you a good idea of what I had in mind
	b.add(
		Taskenv{
			Task:       &b.fetchLogs,
			InitParams: []interface{}{conf.LogFilePath, reader.CommonLogFormatParser(), conf.UpdateFrameDuration},
		},
		Taskenv{
			Task: &b.mostHits,
		},
		Taskenv{
			Task: &b.rates,
		},
		Taskenv{
			Task: &b.countCodes,
		},
		Taskenv{
			Task:       &b.alert,
			InitParams: []interface{}{conf.AlertFrameDuration, conf.AlertThreshold},
		},
	)

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
