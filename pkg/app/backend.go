package app

import (
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/reader"
	"github.com/Juli3nnicolas/http_log_monitor/pkg/task"
)

type Backend struct {
}

func (b *Backend) init() error {
	fetchLogs := task.FetchLogs{}
	err := fetchLogs.Init("/tmp/access.log", reader.CommonLogFormatParser())
	if err != nil {
		return err
	}

	mostHits := task.FindMostHitSections{}
	err = mostHits.Init()
	if err != nil {
		return err
	}

	rates := task.MeasureRates{}
	err = rates.Init()
	if err != nil {
		return err
	}

	countCodes := task.CountErrorCodes{}
	err = countCodes.Init()
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) run(frame time.Duration, outputChan chan ViewFrame) {
	fetchLogs := task.FetchLogs{}
	err := fetchLogs.Init("/tmp/access.log", reader.CommonLogFormatParser())
	if err != nil {
		panic(err)
	}

	mostHits := task.FindMostHitSections{}
	err = mostHits.Init()
	if err != nil {
		panic(err)
	}

	rates := task.MeasureRates{}
	err = rates.Init()
	if err != nil {
		panic(err)
	}

	countCodes := task.CountErrorCodes{}
	err = countCodes.Init()
	if err != nil {
		panic(err)
	}

	start := time.Now()

	for {
		err := fetchLogs.BeforeRun()
		if err != nil {
			panic(err)
		}

		err = mostHits.BeforeRun()
		if err != nil {
			panic(err)
		}

		err = rates.BeforeRun()
		if err != nil {
			panic(err)
		}

		err = countCodes.BeforeRun()
		if err != nil {
			panic(err)
		}

		allDone := false
		resultSent := false
		for time.Now().Sub(start) < frame {
			if !fetchLogs.IsDone() {
				if err = fetchLogs.Run(); err != nil {
					panic(err)
				}
			}
			logs := fetchLogs.Fetch()

			if !mostHits.IsDone() {
				if err = mostHits.Run(logs); err != nil {
					panic(err)
				}
			}

			if !rates.IsDone() {
				if err = rates.Run(logs, uint64(frame.Seconds())); err != nil {
					panic(err)
				}
			}

			if !countCodes.IsDone() {
				if err = countCodes.Run(logs); err != nil {
					panic(err)
				}
				allDone = true
			}

			if allDone && !resultSent {
				view := ViewFrame{
					Hits:  mostHits.Result(),
					Rates: rates.Result(),
					Codes: countCodes.Result(),
				}
				outputChan <- view
				resultSent = true
			}
		}

		start = time.Now()
		err = fetchLogs.AfterRun()
		if err != nil {
			panic(err)
		}

		err = mostHits.AfterRun()
		if err != nil {
			panic(err)
		}

		err = rates.AfterRun()
		if err != nil {
			panic(err)
		}

		err = countCodes.AfterRun()
		if err != nil {
			panic(err)
		}
	}

	err = fetchLogs.Close()
	if err != nil {
		panic(err)
	}

	err = mostHits.Close()
	if err != nil {
		panic(err)
	}

	err = rates.Close()
	if err != nil {
		panic(err)
	}

	err = countCodes.Close()
	if err != nil {
		panic(err)
	}
}

func (b *Backend) shutdown() error {

	return nil
}
