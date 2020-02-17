package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	path            string
	lines, duration uint64
)

const logline string = `83.149.9.216 - - [17/May/2015:10:05:03 +0000] "GET /presentations/logstash-monitorama-2013/images/kibana-search.png HTTP/1.1" 200 203023 "http://semicomplete.com/presentations/logstash-monitorama-2013/" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.77 Safari/537.36"`

func writeLogs() {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	dur := time.Duration(duration) * time.Second
	gap := time.Second / time.Duration((lines / duration))
	start := time.Now()

	var i = 0
	for time.Now().Sub(start) < dur {
		_, err = fmt.Fprintf(file, logline+"\n")
		if err != nil {
			panic(err)
		}

		fmt.Println(logline)
		i++
		time.Sleep(gap)
	}
	fmt.Println("===================================")
	fmt.Printf("Wrote %d lines in %f seconds\n", i, time.Now().Sub(start).Seconds())
}

func main() {

	var rootCmd = &cobra.Command{Use: "writelogs",
		Short: "Writes log following the common log format",
		Long: `Writes log following the common log format. 
It writes --lines lines in --duration seconds
to the file located at --path. The rate is constant between each writing`,
		Args: cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			writeLogs()
		},
	}

	rootCmd.Flags().StringVarP(&path, "path", "p", "", "path to write the logs (required)")
	rootCmd.MarkFlagRequired("path")
	rootCmd.Flags().Uint64VarP(&lines, "lines", "l", 10, "number of lines to be written")
	rootCmd.Flags().Uint64VarP(&duration, "duration", "d", 1, "writing duration - write -l lines in -d seconds")
	rootCmd.Execute()
}
