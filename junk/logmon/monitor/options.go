package monitor

import (
	"container/ring"
	"time"
)

// Option allows the reconfiguration of the tail log reader. See its
// implementors for possible configurations.
type Option func(*Reader)

// Threshold is the number of hits per second, during the timespan of analysis
// which shall be used as criterium to trigger an alert.
func Threshold(t uint64) Option {
	return func(r *Reader) {
		r.alertthresh = t
	}
}

// Timespan is analysis timespan for triggering an alert.
func Timespan(dur time.Duration) Option {
	return func(r *Reader) {
		// TODO: handle truncation later. It is seldom a problem in this
		// particular use case.
		seconds := int(dur / time.Second)
		r.alerttimespan = ring.New(seconds)
	}
}
