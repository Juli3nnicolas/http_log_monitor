package app

import (
	"fmt"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/task"
)

const (
	alertThresholdMessageFormat string = "Threshold: %d req/s"
	alertDurationMessageFormat  string = "Duration: %s"
	alertMessageHeader          string = "Message:"
	alertOnMessageFormat        string = "High traffic generated an alert - hits = %d, triggered at %v"
	alertOffMessageFormat       string = "Traffic is back to normal - recovery time is %v"
	rateMsgHeader               string = "Frame: "
	rateMsgFormat               string = rateMsgHeader + "%ds Max: %d req/s Avg: %d req/s Success: %d Failure: %d"
	mostHitsNoTraffic           string = "No traffic"
	httpCodes100Header          string = "100:\n"
	httpCodes200Header          string = "200:\n"
	httpCodes300Header          string = "300:\n"
	httpCodes400Header          string = "400:\n"
	httpCodes500Header          string = "500:\n"
)

type rateMsgContent struct {
	frameDuration uint64
	maxReqPSec    uint64
	avgReqPSec    uint64
	nbSuccesses   uint64
	nbFailures    uint64
}

func formatRateMsg(r rateMsgContent) string {
	return fmt.Sprintf(rateMsgFormat, r.frameDuration, r.maxReqPSec, r.avgReqPSec, r.nbSuccesses, r.nbFailures)
}

func formatAlertOnMsg(alert *task.AlertState) string {
	return formatAlertInfoMsg(alert) +
		fmt.Sprintf(
			alertMessageHeader+" "+
				alertOnMessageFormat,
			alert.NbReqs, alert.Date)
}

func formatAlertOffMsg(alert *task.AlertState) string {
	return formatAlertInfoMsg(alert) +
		fmt.Sprintf(
			alertMessageHeader+" "+
				alertOffMessageFormat,
			alert.Date)
}

func formatAlertInfoMsg(alert *task.AlertState) string {
	return fmt.Sprintf(
		alertThresholdMessageFormat+" "+
			alertDurationMessageFormat+" ",
		alert.Threshold,
		alert.Duration.String())
}
