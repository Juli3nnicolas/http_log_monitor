package timer

import "time"

// Timer is a wrapper interface for the standard time package
// It has been done so to ease testing
type Timer interface {
	// Now returns the current time
	Now() time.Time
}

// Time is the production timer, relies on the standard time package
type Time struct{}

// Now returns the current time. Calls time.Now()
func (t *Time) Now() time.Time {
	return time.Now()
}

// TimeStub is a struct used for testing
type TimeStub struct {
	NowStub func() time.Time
}

// Now calles NowStub. Set a function to NowStub to customize your tests
func (t *TimeStub) Now() time.Time {
	return t.NowStub()
}
