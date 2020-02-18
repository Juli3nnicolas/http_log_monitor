package app

import (
	"github.com/Juli3nnicolas/http_log_monitor/pkg/logger"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/reader"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/task"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/timer"
)

type Backend struct {
}

func (b *Backend) run(conf *Config, outputChan chan ViewFrame) {
	l := logger.Get()

	fetchLogs := task.FetchLogs{}
	err := fetchLogs.Init(conf.LogFilePath, reader.CommonLogFormatParser(), conf.UpdateFrameDuration)
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	mostHits := task.FindMostHitSections{}
	err = mostHits.Init()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	rates := task.MeasureRates{}
	err = rates.Init()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	countCodes := task.CountErrorCodes{}
	err = countCodes.Init()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	alert := task.Alert{}
	err = alert.Init(conf.AlertFrameDuration, conf.AlertThreshold)
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	frame := conf.UpdateFrameDuration
	t := &timer.Time{}
	start := t.Now()
	for {
		err := fetchLogs.BeforeRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = mostHits.BeforeRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = rates.BeforeRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = countCodes.BeforeRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = alert.BeforeRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		allDone := false
		resultSent := false
		for t.Now().Sub(start) < frame {
			if !fetchLogs.IsDone() {
				if err = fetchLogs.Run(); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}
			logs := fetchLogs.Fetch()

			if !mostHits.IsDone() {
				if err = mostHits.Run(logs); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}

			if !rates.IsDone() {
				if err = rates.Run(logs, uint64(frame.Seconds())); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}

			if !countCodes.IsDone() {
				if err = countCodes.Run(logs); err != nil {
					l.Fatalf(err.Error())
					return
				}
			}

			if !alert.IsDone() && rates.IsDone() {
				if err = alert.Run(rates.Result(), t); err != nil {
					l.Fatalf(err.Error())
					return
				}
				allDone = true
			}

			if allDone && !resultSent {
				view := ViewFrame{
					Hits:  mostHits.Result(),
					Rates: rates.Result(),
					Codes: countCodes.Result(),
					Alert: alert.Result(),
				}
				outputChan <- view
				resultSent = true
			}
		}

		start = t.Now()
		err = fetchLogs.AfterRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = mostHits.AfterRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = rates.AfterRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = countCodes.AfterRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}

		err = alert.AfterRun()
		if err != nil {
			l.Fatalf(err.Error())
			return
		}
	}

	err = fetchLogs.Close()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	err = mostHits.Close()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	err = rates.Close()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	err = countCodes.Close()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}

	err = alert.Close()
	if err != nil {
		l.Fatalf(err.Error())
		return
	}
}

func (b *Backend) shutdown() error {

	return nil
}
