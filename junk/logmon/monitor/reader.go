package monitor

import (
	"container/ring"
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"xojoc.pw/logparse"
)

var errImproperUse = errors.New("cannot use tail reader directly. It must be started with tail.Open()")

// Stat describes the relevant information about the incoming traffic of a
// section.
type Stat struct {
	Hits                                       uint64
	Status5xx, Status4xx, Status3xx, Status2xx uint64
}

// Alert encapsulates one alert event, from its start to its conclusion.
type Alert struct {
	Start, End time.Time
	Hits       uint64
}

// Reader implements a Common Format log reader, that internally processes
// further information from each parsed line. It will extract a section from the
// incoming traffic assuming that everything that is under the first level deep
// in URL is a section. It also does traffic-triggered alerts. Default: 1000
// hits per second in average for more than 2 minutes.
type Reader struct {
	alertthresh   uint64
	alerttimespan *ring.Ring

	err error

	mu        sync.Mutex
	stats     map[string]Stat
	totalhits uint64
	alerts    []Alert
}

// Scanner describes the expected interface for any datasource to be use by this
// monitor. Currently only cirello.io/logmon/tail implements it, but in practice
// it can be anything.
type Scanner interface {
	Scan() bool
	Text() string
}

// Open processes the "content"and read it until context cancelation.
// Make sure to check Err() after the use of this reader.
func Open(ctx context.Context, content Scanner, opts ...Option) *Reader {
	lr := &Reader{
		stats:         make(map[string]Stat),
		alerttimespan: ring.New(120), // 120 seconds - 2min
		alertthresh:   1000,          // 1000 hits per second
	}
	for _, opt := range opts {
		opt(lr)
	}
	go lr.poll(ctx, content)
	return lr
}

// Alerts return the list of all alerts that happened during the lifespan of
// this Reader.
func (r *Reader) Alerts() []Alert {
	if r.stats == nil {
		r.err = errImproperUse
		return nil
	}

	r.mu.Lock()
	alerts := make([]Alert, len(r.alerts))
	copy(alerts, r.alerts)
	r.mu.Unlock()
	return alerts
}

// Stats returns an unordered map of sections with their respective statistics.
func (r *Reader) Stats() map[string]Stat {
	if r.stats == nil {
		r.err = errImproperUse
		return nil
	}

	statscopy := make(map[string]Stat)
	r.mu.Lock()
	for k, v := range r.stats {
		statscopy[k] = v
	}
	r.mu.Unlock()
	return statscopy
}

// Err shall return non-nil in case of any internal failure. Must be invoked
// as soon as this Reader is done.
func (r *Reader) Err() error {
	return r.err
}

// Here is where we put all the components of the clockwork together. It all
// starts with parselog. The flow is rather simple: parselog reads each incoming
// line and keeps an internal track of their meaning (sections, hit counts, per
// status response hit count and whatnots). Stats and Alerts are a simple
// snapshots of this internal state, coordinated through a struct-level mutex.
// Monitor has an internal clock that takes an snapshot of the overall hit
// counter per second and then compares with the configured threshold in order
// to trigger or recover from alerts.
func (r *Reader) poll(ctx context.Context, content Scanner) {
	ctx, cancel := context.WithCancel(ctx)

	go r.parselog(ctx, cancel, content)
	go r.monitor(ctx)
}

func (r *Reader) parselog(ctx context.Context, cancel context.CancelFunc, content Scanner) {
	defer cancel()

	for content.Scan() {
		select {
		case <-ctx.Done():
			return

		default:
			if err := r.parseline(content.Text()); err != nil {
				r.err = err
				return
			}
		}
	}
}

func (r *Reader) parseline(line string) error {
	ll, err := logparse.Common(line)
	if err != nil {
		return err
	}

	section := extractSection(ll.Request.URL)
	r.mu.Lock()
	r.totalhits++
	s, ok := r.stats[section]
	if !ok {
		r.stats[section] = Stat{}
		s = r.stats[section]
	}
	r.mu.Unlock()

	s.Hits++
	if ll.Status >= http.StatusInternalServerError {
		s.Status5xx++
	} else if ll.Status >= http.StatusBadRequest {
		s.Status4xx++
	} else if ll.Status >= http.StatusMultipleChoices {
		s.Status3xx++
	} else if ll.Status >= http.StatusOK {
		s.Status2xx++
	}

	r.mu.Lock()
	r.stats[section] = s
	r.mu.Unlock()
	return nil
}

func (r *Reader) monitor(ctx context.Context) {
	var alarmed bool
	t := time.Tick(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return

		case <-t:
			alarmed = r.checkalert(alarmed)
		}
	}
}

func (r *Reader) checkalert(alarmed bool) bool {
	r.mu.Lock()
	th := r.totalhits
	r.mu.Unlock()

	r.alerttimespan = r.alerttimespan.Next()
	r.alerttimespan.Value = th

	var max, min uint64
	var l int
	r.alerttimespan.Do(func(e interface{}) {
		if e == nil {
			return
		}
		l++
		v := e.(uint64)
		if v > max {
			max = v
		} else if v < min || min == 0 {
			min = v
		}
	})
	if l < r.alerttimespan.Len() {
		return false
	}

	avg := (max - min) / (uint64(l) - 1)

	if avg >= r.alertthresh && !alarmed {
		r.mu.Lock()
		r.alerts = append(r.alerts, Alert{Start: time.Now(), Hits: avg})
		r.mu.Unlock()
		return true
	} else if avg < r.alertthresh && alarmed {
		r.mu.Lock()
		r.alerts[len(r.alerts)-1].End = time.Now()
		r.mu.Unlock()
		return false
	}

	return alarmed
}

func extractSection(u *url.URL) string {
	c := path.Clean(u.Path)
	parts := strings.SplitN(c, "/", 3)
	return parts[0] + "/" + parts[1]
}
