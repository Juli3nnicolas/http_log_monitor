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

	const frame time.Duration = 1 * time.Second
	start := time.Now()

	for {
		err = fetchLogs.BeforeRun()
		if err != nil {
			panic(err)
		}

		for time.Now().Sub(start) < frame {
			if !fetchLogs.IsDone() {
				fetchLogs.Run()
				fmt.Println(fetchLogs.Fetch())
			}
		}

		start = time.Now()
		err = fetchLogs.AfterRun()
	}

	err = fetchLogs.Close()
	if err != nil {
		panic(err)
	}
}
