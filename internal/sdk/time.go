package sdk

import "time"

// DefaultNowTime returns the curren wall clock time.
func DefaultNowTime() time.Time {
	return time.Now()
}

// NowTime is used to get the current time. Can be overriden for tests.
var NowTime = DefaultNowTime
