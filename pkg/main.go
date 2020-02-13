package main

import (
	"fmt"
	"time"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/task"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/reader"
)

func main() {
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

	const frame time.Duration = 1 * time.Second
	start := time.Now()

	for {
		err = fetchLogs.BeforeRun()
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
				fmt.Println("Most hits ", mostHits.Result())
			}

			if !rates.IsDone() {
				if err = rates.Run(logs, uint64(frame.Seconds())); err != nil {
					panic(err)
				}
				fmt.Println("Rates ", rates.Result())
			}

			if !countCodes.IsDone() {
				if err = countCodes.Run(logs); err != nil {
					panic(err)
				}
				fmt.Println("Rates ", countCodes.Result())
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
