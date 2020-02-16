package app

import "fmt"

const (
	alertThresholdHeader string = "Threshold (req/s): "
	alertDurationHeader  string = "Duration (s): "
	alertMessageHeader   string = "Message: "
	rateMsgHeader        string = "Frame: "
	rateMsgFormat        string = rateMsgHeader + "%d s Max: %d req/s Avg: %d req/s Success: %d Failure: %d"
	mostHitsNoTraffic    string = "No traffic"
	httpCodes100Header   string = "100:\n"
	httpCodes200Header   string = "200:\n"
	httpCodes300Header   string = "300:\n"
	httpCodes400Header   string = "400:\n"
	httpCodes500Header   string = "500:\n"
)

type rateMsgContent struct {
	frameDuration uint64
	maxReqPSec    uint64
	avgReqPSec    uint64
	nbSuccesses   uint64
	nbFailures    uint64
}

func formatRatesMsg(r rateMsgContent) string {
	return fmt.Sprintf(rateMsgFormat, r.frameDuration, r.maxReqPSec, r.avgReqPSec, r.nbSuccesses, r.nbFailures)
}
