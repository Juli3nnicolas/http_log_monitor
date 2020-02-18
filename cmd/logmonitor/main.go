package main

import (
	"github.com/Juli3nnicolas/http_log_monitor/pkg/app"
	"github.com/spf13/cobra"
)

var conf *app.Config = &app.Config{}

func main() {

	var rootCmd = &cobra.Command{Use: "logmonitor",
		Short: "Monitors a log file following the common log format",
		Long: `Monitors a log file following the common log format
It displays various metrics such as the request-rate, the number of
successful and failed calls, live details of http-return-codes
and the website with the most traffic. An alerting feature is also
present to be notified when traffic gets awry.`,
		Args: cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			app.Run(conf)
		},
	}

	rootCmd.Flags().StringVarP(&conf.LogFilePath, "path", "p", app.DefaultLogFilePath, "path to the log file to monitor traffic from")
	rootCmd.Flags().DurationVarP(&conf.UpdateFrameDuration, "update", "u", app.DefaultUpdateFrameDuration, "app's refresh rate - rate at which data are going to be fetched and displayed")
	rootCmd.Flags().DurationVarP(&conf.AlertFrameDuration, "alert-period", "T", app.DefaultAlertFrameDuration, "configure alerts' monitoring interval - if the request-rate is above it fro -T, an alert is given")
	rootCmd.Flags().Uint64VarP(&conf.AlertThreshold, "alert-threshold", "t", app.DefaultAlertThreshold, "threshold value, if the request rate is above for -T time, an alert is switched on")
	rootCmd.Execute()
}
