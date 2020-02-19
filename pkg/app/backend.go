package app

import (
	"os"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/logger"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/task"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/timer"
)

// Backend is the structure that holds all task to have metrics displayed
type Backend struct {
	fetchLogs  task.FetchLogs
	mostHits   task.FindMostHitSections
	rates      task.MeasureRates
	countCodes task.CountErrorCodes
	alert      task.Alert
	tasks      []Taskenv
}

// Taskenv is a task and all its necessary environment to be executed
type Taskenv struct {
	Task       task.Task
	InitParams []interface{}
}

func (b *Backend) add(tasks ...Taskenv) {
	b.tasks = tasks
}

func (b *Backend) init(conf *Config) error {
	// Create input log file if it doesn't exist
	f, err := os.OpenFile(conf.LogFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	f.Close()

	// Initialise all tasks
	for _, t := range b.tasks {
		if err := t.Task.Init(t.InitParams...); err != nil {
			return err
		}
	}

	return nil
}

func (b *Backend) run(conf *Config, outputChan chan ViewFrame) {
	l := logger.Get()

	frame := conf.UpdateFrameDuration
	t := &timer.Time{}
	start := t.Now()
	for {

		for _, t := range b.tasks {
			if err := t.Task.BeforeRun(); err != nil {
				l.Fatalf(err.Error())
				return
			}
		}

		var err error
		allDone := false
		resultSent := false

		for t.Now().Sub(start) < frame {
			if !b.fetchLogs.IsDone() {
				if err = b.fetchLogs.Run(); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}
			logs := b.fetchLogs.Fetch()

			if !b.mostHits.IsDone() {
				if err = b.mostHits.Run(logs); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}

			if !b.rates.IsDone() {
				if err = b.rates.Run(logs, uint64(frame.Seconds())); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}

			if !b.countCodes.IsDone() {
				if err = b.countCodes.Run(logs); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}

			if !b.alert.IsDone() && b.rates.IsDone() {
				if err = b.alert.Run(b.rates.Result(), t); err != nil {
					l.Fatalf(err.Error())
					return
				}
				allDone = true
			}

			if allDone && !resultSent {
				view := ViewFrame{
					Hits:  b.mostHits.Result(),
					Rates: b.rates.Result(),
					Codes: b.countCodes.Result(),
					Alert: b.alert.Result(),
				}
				outputChan <- view
				resultSent = true
			}
		}

		start = t.Now()

		for _, t := range b.tasks {
			if err := t.Task.AfterRun(); err != nil {
				l.Fatalf(err.Error())
				return
			}
		}
	}
}

func (b *Backend) shutdown() error {
	for _, t := range b.tasks {
		if err := t.Task.Close(); err != nil {
			return err
		}
	}

	return nil
}
